package avm

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/avm"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"

	"github.com/figment-networks/avalanche-indexer/indexer/shared"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/model/types"
	"github.com/figment-networks/avalanche-indexer/util"
)

func prepareBaseTx(tx *avm.BaseTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:   tx.ID().String(),
		Type: model.TxTypeBase,
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(transaction, tx.BaseTx.ID(), tx.Ins, tx.Outs)
}

func prepareImportTx(tx *avm.ImportTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:          tx.ID().String(),
		Type:        model.TxTypeXImport,
		SourceChain: util.StringPtr(tx.SourceChain.String()),
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(
		transaction,
		tx.BaseTx.ID(),
		append(tx.Ins, tx.ImportedIns...), tx.Outs,
	)
}

func prepareExportTx(tx *avm.ExportTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:               tx.ID().String(),
		Type:             model.TxTypeXExport,
		DestinationChain: util.StringPtr(tx.DestinationChain.String()),
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(
		transaction,
		tx.BaseTx.ID(),
		tx.Ins,
		append(tx.Outs, tx.ExportedOuts...),
	)
}

func prepareCreateAssetTx(tx *avm.CreateAssetTx) (*model.Transaction, error) {
	assetID := tx.ID().String()
	assetType := model.AssetTypeFixed
	assetInitialSupply := uint64(0)

	inputs, err := shared.PrepareInputs(tx.Ins, tx.BaseTx.ID())
	if err != nil {
		return nil, err
	}

	outputs, err := shared.PrepareOutputs(tx.Outs, tx.BaseTx.ID())
	if err != nil {
		return nil, err
	}

	startIdx := len(outputs)

	for _, state := range tx.States {
		for idx, stateOut := range state.Outs {
			out, err := shared.PrepareOutput(stateOut, tx.ID().String(), startIdx+idx, tx.BaseTx.ID())
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, *out)
		}
	}

	for _, out := range outputs {
		if assetID != out.Asset {
			continue
		}

		switch out.Type {
		case model.OutTypeNftMint, model.OutTypeNftTransfer:
			assetType = model.AssetTypeNFT
			break
		case model.OutTypeMint:
			assetType = model.AssetTypeVariable
			assetInitialSupply += out.Amount
		}
	}

	transaction := &model.Transaction{
		ID:   assetID,
		Type: model.TxTypeCreateAsset,
		Metadata: types.Map{
			"asset_id":             assetID,
			"asset_type":           assetType,
			"asset_name":           tx.Name,
			"asset_denomination":   int(tx.Denomination),
			"asset_symbol":         tx.Symbol,
			"asset_initial_supply": assetInitialSupply,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}

	transaction.SetRawMemo(tx.Memo)

	return transaction, nil
}

func prepareOperationTx(tx *avm.OperationTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:    tx.ID().String(),
		Chain: tx.BlockchainID.String(),
		Type:  model.TxTypeOperation,
	}
	transaction.SetRawMemo(tx.Memo)

	inputs, err := shared.PrepareInputs(tx.Ins, tx.BaseTx.ID())
	if err != nil {
		return nil, err
	}

	outputs, err := shared.PrepareOutputs(tx.Outs, tx.BaseTx.ID())
	if err != nil {
		return nil, err
	}

	startIdx := len(outputs)
	for _, op := range tx.Ops {
		for idx, utxo := range op.UTXOIDs {
			in, err := shared.PrepareInput(&avax.TransferableInput{
				Asset:  op.Asset,
				UTXOID: *utxo,
				In:     &secp256k1fx.TransferInput{},
			}, len(inputs)+idx, tx.BaseTx.ID())
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, in)
		}

		for idx, txOut := range op.Op.Outs() {
			out, err := shared.PrepareOutput(txOut, op.AssetID().String(), startIdx+idx, tx.BaseTx.ID())
			if err != nil {
				return nil, err
			}
			outputs = append(outputs, *out)
		}
	}

	transaction.Inputs = inputs
	transaction.Outputs = outputs

	return transaction, err
}

func setTxInsOuts(tx *model.Transaction, txID ids.ID, ins []*avax.TransferableInput, outs []*avax.TransferableOutput) (*model.Transaction, error) {
	inputs, err := shared.PrepareInputs(ins, txID)
	if err != nil {
		return nil, err
	}

	outputs, err := shared.PrepareOutputs(outs, txID)
	if err != nil {
		return nil, err
	}

	tx.Inputs = inputs
	tx.Outputs = outputs

	return tx, nil
}

func updateTransactionTotals(tx *model.Transaction, avaxAssetID string) {
	if tx.InputAmounts == nil {
		tx.InputAmounts = map[string]uint64{}
	}
	if tx.OutputAmounts == nil {
		tx.OutputAmounts = map[string]uint64{}
	}

	totalIn := uint64(0)
	totalOut := uint64(0)

	for _, in := range tx.Inputs {
		tx.InputAmounts[in.Asset] += in.Amount

		if in.Asset == avaxAssetID {
			totalIn += in.Amount
		}
	}

	for _, out := range tx.Outputs {
		tx.OutputAmounts[out.Asset] += out.Amount

		if out.Asset == avaxAssetID {
			totalOut += out.Amount
		}
	}

	fee := totalIn - totalOut
	if fee < 0 {
		fee = 0
	}
	tx.Fee = fee
}
