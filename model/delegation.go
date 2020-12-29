package model

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model/types"
)

type Delegation struct {
	ID              int          `json:"-"`
	ReferenceID     string       `json:"reference_id"`
	NodeID          string       `json:"node_id"`
	StakeAmount     types.Amount `json:"stake_amount"`
	PotentialReward types.Amount `json:"potential_reward"`
	RewardAddress   string       `json:"reward_address"`
	Active          bool         `json:"active"`
	ActiveStartTime time.Time    `json:"active_start_time"`
	ActiveEndTime   time.Time    `json:"active_end_time"`
	FirstHeight     int64        `json:"first_height"`
	LastHeight      int64        `json:"last_height"`
	CreatedAt       time.Time    `json:"-"`
	UpdatedAt       time.Time    `json:"-"`
}

func (Delegation) TableName() string {
	return "delegations"
}
