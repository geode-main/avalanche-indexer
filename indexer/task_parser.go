package indexer

import (
	"context"

	"github.com/figment-networks/indexing-engine/pipeline"
	"github.com/sirupsen/logrus"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/util"
)

type ParserTask struct {
	logger *logrus.Logger
}

func (t *ParserTask) GetName() string {
	return taskParser
}

func (t *ParserTask) Run(ctx context.Context, p pipeline.Payload) error {
	logStart(t, t.logger)
	defer logDone(t, t.logger)

	payload := p.(*Payload)

	return runChain(
		payload,
		t.prepareValidatorShares,
		t.prepareValidators,
		t.prepareDelegations,
		t.parseMinStake,
		t.parseTxFee,
		t.prepareNetworkMetric,
	)
}

func (t ParserTask) prepareValidatorShares(payload *Payload) error {
	activeStake := int64(0)
	activeAmounts := map[string]int64{}

	for _, validator := range payload.CurrentValidators {
		val, err := util.ParseInt64(validator.StakeAmount)
		if err != nil {
			return err
		}
		activeAmounts[validator.NodeID] = val
		activeStake += val
	}
	payload.ActiveValidatorShare = map[string]float64{}

	for _, validator := range payload.CurrentValidators {
		val := util.PercentOf(activeAmounts[validator.NodeID], activeStake)
		payload.ActiveValidatorShare[validator.NodeID] = val
	}
	payload.ActiveStakeAmount = activeStake

	return nil
}

// parseMinStake parses the min validator and delegator stake values
func (t ParserTask) parseMinStake(payload *Payload) error {
	minValidatorStake, err := util.ParseInt64(payload.MinStake.MinValidatorStake)
	if err != nil {
		return err
	}

	minDelegatorStake, err := util.ParseInt64(payload.MinStake.MinDelegatorStake)
	if err != nil {
		return err
	}

	payload.MinValidatorStake = minValidatorStake
	payload.MinDelegatorStake = minDelegatorStake

	return nil
}

// parseTxFee parses the current fee amounts
func (t ParserTask) parseTxFee(payload *Payload) error {
	txFee, err := util.ParseInt64(payload.RawTxFee.TxFee)
	if err != nil {
		return err
	}

	creationTxFee, err := util.ParseInt64(payload.RawTxFee.CreationTxFee)
	if err != nil {
		return err
	}

	payload.TxFee = txFee
	payload.CreationTxFee = creationTxFee

	return nil
}

// prepareNetworkMetric builds a new network metric record
func (t ParserTask) prepareNetworkMetric(payload *Payload) error {
	validatorsCount := len(payload.CurrentValidators)
	avgUptime := float64(0)
	avgDelegationFee := float64(0)

	for _, validator := range payload.CurrentValidators {
		uptime, err := util.ParseFloat32(validator.Uptime)
		if err != nil {
			return err
		}

		fee, err := util.ParseFloat32(validator.DelegationFee)
		if err != nil {
			return err
		}

		avgUptime += uptime * 100
		avgDelegationFee += fee
	}

	avgUptime = avgUptime / float64(validatorsCount)
	avgDelegationFee = avgDelegationFee / float64(validatorsCount)

	payload.NetworkMetric = &model.NetworkMetric{
		Time:                   payload.SyncTime,
		Height:                 payload.Height,
		PeersCount:             len(payload.Peers),
		BlockchainsCount:       len(payload.Blockchains),
		ActiveValidatorsCount:  len(payload.CurrentValidators),
		PendingValidatorsCount: len(payload.PendingValidators),
		MinValidatorStake:      payload.MinValidatorStake,
		MinDelegationStake:     payload.MinDelegatorStake,
		TxFee:                  int(payload.TxFee),
		CreationTxFee:          int(payload.CreationTxFee),
		Uptime:                 avgUptime,
		DelegationFee:          avgDelegationFee,
	}

	return nil
}

// prepareValidators builds a new set of validator records
func (t ParserTask) prepareValidators(payload *Payload) error {
	delegationsCount := 0
	delegatedAmount := int64(0)
	delegatedAmountMap := map[string]int64{}

	for _, validator := range payload.CurrentValidators {
		validatorAmount := int64(0)

		for _, d := range validator.Delegators {
			amount, err := util.ParseInt64(d.StakeAmount)
			if err != nil {
				return err
			}

			delegationsCount++
			delegatedAmount += amount
			validatorAmount += amount
		}

		delegatedAmountMap[validator.NodeID] = validatorAmount
	}

	for _, validator := range payload.CurrentValidators {
		record, err := initValidator(&validator, payload.SyncTime)
		if err != nil {
			return err
		}

		record.FirstHeight = payload.Height
		record.LastHeight = payload.Height
		record.StakePercent = payload.ActiveValidatorShare[validator.NodeID]
		record.DelegationsPercent = util.PercentOf(int64(len(validator.Delegators)), int64(delegationsCount))
		record.DelegatedAmount = delegatedAmountMap[validator.NodeID]
		record.DelegatedAmountPercent = util.PercentOf(delegatedAmountMap[validator.NodeID], delegatedAmount)
		record.Capacity = record.StakeAmount*4 - record.DelegatedAmount
		record.CapacityPercent = util.PercentOf(record.DelegatedAmount, record.StakeAmount*4)

		seqRecord := model.ValidatorSeq{
			Time:                   payload.SyncTime,
			Height:                 payload.Height,
			NodeID:                 validator.NodeID,
			StakeAmount:            record.StakeAmount,
			StakePercent:           record.StakePercent,
			PotentialReward:        record.PotentialReward,
			RewardAddress:          record.RewardAddress,
			ActiveStartTime:        record.ActiveStartTime,
			ActiveEndTime:          record.ActiveEndTime,
			Active:                 true,
			ActiveProgressPercent:  record.ActiveProgressPercent,
			Uptime:                 record.Uptime,
			DelegationFee:          record.DelegationFee,
			DelegationsCount:       record.DelegationsCount,
			DelegationsPercent:     record.DelegationsPercent,
			DelegatedAmount:        record.DelegatedAmount,
			DelegatedAmountPercent: record.DelegatedAmountPercent,
		}

		payload.Validators = append(payload.Validators, *record)
		payload.ValidatorSeq = append(payload.ValidatorSeq, seqRecord)
	}

	return nil
}

// prepareDelegations builds a new set of delegation records
func (t ParserTask) prepareDelegations(payload *Payload) error {
	for _, validator := range payload.CurrentValidators {

		delegations, err := initDelegations(&validator, payload.SyncTime)
		if err != nil {
			return err
		}

		for idx := range delegations {
			delegations[idx].FirstHeight = payload.Height
			delegations[idx].LastHeight = payload.Height
			delegations[idx].CreatedAt = payload.SyncTime
			delegations[idx].UpdatedAt = payload.SyncTime
		}

		payload.Delegations = append(payload.Delegations, delegations...)
	}

	return nil
}
