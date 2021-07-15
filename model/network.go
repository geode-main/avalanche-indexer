package model

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model/types"
)

type NetworkMetric struct {
	ID                      int          `json:"-"`
	Time                    time.Time    `json:"time"`
	Height                  int64        `json:"height"`
	PeersCount              int          `json:"peers_count"`
	BlockchainsCount        int          `json:"blockchains_count"`
	ActiveValidatorsCount   int          `json:"active_validators_count"`
	PendingValidatorsCount  int          `json:"pending_validators_count"`
	ActiveDelegationsCount  int          `json:"active_delegations_count"`
	PendingDelegationsCount int          `json:"pending_delegations_count"`
	MinValidatorStake       int64        `json:"min_validator_stake"`
	MinDelegationStake      int64        `json:"min_delegation_stake"`
	TxFee                   int          `json:"tx_fee"`
	CreationTxFee           int          `json:"creation_tx_fee"`
	Uptime                  float64      `json:"uptime"`
	DelegationFee           float64      `json:"delegation_fee"`
	TotalStaked             types.Amount `json:"total_staked"`
	TotalDelegated          types.Amount `json:"total_delegated"`
}

func (NetworkMetric) TableName() string {
	return "network_metrics"
}

type NetworkStat struct {
	ID                 int          `json:"-"`
	Time               time.Time    `json:"time"`
	Bucket             string       `json:"bucket"`
	HeightChange       int          `json:"height_change"`
	Peers              int          `json:"peers"`
	Blockchains        int          `json:"blockchains"`
	ActiveValidators   int          `json:"active_validators"`
	PendingValidators  int          `json:"pending_validators"`
	ValidatorUptime    float64      `json:"validator_uptime"`
	ActiveDelegations  int          `json:"active_delegations"`
	PendingDelegations int          `json:"pending_delegations"`
	MinValidatorStake  int64        `json:"min_validator_stake"`
	MinDelegationStake int64        `json:"min_delegator_stake"`
	TxFee              int64        `json:"tx_fee"`
	CreateTxFee        int64        `json:"create_tx_fee"`
	TotalStaked        types.Amount `json:"total_staked"`
	TotalDelegated     types.Amount `json:"total_delegated"`
}

func (NetworkStat) TableName() string {
	return "network_stats"
}
