package model

import "time"

type NetworkStat struct {
	ID                 int       `json:"-"`
	Time               time.Time `json:"time"`
	Bucket             string    `json:"bucket"`
	HeightChange       int       `json:"height_change"`
	Peers              int       `json:"peers"`
	Blockchains        int       `json:"blockchains"`
	ActiveValidators   int       `json:"active_validators"`
	PendingValidators  int       `json:"pending_validators"`
	ValidatorUptime    float64   `json:"validator_uptime"`
	ActiveDelegations  int       `json:"active_delegations"`
	PendingDelegations int       `json:"pending_delegations"`
	MinValidatorStake  int64     `json:"min_validator_stake"`
	MinDelegationStake int64     `json:"min_delegator_stake"`
	TxFee              int64     `json:"tx_fee"`
	CreateTxFee        int64     `json:"create_tx_fee"`
}

func (NetworkStat) TableName() string {
	return "network_stats"
}
