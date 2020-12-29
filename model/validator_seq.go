package model

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model/types"
)

type ValidatorSeq struct {
	ID                     int
	Time                   time.Time
	Height                 int64
	NodeID                 string
	StakeAmount            types.Amount
	StakePercent           float64
	PotentialReward        types.Amount
	RewardAddress          string
	Active                 bool
	ActiveStartTime        time.Time
	ActiveEndTime          time.Time
	ActiveProgressPercent  float64
	DelegationsCount       int
	DelegationsPercent     float64
	DelegatedAmount        types.Amount
	DelegatedAmountPercent float64
	DelegationFee          float64
	Uptime                 float64
}

func (ValidatorSeq) TableName() string {
	return "validator_sequences"
}
