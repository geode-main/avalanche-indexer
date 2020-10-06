package model

import "time"

type ValidatorSeq struct {
	ID                     int
	Time                   time.Time
	Height                 int64
	NodeID                 string
	StakeAmount            int64
	StakePercent           float64
	PotentialReward        int64
	RewardAddress          string
	Active                 bool
	ActiveStartTime        time.Time
	ActiveEndTime          time.Time
	ActiveProgressPercent  float64
	DelegationsCount       int
	DelegationsPercent     float64
	DelegatedAmount        int64
	DelegatedAmountPercent float64
	DelegationFee          float64
	Uptime                 float64
}

func (ValidatorSeq) TableName() string {
	return "validator_sequences"
}
