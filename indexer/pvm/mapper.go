package pvm

import (
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/platformvm"

	"github.com/figment-networks/avalanche-indexer/indexer/shared"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/model/types"
	"github.com/figment-networks/avalanche-indexer/util"
)

func prepareAddValidatorTx(tx *platformvm.UnsignedAddValidatorTx) (*model.Transaction, *model.RewardsOwner, error) {
	transaction := &model.Transaction{
		ID:       tx.ID().String(),
		Type:     model.TxTypeAddValidator,
		Metadata: makeStakingMetadata(tx.Validator),
	}
	transaction.SetRawMemo(tx.Memo)
	transaction.Metadata["validator_commission_rate"] = tx.Shares / 10000

	_, err := setTxInsOuts(transaction, tx.BaseTx.ID(), tx.Ins, append(tx.Outs, tx.Stake...))
	if err != nil {
		return nil, nil, err
	}

	for idx := range transaction.Outputs {
		for _, stakeOut := range tx.Stake {
			if stakeOut.ID.String() == transaction.Outputs[idx].ID {
				transaction.Outputs[idx].Stake = true
			}
		}
	}

	rewardsOwner, err := shared.PrepareRewardsOwner(tx.ID(), len(transaction.Outputs), tx.RewardsOwner)
	if err != nil {
		return nil, nil, err
	}

	return transaction, rewardsOwner, nil
}

func prepareRewardValidatorTx(tx *platformvm.UnsignedRewardValidatorTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:            tx.ID().String(),
		Type:          model.TxTypeRewardValidator,
		ReferenceTxID: util.StringPtr(tx.TxID.String()),
	}
	return transaction, nil
}

func prepareAddDelegatorTx(tx *platformvm.UnsignedAddDelegatorTx) (*model.Transaction, *model.RewardsOwner, error) {
	transaction := &model.Transaction{
		ID:       tx.ID().String(),
		Type:     model.TxTypeAddDelegator,
		Metadata: makeStakingMetadata(tx.Validator),
	}
	transaction.SetRawMemo(tx.Memo)

	for idx, txInput := range tx.Ins {
		input, err := shared.PrepareInput(txInput, idx, tx.BaseTx.ID())
		if err != nil {
			return nil, nil, err
		}
		transaction.Inputs = append(transaction.Inputs, input)
	}

	outIdx := 0
	for _, txOut := range tx.Outs {
		out, err := shared.PrepareOutput(txOut.Out, txOut.Asset.ID.String(), outIdx, tx.ID())
		if err != nil {
			return nil, nil, err
		}

		transaction.Outputs = append(transaction.Outputs, *out)
		outIdx++
	}

	for _, txOut := range tx.Stake {
		out, err := shared.PrepareOutput(txOut.Out, txOut.Asset.ID.String(), outIdx, tx.ID())
		if err != nil {
			return nil, nil, err
		}

		out.Stake = true
		transaction.Outputs = append(transaction.Outputs, *out)
		outIdx++
	}

	rewardsOwner, err := shared.PrepareRewardsOwner(tx.ID(), len(transaction.Outputs), tx.RewardsOwner)
	if err != nil {
		return nil, nil, err
	}

	return transaction, rewardsOwner, nil
}

func prepareCreateChainTx(tx *platformvm.UnsignedCreateChainTx) (*model.Transaction, *model.Chain, error) {
	chain := &model.Chain{
		ChainID: tx.ID().String(),
		Name:    tx.ChainName,
		VM:      tx.VMID.String(),
		Network: tx.NetworkID,
		Subnet:  tx.SubnetID.String(),
	}

	transaction := &model.Transaction{
		ID:       tx.ID().String(),
		Type:     model.TxTypeCreateChain,
		Metadata: makeChainMetadata(tx),
	}
	transaction.SetRawMemo(tx.Memo)

	if _, err := setTxInsOuts(transaction, tx.ID(), tx.Ins, tx.Outs); err != nil {
		return nil, nil, err
	}

	return transaction, chain, nil
}

func prepareCreateSubnetTx(tx *platformvm.UnsignedCreateSubnetTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:   tx.ID().String(),
		Type: model.TxTypeCreateSubnet,
		Metadata: types.Map{
			"subnet_id": tx.ID(),
		},
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(transaction, tx.ID(), tx.Ins, tx.Outs)
}

func prepareAddSubnetValidatorTx(tx *platformvm.UnsignedAddSubnetValidatorTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:   tx.ID().String(),
		Type: model.TxTypeAddSubnetValidator,
		Metadata: types.Map{
			"validator_node_id": tx.Validator.NodeID,
			"subnet_id":         tx.Validator.Subnet.String(),
		},
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(transaction, tx.ID(), tx.Ins, tx.Outs)
}

func prepareImportTx(tx *platformvm.UnsignedImportTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:          tx.ID().String(),
		Type:        model.TxTypePImport,
		SourceChain: util.StringPtr(tx.SourceChain.String()),
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(transaction, tx.BaseTx.ID(), append(tx.Ins, tx.ImportedInputs...), tx.Outs)
}

func prepareExportTx(tx *platformvm.UnsignedExportTx) (*model.Transaction, error) {
	transaction := &model.Transaction{
		ID:               tx.ID().String(),
		Type:             model.TxTypePExport,
		DestinationChain: util.StringPtr(tx.DestinationChain.String()),
	}
	transaction.SetRawMemo(tx.Memo)

	return setTxInsOuts(transaction, tx.BaseTx.ID(), tx.Ins, append(tx.Outs, tx.ExportedOutputs...))
}

func prepareAdvanceTimeTx(tx *platformvm.UnsignedAdvanceTimeTx) (*model.Transaction, error) {
	return nil, nil
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

func makeStakingMetadata(validator platformvm.Validator) types.Map {
	return types.Map{
		"node_id":    validator.NodeID.String(),
		"start_time": validator.StartTime().UTC().Format(time.RFC3339),
		"end_time":   validator.EndTime().UTC().Format(time.RFC3339),
		"duration":   validator.EndTime().Unix() - validator.StartTime().Unix(),
		"weight":     validator.Weight(),
	}
}

func makeChainMetadata(tx *platformvm.UnsignedCreateChainTx) types.Map {
	return types.Map{
		"chain_name":       tx.ChainName,
		"chain_id":         tx.ID().String(),
		"chain_vm_id":      tx.VMID.String(),
		"chain_network_id": tx.NetworkID,
		"chain_subnet_id":  tx.SubnetID,
	}
}

func setTxInsOuts(
	tx *model.Transaction, txID ids.ID, ins []*avax.TransferableInput, outs []*avax.TransferableOutput) (*model.Transaction, error) {
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

type inputGroup struct {
	inputs []*avax.TransferableInput
}

type outputGroup struct {
	outputs []*avax.TransferableInput
	stake   bool
	reward  bool
}
