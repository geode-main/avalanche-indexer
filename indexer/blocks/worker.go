package blocks

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/model/types"
	"github.com/figment-networks/avalanche-indexer/store"
)

type Worker struct {
	log           *logrus.Logger
	rpc           *client.Client
	db            *store.DB
	chain         string
	status        *model.SyncStatus
	syncStatusKey string

	errWaitTime time.Duration
	syncTime    time.Duration
	cycleTime   time.Duration
}

func NewWorker(db *store.DB, rpc *client.Client, log *logrus.Logger, chain string) Worker {
	return Worker{
		db:            db,
		rpc:           rpc,
		log:           log,
		chain:         chain,
		syncStatusKey: fmt.Sprintf("%s_events", chain),

		errWaitTime: time.Second,
		syncTime:    time.Second * 3,
		cycleTime:   time.Millisecond * 10,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.log.WithField("chain", w.syncStatusKey).Info("starting worker")

	timer := time.NewTimer(time.Second)
	defer func() {
		timer.Stop()
		w.log.WithField("chain", w.syncStatusKey).Info("worker stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			w.log.WithField("chain", w.syncStatusKey).Info("stopping worker")
			return
		case <-timer.C:
			if err := w.Run(); err != nil {
				w.log.WithField("chain", w.syncStatusKey).WithError(err).Info("worker run failed")
				timer.Reset(w.errWaitTime)
				break
			}

			w.log.
				WithField("chain", w.syncStatusKey).
				WithField("index", w.status.IndexID).
				WithField("lag", w.status.Lag()).
				Info("finished run")

			if w.status.AtTip() {
				timer.Reset(w.syncTime)
			} else {
				timer.Reset(w.cycleTime)
			}
		}
	}
}

func (w *Worker) Run() error {
	status, err := w.getSyncStatus()
	if err != nil {
		return err
	}
	w.status = status

	if w.status.AtTip() {
		return nil
	}

	search := &store.BlocksSearch{
		Chain:       w.chain,
		StartHeight: int(w.status.NextID()),
		Order:       "height_asc",
	}

	blocks, err := w.db.Platform.GetBlocks(search)
	if err != nil {
		return err
	}
	if len(blocks) == 0 {
		w.log.
			WithFields(logrus.Fields{"chain": w.chain, "start_height": search.StartHeight}).
			Debug("no new blocks found")

		return nil
	}

	for _, block := range blocks {
		txResult, err := w.db.Transactions.Search(store.TxSearchInput{
			BlockHash: block.Parent,
		})
		if err != nil {
			return err
		}

		if len(txResult.Transactions) > 0 {
			if err := w.processBlockTx(&block, &txResult.Transactions[0]); err != nil {
				return err
			}
		}

		w.status.IndexID = int64(block.Height)
		w.status.IndexTime = block.Timestamp
	}

	return w.db.Platform.UpdateSyncStatus(status)
}

func (w *Worker) getSyncStatus() (*model.SyncStatus, error) {
	lastBlock, err := w.db.Platform.LastBlock(w.chain)
	if err != nil {
		return nil, err
	}

	status, err := w.db.Platform.GetSyncStatus(w.syncStatusKey)
	if err != nil {
		if err != store.ErrNotFound {
			return nil, err
		}

		status = &model.SyncStatus{
			ID:      w.syncStatusKey,
			IndexID: 0,
			TipID:   int64(lastBlock.Height),
			TipTime: lastBlock.Timestamp,
		}

		if err := w.db.Platform.UpdateSyncStatus(status); err != nil {
			return nil, err
		}
	}

	status.TipID = int64(lastBlock.Height)
	status.TipTime = lastBlock.Timestamp

	return status, nil
}

func (w Worker) processBlockTx(block *model.Block, tx *model.Transaction) error {
	switch tx.Type {
	case model.TxTypeAddValidator:
		return w.createAddValidatorEvent(block, tx)
	case model.TxTypeAddDelegator:
		return w.createAddDelegatorEvent(block, tx)
	case model.TxTypeRewardValidator:
		return w.createFinishValidatorEvent(block, tx)
	case model.TxTypeAddSubnetValidator:
		return w.createAddSubnetValidatorEvent(block, tx)
	default:
		return nil
	}
}

func (w Worker) initEvent(block *model.Block, tx *model.Transaction) *model.Event {
	return &model.Event{
		Chain:       block.Chain,
		BlockHash:   block.ID,
		BlockHeight: block.Height,
		TxHash:      tx.ID,
		Timestamp:   block.Timestamp,
	}
}

func (w Worker) initCommissionChangeEvent(block *model.Block, tx *model.Transaction, current *model.Event, prev *model.Event) *model.Event {
	if current.Type != prev.Type {
		return nil
	}

	beforeVal := prev.Data.GetInt("commission_rate")
	afterVal := current.Data.GetInt("commission_rate")

	if afterVal == beforeVal {
		return nil
	}

	commEvent := w.initEvent(block, tx)
	commEvent.Scope = current.Scope
	commEvent.Type = model.EventTypeValidatorCommissionChanged
	commEvent.ItemID = current.ItemID
	commEvent.ItemType = current.ItemType

	commEvent.Data = types.NewMap()
	commEvent.Data["before"] = beforeVal
	commEvent.Data["after"] = afterVal
	commEvent.Data["change"] = afterVal - beforeVal

	return commEvent
}

func (w Worker) createAddValidatorEvent(block *model.Block, tx *model.Transaction) error {
	if block.Type != model.BlockTypeCommit {
		return nil
	}

	event := w.initEvent(block, tx)
	event.Scope = model.EventScopeStaking
	event.Type = model.EventTypeValidatorAdded
	event.ItemID = tx.Metadata.GetString("node_id")
	event.ItemType = model.EventItemTypeValidator
	event.Data = tx.Metadata

	recentEvents, err := w.db.Events.Search(&store.EventSearchInput{
		Type:      model.EventTypeValidatorAdded,
		ItemID:    event.ItemID,
		ItemType:  event.ItemType,
		EndHeight: int(event.BlockHeight),
		Limit:     1,
	})
	if err != nil {
		return err
	}

	var commEvent *model.Event
	if len(recentEvents) > 0 {
		commEvent = w.initCommissionChangeEvent(block, tx, event, &recentEvents[0])
	}

	err = w.createEvent(event)

	if err == nil && commEvent != nil {
		err = w.createEvent(commEvent)
	}

	return err
}

func (w Worker) createAddDelegatorEvent(block *model.Block, tx *model.Transaction) error {
	if block.Type != model.BlockTypeCommit {
		return nil
	}

	event := w.initEvent(block, tx)
	event.Scope = model.EventScopeStaking
	event.Type = model.EventTypeDelegatorAdded
	event.ItemID = tx.Metadata.GetString("node_id")
	event.ItemType = model.EventItemTypeValidator
	event.Data = tx.Metadata

	return w.createEvent(event)
}

func (w Worker) createFinishValidatorEvent(block *model.Block, tx *model.Transaction) error {
	if tx.ReferenceTxID == nil {
		return nil
	}

	refTx, err := w.db.Transactions.GetByID(*tx.ReferenceTxID)
	if err != nil && err != store.ErrNotFound {
		return err
	}
	if refTx == nil {
		return nil
	}

	event := w.initEvent(block, tx)
	event.Scope = model.EventScopeStaking
	event.ItemID = refTx.Metadata.GetString("node_id")
	event.ItemType = model.EventItemTypeValidator
	event.Data = types.NewMap()
	event.Data["rewarded"] = block.Type == model.BlockTypeCommit

	switch refTx.Type {
	case model.TxTypeAddValidator:
		event.Type = model.EventTypeValidatorFinished
	case model.TxTypeAddDelegator:
		event.Type = model.EventTypeDelegatorFinished
	default:
		return fmt.Errorf("unhandled reward validator tx type: %s", refTx.Type)
	}

	return w.createEvent(event)
}

func (w Worker) createAddSubnetValidatorEvent(block *model.Block, tx *model.Transaction) error {
	event := w.initEvent(block, tx)
	event.Scope = model.EventScopeNetwork
	event.Type = model.EventTypeSubnetValidatorAdded
	event.ItemID = tx.Metadata.GetString("validator_node_id")
	event.ItemType = model.EventItemTypeValidator

	return w.createEvent(event)
}

func (w Worker) createEvent(event *model.Event) error {
	w.log.
		WithField("chain", w.syncStatusKey).
		WithField("type", event.Type).
		Debug("creating event")

	return w.db.Events.Create(event)
}
