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
