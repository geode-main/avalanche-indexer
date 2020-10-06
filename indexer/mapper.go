package indexer

import (
	"fmt"
	"time"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/util"
)

// initValidator builds a new validator record from the raw client data
func initValidator(validator *client.Validator, ts time.Time) (*model.Validator, error) {
	stake, err := util.ParseInt64(validator.StakeAmount)
	if err != nil {
		return nil, err
	}

	reward, err := util.ParseInt64(validator.PotentialReward)
	if err != nil {
		return nil, err
	}

	startTime, err := util.ParseUnixTime(validator.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := util.ParseUnixTime(validator.EndTime)
	if err != nil {
		return nil, err
	}

	delegationFee, err := util.ParseFloat32(validator.DelegationFee)
	if err != nil {
		return nil, err
	}

	uptime, err := util.ParseFloat32(validator.Uptime)
	if err != nil {
		return nil, err
	}

	// calculate the current validation progress
	timeLeft := endTime.Sub(ts).Milliseconds()
	timeTotal := endTime.Sub(startTime).Milliseconds()
	progressPercent := util.PercentOf(timeLeft, timeTotal)

	return &model.Validator{
		NodeID:                 validator.NodeID,
		StakeAmount:            stake,
		StakePercent:           0, // filled later in the pipeline
		PotentialReward:        reward,
		RewardAddress:          validator.RewardOwner.Addresses[0],
		ActiveStartTime:        startTime,
		ActiveEndTime:          endTime,
		Active:                 true,
		ActiveProgressPercent:  progressPercent,
		DelegationFee:          delegationFee,
		DelegationsCount:       len(validator.Delegators),
		DelegationsPercent:     0,            // filled later in the pipeline
		DelegatedAmount:        0,            // filled later in the pipeline
		DelegatedAmountPercent: 0,            // filled later in the pipeline
		Uptime:                 uptime * 100, // we want this in %
	}, nil
}

func initDelegation(delegator *client.Delegator) (result model.Delegation, err error) {
	amount, err := util.ParseInt64(delegator.StakeAmount)
	if err != nil {
		return result, err
	}

	reward, err := util.ParseInt64(delegator.PotentialReward)
	if err != nil {
		return result, err
	}

	startTime, err := util.ParseUnixTime(delegator.StartTime)
	if err != nil {
		return result, err
	}

	endTime, err := util.ParseUnixTime(delegator.EndTime)
	if err != nil {
		return result, err
	}

	// note: avax does not expose internal ID for delegations so we must build one
	referenceID := fmt.Sprintf(
		"%s_%s_%v_%v",
		delegator.NodeID,
		delegator.RewardOwner.Addresses[0],
		startTime.Unix(),
		endTime.Unix(),
	)

	return model.Delegation{
		ReferenceID:     util.StringSHA1(referenceID),
		NodeID:          delegator.NodeID,
		StakeAmount:     amount,
		PotentialReward: reward,
		RewardAddress:   delegator.RewardOwner.Addresses[0],
		Active:          true,
		ActiveStartTime: startTime,
		ActiveEndTime:   endTime,
	}, nil
}

func initDelegations(validator *client.Validator, ts time.Time) ([]model.Delegation, error) {
	result := []model.Delegation{}

	for _, d := range validator.Delegators {
		delegator, err := initDelegation(&d)
		if err != nil {
			return nil, err
		}
		delegator.CreatedAt = ts

		result = append(result, delegator)
	}

	return result, nil
}
