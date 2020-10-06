package model

import "time"

type DelegatorSeq struct {
	ID              int
	NodeID          string
	StakeAmount     int64
	PotentialReward int64
	ActiveStartTime time.Time
	ActiveEndTime   time.Time
	CreatedAt       time.Time
}

func (DelegatorSeq) TableName() string {
	return "delegator_sequences"
}
