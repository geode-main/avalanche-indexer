package shared

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/verify"
	"github.com/ava-labs/avalanchego/vms/nftfx"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"

	"github.com/figment-networks/avalanche-indexer/model"
)

var (
	ecdsaRecovery = crypto.FactorySECP256K1R{}
)

func PrepareInputs(inputs []*avax.TransferableInput, txID ids.ID) ([]model.Output, error) {
	result := make([]model.Output, len(inputs))

	for idx, input := range inputs {
		record, err := PrepareInput(input, idx, txID)
		if err != nil {
			return nil, err
		}
		result[idx] = record
	}

	return result, nil
}

func PrepareInput(input *avax.TransferableInput, idx int, txID ids.ID) (model.Output, error) {
	result := model.Output{
		ID:     input.TxID.Prefix(uint64(input.OutputIndex)).String(),
		TxID:   input.TxID.String(),
		Asset:  input.Asset.ID.String(),
		Amount: input.In.Amount(),
	}
	return result, nil
}

func PrepareOutputs(outputs []*avax.TransferableOutput, txID ids.ID) ([]model.Output, error) {
	result := make([]model.Output, len(outputs))

	for idx, out := range outputs {
		output, err := PrepareOutput(out.Out, out.Asset.ID.String(), idx, txID)
		if err != nil {
			return nil, err
		}
		result[idx] = *output
	}

	return result, nil
}

func PrepareOutput(output verify.State, asset string, idx int, txID ids.ID) (*model.Output, error) {
	outputID := txID.Prefix(uint64(idx)).String()

	result := &model.Output{
		ID:    outputID,
		TxID:  txID.String(),
		Index: uint64(idx),
		Asset: asset,
	}

	switch out := output.(type) {
	case *secp256k1fx.TransferOutput:
		addrs, err := bech32addrs(out.Addresses())
		if err != nil {
			return nil, err
		}

		result.Type = model.OutTypeTransfer
		result.Amount = out.Amt
		result.Locktime = out.Locktime
		result.Threshold = out.Threshold
		result.Addresses = addrs

	case *secp256k1fx.MintOutput:
		addrs, err := bech32addrs(out.Addresses())
		if err != nil {
			return nil, err
		}

		result.Type = model.OutTypeMint
		result.Locktime = out.Locktime
		result.Threshold = out.Threshold
		result.Addresses = addrs

	case *platformvm.StakeableLockOut:
		addrs, err := bech32addrs(out.Addresses())
		if err != nil {
			return nil, err
		}

		result.Type = model.OutTypeStakeableLock
		result.Amount = out.Amount()
		result.Locktime = out.Locktime
		result.Addresses = addrs

	case *nftfx.MintOutput:
		addrs, err := bech32addrs(out.Addresses())
		if err != nil {
			return nil, err
		}

		result.Type = model.OutTypeNftMint
		result.Locktime = out.Locktime
		result.Threshold = out.Threshold
		result.Group = out.GroupID
		result.Addresses = addrs

	case *nftfx.TransferOutput:
		payload := base64.StdEncoding.EncodeToString(out.Payload)

		addrs, err := bech32addrs(out.Addresses())
		if err != nil {
			return nil, err
		}

		result.Type = model.OutTypeNftTransfer
		result.Locktime = out.Locktime
		result.Threshold = out.Threshold
		result.Group = out.GroupID
		result.Addresses = addrs
		result.Payload = &payload

	default:
		return nil, fmt.Errorf("unknown output type: %s (type=%s)", outputID, reflect.TypeOf(out))
	}

	return result, nil
}

func PrepareTxCredentials(verifiables []verify.Verifiable, unsignedBytes []byte) ([]model.Credential, error) {
	result := make([]model.Credential, len(verifiables))

	for idx, ver := range verifiables {
		cred, ok := ver.(*secp256k1fx.Credential)
		if !ok {
			continue
		}

		for _, sig := range cred.Sigs {
			publicKey, err := ecdsaRecovery.RecoverPublicKey(unsignedBytes, sig[:])
			if err != nil {
				return nil, err
			}

			addr, err := bech32address(publicKey.Address())
			if err != nil {
				panic(err)
			}

			result[idx].Address = addr
			result[idx].PublicKey = base64.StdEncoding.EncodeToString(publicKey.Bytes())
			result[idx].Signature = base64.StdEncoding.EncodeToString(sig[:])
		}
	}

	return result, nil
}

func PrepareRewardsOwner(txID ids.ID, outputsCount int, owner verify.Verifiable) (*model.RewardsOwner, error) {
	outputOwners, ok := owner.(*secp256k1fx.OutputOwners)
	if !ok {
		return nil, errors.New("invalid rewards owner type")
	}

	addrs, err := bech32addrs(outputOwners.Addresses())
	if err != nil {
		return nil, err
	}

	rewardsAddresses := make([]model.RewardsOwnerAddress, len(addrs))
	for idx, addr := range addrs {
		rewardsAddresses[idx].ID = txID.String()
		rewardsAddresses[idx].Address = addr
		rewardsAddresses[idx].Index = uint32(idx)
	}

	rewardsOutputs := []model.RewardsOwnerOutput{}
	for pos := outputsCount; pos < outputsCount+2; pos++ {
		rewardsOutputs = append(rewardsOutputs, model.RewardsOwnerOutput{
			ID:            txID.Prefix(uint64(pos)).String(),
			TransactionID: txID.String(),
			Index:         uint32(pos),
		})
	}

	return &model.RewardsOwner{
		ID:        txID.String(),
		Locktime:  outputOwners.Locktime,
		Threshold: outputOwners.Threshold,
		Addresses: rewardsAddresses,
		Outputs:   rewardsOutputs,
	}, nil
}
