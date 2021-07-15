package cvm

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ava-labs/avalanchego/codec"
	"github.com/ava-labs/avalanchego/utils/formatting"
	"github.com/ava-labs/coreth/plugin/evm"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"

	corethTypes "github.com/ava-labs/coreth/core/types"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/indexer/shared"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/model/types"
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/figment-networks/avalanche-indexer/util"
)

const (
	containerType  = "C"
	fetchBatchSize = 256
)

type Worker struct {
	chain       string
	avaxAsset   string
	store       *store.DB
	codec       codec.Manager
	indexClient *client.IndexClient
	evmClient   *client.EvmClient
	status      *model.SyncStatus
	log         *logrus.Logger
	ethChainID  *big.Int
	ethSigner   corethTypes.Signer
}

func NewWorker(
	store *store.DB,
	codec codec.Manager,
	indexClient *client.IndexClient,
	evmClient *client.EvmClient,
	chain string,
	avaxAsset string,
	ethChainID *big.Int,
) Worker {
	return Worker{
		chain:       chain,
		avaxAsset:   avaxAsset,
		store:       store,
		codec:       codec,
		indexClient: indexClient,
		evmClient:   evmClient,
		log:         logrus.StandardLogger(),
		ethChainID:  ethChainID,
		ethSigner:   corethTypes.NewEIP155Signer(ethChainID),
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

	block := &corethTypes.Block{}
	if err := rlp.DecodeBytes(raw, block); err != nil {
		return err
	}

	ourBlock, err := w.prepareBlock(block)
	if err != nil {
		return err
	}

	atomicTx, err := w.prepareTx(block)
	if err != nil {
		return err
	}
	if atomicTx != nil {
		if err := w.saveAtomicTx(atomicTx); err != nil {
			return err
		}
	}

	if err := w.store.Platform.CreateBlock(ourBlock); err != nil {
		return err
	}

	for _, ethTx := range block.Transactions() {
		tx, err := w.prepareEvmTx(block, ethTx)
		if err != nil {
			return err
		}

		if err := w.store.Platform.CreateTransaction(tx); err != nil {
			return err
		}
	}

	return nil
}

func (w Worker) prepareBlock(ethBlock *corethTypes.Block) (*model.Block, error) {
	return &model.Block{
		ID:        ethBlock.Hash().String(),
		Parent:    ethBlock.ParentHash().String(),
		Chain:     w.chain,
		Height:    ethBlock.NumberU64(),
		Type:      model.BlockTypeEvm,
		Timestamp: time.Unix(int64(ethBlock.Time()), 0),
	}, nil
}

func (w Worker) prepareTx(block *corethTypes.Block) (*model.Transaction, error) {
	rawBytes := block.ExtData()
	if len(rawBytes) == 0 {
		return nil, nil
	}

	evmTx := &evm.Tx{}

	version, err := w.codec.Unmarshal(rawBytes, evmTx)
	if err != nil {
		return nil, err
	}

	unsignedBytes, err := w.codec.Marshal(version, &evmTx.UnsignedAtomicTx)
	if err != nil {
		return nil, err
	}
	evmTx.Initialize(unsignedBytes, rawBytes)

	var tx *model.Transaction

	switch typedTx := evmTx.UnsignedAtomicTx.(type) {
	case *evm.UnsignedImportTx:
		tx, err = w.prepareAtomicImportTx(typedTx)
	case *evm.UnsignedExportTx:
		tx, err = w.prepareAtomicExportTx(typedTx)
	default:
		err = fmt.Errorf("unsupported transaction type: %s", reflect.TypeOf(tx).String())
	}

	if err != nil {
		return nil, err
	}

	blockHash := block.Hash().String()
	blockHeight := block.Number().Uint64()

	tx.Status = model.TxStatusAccepted
	tx.Block = &blockHash
	tx.BlockHeight = &blockHeight
	tx.Timestamp = time.Unix(int64(block.Time()), 0)

	return tx, nil
}

func (w Worker) saveAtomicTx(tx *model.Transaction) error {
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

	return nil
}

func (w Worker) prepareAtomicImportTx(tx *evm.UnsignedImportTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:          tx.ID().String(),
		Type:        model.TxTypeAtomicImport,
		Chain:       tx.BlockchainID.String(),
		SourceChain: util.StringPtr(tx.SourceChain.String()),
	}

	ins, err := shared.PrepareInputs(tx.ImportedInputs, tx.ID())
	if err != nil {
		return nil, err
	}
	transaction.Inputs = ins

	for idx, txOut := range tx.Outs {
		out := model.Output{
			ID:        tx.ID().Prefix(uint64(idx)).String(),
			Chain:     transaction.Chain,
			Type:      model.OutTypeTransfer,
			TxID:      tx.ID().String(),
			Index:     uint64(idx),
			Asset:     txOut.AssetID.String(),
			Amount:    txOut.Amount,
			Addresses: []string{txOut.Address.String()},
		}
		transaction.Outputs = append(transaction.Outputs, out)
	}

	return transaction, nil
}

func (w Worker) prepareAtomicExportTx(tx *evm.UnsignedExportTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:               tx.ID().String(),
		Type:             model.TxTypeAtomicExport,
		Chain:            tx.BlockchainID.String(),
		DestinationChain: util.StringPtr(tx.DestinationChain.String()),
	}

	for idx, txIn := range tx.Ins {
		in := model.Output{
			ID:        tx.ID().Prefix(uint64(idx)).String(),
			Type:      model.OutTypeTransfer,
			TxID:      tx.ID().String(),
			Index:     uint64(idx),
			Asset:     txIn.AssetID.String(),
			Amount:    txIn.Amount,
			Addresses: []string{txIn.Address.String()},
		}
		transaction.Inputs = append(transaction.Inputs, in)
	}

	outs, err := shared.PrepareOutputs(tx.ExportedOutputs, tx.ID())
	if err != nil {
		return nil, err
	}
	transaction.Outputs = outs

	return transaction, nil
}

func (w Worker) prepareEvmTx(block *corethTypes.Block, ethTx *corethTypes.Transaction) (*model.Transaction, error) {
	msg, err := ethTx.AsMessage(w.ethSigner)
	if err != nil {
		return nil, err
	}
	nonce := msg.Nonce()

	meta := types.Map{
		"sender":    msg.From(),
		"receiver":  msg.To(),
		"nonce":     nonce,
		"amount":    ethTx.Value().String(),
		"gas":       ethTx.Gas(),
		"gas_price": ethTx.GasPrice().String(),
		"cost":      ethTx.Cost().String(),
	}

	tx := &model.Transaction{
		ID:          ethTx.Hash().String(),
		Type:        model.TxTypeEvm,
		Status:      model.TxStatusAccepted,
		Timestamp:   time.Unix(int64(block.Time()), 0),
		Chain:       w.chain,
		Block:       util.StringPtr(block.Hash().String()),
		BlockHeight: util.Uint64Prt(block.NumberU64()),
		Nonce:       &nonce,
		Metadata:    meta,
	}

	return tx, nil
}
