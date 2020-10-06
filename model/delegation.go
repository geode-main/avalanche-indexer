package model

import "time"

type Delegation struct {
	ID              string    `json:"id"`
	ReferenceID     string    `json:"reference_id"`
	NodeID          string    `json:"node_id"`
	StakeAmount     int64     `json:"stake_amount"`
	PotentialReward int64     `json:"potential_reward"`
	RewardAddress   string    `json:"reward_address"`
	Active          bool      `json:"active"`
	ActiveStartTime time.Time `json:"active_start_time"`
	ActiveEndTime   time.Time `json:"active_end_time"`
	FirstHeight     int64     `json:"first_height"`
	LastHeight      int64     `json:"last_height"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}

func (Delegation) TableName() string {
	return "delegations"
}
