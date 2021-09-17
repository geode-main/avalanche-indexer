package pvm

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer/shared"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/avalanche-indexer/util"
)

const (
	containerType  = "P"
	fetchBatchSize = 256
)

type Worker struct {
	chain       string
	avaxAsset   string
	codec       codec.Manager
	store       *store.DB
	indexClient *client.IndexClient
	log         *logrus.Logger
	status      *model.SyncStatus
}

type BlockData struct {
	CodecVersion uint16
	BlockID      string
	Block        *model.Block
	Transactions []*model.Transaction
	RewardsOwner *model.RewardsOwner
	Chain        *model.Chain
}

func NewWorker(
	indexClient *client.IndexClient,
	store *store.DB,
	codec codec.Manager,
	chain string,
	asset string,
) Worker {
	return Worker{
		chain:       chain,
		avaxAsset:   asset,
		store:       store,
		codec:       codec,
		log:         logrus.StandardLogger(),
		indexClient: indexClient,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.log.WithField("chain", w.chain).Info("starting worker")

	timer := time.NewTimer(time.Second)
	defer func() {
		timer.Stop()
		w.log.WithField("chain", w.chain).Info("worker stopped")
	}()

	for {
		select {
		case <-ctx.Done():
			w.log.WithField("chain", w.chain).Info("stopping worker")
			return
		case <-timer.C:
			if err := w.Run(); err != nil {
				w.log.WithField("chain", w.chain).WithError(err).Info("worker run failed")
				timer.Reset(time.Second)
				break
			}

			w.log.
				WithField("chain", w.chain).
				WithField("index", w.status.IndexID).
				WithField("lag", w.status.Lag()).
				Info("finished run")

			if w.status.AtTip() {
				timer.Reset(time.Second * 3)
			} else {
				timer.Reset(time.Millisecond * 10)
			}
		}
	}
}

func (w *Worker) Run() error {
	status, err := shared.GetSyncStatus(w.indexClient, w.store, containerType, w.chain)
	if err != nil {
		return err
	}
	w.status = status

	if w.status.AtTip() {
		return nil
	}

	return shared.ProcessContainerRange(
		w.status,
		w.indexClient,
		w.store,
		containerType,
		fetchBatchSize,
		w.ProcessMessage,
	)
}

func (w Worker) ProcessMessage(message *model.RawMessage) error {
	var (
		raw []byte
		err error
	)

	raw, err = formatting.Decode(formatting.Hex, message.Data)
	if err != nil {
		return err
	}

	blockData, err := w.prepareBlockData(raw, message.CreatedAt)
	if err != nil {
		return err
	}

	return w.saveBlockData(blockData)
}

func (w Worker) prepareBlockData(raw []byte, blockTime time.Time) (*BlockData, error) {
	var genericBlock platformvm.Block

	version, err := w.codec.Unmarshal(raw, &genericBlock)
	if err != nil {
		return nil, err
	}

	blockID, err := util.AvalancheID(raw)
	if err != nil {
		return nil, err
	}

	blockData := &BlockData{
		BlockID:      blockID,
		CodecVersion: version,
	}

	switch block := genericBlock.(type) {
	case *platformvm.ProposalBlock:
		err = w.buildBlock(blockData, model.BlockTypeProposal, blockTime, block.CommonBlock, &block.Tx)
	case *platformvm.StandardBlock:
		err = w.buildBlock(blockData, model.BlockTypeStandard, blockTime, block.CommonBlock, block.Txs...)
	case *platformvm.AtomicBlock:
		err = w.buildBlock(blockData, model.BlockTypeAtomic, blockTime, block.CommonBlock, &block.Tx)
	case *platformvm.AbortBlock:
		err = w.buildBlock(blockData, model.BlockTypeAbort, blockTime, block.CommonBlock)
	case *platformvm.CommitBlock:
		err = w.buildBlock(blockData, model.BlockTypeCommit, blockTime, block.CommonBlock)
	default:
		return nil, fmt.Errorf("unsupported block: %v %v", blockID, reflect.TypeOf(block))
	}

	return blockData, nil
}

func (w Worker) saveBlockData(data *BlockData) error {
	w.log.
		WithFields(logrus.Fields{"block": data.Block.ID, "height": data.Block.Height}).
		Debug("importing block")

	if err := w.store.Platform.CreateBlock(data.Block); err != nil {
		return err
	}

	for _, tx := range data.Transactions {
		w.log.
			WithFields(logrus.Fields{"type": tx.Type, "id": tx.ID}).
			Debug("importing transaction")

		if err := w.store.Platform.CreateTransaction(tx); err != nil {
			return err
		}

		if err := w.store.Platform.CreateTxOutputs(tx.Outputs); err != nil {
			return err
		}

		spentIDs := make([]string, len(tx.Inputs))
		for idx, input := range tx.Inputs {
			spentIDs[idx] = input.ID
		}

		if err := w.store.Platform.CreateTxInputs(spentIDs, tx.ID); err != nil {
			return err
		}

		if err := w.store.Platform.MarkOutputsSpent(spentIDs, tx.ID, tx.Timestamp); err != nil {
			return err
		}
	}

	if data.Chain != nil {
		if err := w.store.Platform.CreateChain(data.Chain); err != nil {
			return err
		}
	}

	if data.RewardsOwner != nil {
		w.log.WithField("tx", data.RewardsOwner.ID).Debug("creating rewards owner records")
		if err := w.store.Platform.CreateRewardsOwner(data.RewardsOwner); err != nil {
			return err
		}
	}

	return nil
}

func (w Worker) buildBlock(data *BlockData, blockType string, blockTime time.Time, block platformvm.CommonBlock, transactions ...*platformvm.Tx) error {
	data.Block = &model.Block{
		ID:        data.BlockID,
		Type:      blockType,
		Parent:    block.Parent().String(),
		Height:    block.Height(),
		Timestamp: blockTime,
		Chain:     w.chain,
	}

	for _, tx := range transactions {
		if err := w.initTransaction(data.CodecVersion, tx); err != nil {
			return err
		}

		if err := w.buildTransaction(data, tx); err != nil {
			return err
		}
	}

	return nil
}

func (w Worker) buildTransaction(data *BlockData, pvmTx *platformvm.Tx) error {
	var (
		transaction  *model.Transaction
		rewardsOwner *model.RewardsOwner
		chain        *model.Chain
		err          error
	)

	switch tx := pvmTx.UnsignedTx.(type) {
	case *platformvm.UnsignedCreateChainTx:
		transaction, chain, err = prepareCreateChainTx(tx)
	case *platformvm.UnsignedCreateSubnetTx:
		transaction, err = prepareCreateSubnetTx(tx)
	case *platformvm.UnsignedAddSubnetValidatorTx:
		transaction, err = prepareAddSubnetValidatorTx(tx)
	case *platformvm.UnsignedRewardValidatorTx:
		transaction, err = prepareRewardValidatorTx(tx)
	case *platformvm.UnsignedAddValidatorTx:
		transaction, rewardsOwner, err = prepareAddValidatorTx(tx)
	case *platformvm.UnsignedAddDelegatorTx:
		transaction, rewardsOwner, err = prepareAddDelegatorTx(tx)
	case *platformvm.UnsignedImportTx:
		transaction, err = prepareImportTx(tx)
	case *platformvm.UnsignedExportTx:
		transaction, err = prepareExportTx(tx)
	case *platformvm.UnsignedAdvanceTimeTx:
		transaction, err = prepareAdvanceTimeTx(tx)
	default:
		return fmt.Errorf("unsupported pvm tx type: %s", reflect.TypeOf(pvmTx.UnsignedTx))
	}

	if err != nil {
		return err
	}
	if transaction == nil {
		return nil
	}

	transaction.Status = model.TxStatusAccepted
	transaction.Chain = data.Block.Chain
	transaction.Block = &data.Block.ID
	transaction.BlockHeight = &data.Block.Height
	transaction.Timestamp = data.Block.Timestamp

	updateTransactionTotals(transaction, w.avaxAsset)

	data.Transactions = append(data.Transactions, transaction)
	data.RewardsOwner = rewardsOwner
	data.Chain = chain

	return nil
}

func (w Worker) initTransaction(codecVersion uint16, tx *platformvm.Tx) error {
	unsignedBytes, err := w.codec.Marshal(codecVersion, tx.UnsignedTx)
	if err != nil {
		return err
	}

	signedBytes, err := w.codec.Marshal(codecVersion, tx)
	if err != nil {
		return err
	}

	tx.Initialize(unsignedBytes, signedBytes)
	return nil
}
