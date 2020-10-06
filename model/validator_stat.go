package model

import "time"

type ValidatorStat struct {
	ID                     int       `json:"-"`
	Time                   time.Time `json:"time"`
	Bucket                 string    `json:"bucket"`
	NodeID                 string    `json:"-"`
	UptimeMin              float64   `json:"uptime_min"`
	UptimeMax              float64   `json:"uptime_max"`
	UptimeAvg              float64   `json:"uptime_avg"`
	StakeAmount            int64     `json:"stake_amount"`
	StakePercent           float64   `json:"stake_percent"`
	DelegationsCount       int       `json:"delegations_count"`
	DelegationsPercent     float64   `json:"delegations_percent"`
	DelegatedAmount        int64     `json:"delegated_amount"`
	DelegatedAmountPercent float64   `json:"delegated_amount_percent"`
}

func (ValidatorStat) TableName() string {
	return "validator_stats"
}
