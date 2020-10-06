package model

import "time"

type Validator struct {
	ID                     int       `json:"-"`
	NodeID                 string    `json:"node_id"`
	StakeAmount            int64     `json:"stake_amount"`
	StakePercent           float64   `json:"stake_percent"`
	PotentialReward        int64     `json:"potential_reward"`
	RewardAddress          string    `json:"reward_address"`
	Active                 bool      `json:"active"`
	ActiveStartTime        time.Time `json:"active_start_time"`
	ActiveEndTime          time.Time `json:"active_end_time"`
	ActiveProgressPercent  float64   `json:"active_progress_percent"`
	Uptime                 float64   `json:"uptime"`
	DelegationsCount       int       `json:"delegations_count"`
	DelegationsPercent     float64   `json:"delegations_percent"`
	DelegatedAmount        int64     `json:"delegated_amount"`
	DelegatedAmountPercent float64   `json:"delegated_amount_percent"`
	DelegationFee          float64   `json:"delegation_fee"`
	Capacity               int64     `json:"capacity"`
	CapacityPercent        float64   `json:"capacity_percent"`
	FirstHeight            int64     `json:"first_height"`
	LastHeight             int64     `json:"last_height"`
	CreatedAt              time.Time `json:"-"`
	UpdatedAt              time.Time `json:"-"`
}

func (Validator) TableName() string {
	return "validators"
}
