package avm

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/avalanchego/vms/avm"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer/shared"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store"
)

const (
	containerType  = "X"
	fetchBatchSize = 256
)

type Worker struct {
	chain       string
	avaxAsset   string
	status      *model.SyncStatus
	codec       codec.Manager
	store       *store.DB
	indexClient *client.IndexClient
	log         *logrus.Logger
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
	w.log.WithField("index", message.IndexID).Debug("processing message")

	var (
		genericTx avm.Tx
		tx        *model.Transaction
	)

	raw, err := formatting.Decode(formatting.Hex, message.Data)
	if err != nil {
		return err
	}

	version, err := w.codec.Unmarshal(raw, &genericTx)
	if err != nil {
		return err
	}

	unsignedBytes, err := w.codec.Marshal(version, &genericTx.UnsignedTx)
	if err != nil {
		return err
	}
	genericTx.Initialize(unsignedBytes, raw)

	switch typedTx := genericTx.UnsignedTx.(type) {
	case *avm.BaseTx:
		tx, err = prepareBaseTx(typedTx)
	case *avm.ImportTx:
		tx, err = prepareImportTx(typedTx)
	case *avm.ExportTx:
		tx, err = prepareExportTx(typedTx)
	case *avm.CreateAssetTx:
		tx, err = prepareCreateAssetTx(typedTx)
	case *avm.OperationTx:
		tx, err = prepareOperationTx(typedTx)
	default:
		return fmt.Errorf("unsupported tx type: %v", typedTx)
	}

	if err != nil {
		return err
	}
	if tx == nil {
		return fmt.Errorf("tx %s was not decoded", genericTx.ID().String())
	}

	tx.Status = model.TxStatusAccepted
	tx.Chain = w.chain
	tx.Timestamp = message.CreatedAt

	for idx := range tx.Outputs {
		tx.Outputs[idx].Chain = w.chain
	}

	updateTransactionTotals(tx, w.avaxAsset)

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

	switch tx.Type {
	case model.TxTypeCreateAsset:
		if err := w.createAssetFromTx(tx); err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) createAssetFromTx(tx *model.Transaction) error {
	asset := &model.Asset{
		AssetID:      tx.Metadata.GetString("asset_id"),
		Type:         tx.Metadata.GetString("asset_type"),
		Name:         tx.Metadata.GetString("asset_name"),
		Symbol:       tx.Metadata.GetString("asset_symbol"),
		Denomination: tx.Metadata.GetInt("asset_denomination"),
	}

	return w.store.Assets.Create(asset)
}
