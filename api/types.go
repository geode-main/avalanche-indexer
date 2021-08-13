package api

import (
	"math/big"
	"time"

	"github.com/figment-networks/avalanche-indexer/client"
	"github.com/figment-networks/avalanche-indexer/model"
)

type StatusResponse struct {
	AppName     string     `json:"app_name"`
	AppVersion  string     `json:"app_version"`
	GitCommit   string     `json:"git_commit"`
	GoVersion   string     `json:"go_version"`
	SyncStatus  string     `json:"sync_status"`
	SyncTime    *time.Time `json:"sync_time"`
	NodeVersion string     `json:"node_version"`
	NetworkName string     `json:"network_name"`
}

type ValidatorResponse struct {
	Validator   *model.Validator      `json:"validator"`
	Delegations []model.Delegation    `json:"delegations"`
	HourlyStats []model.ValidatorStat `json:"stats_24h"`
	DailyStats  []model.ValidatorStat `json:"stats_30d"`
}

type AddressBalancesResponse struct {
	Platform client.Balance      `json:"P"`
	Exchance []client.AvmBalance `json:"X"`
}

type CBalanceResponse struct {
	Balance string   `json:"balance"`
	Height  *big.Int `json:"height"`
}

type TxTraceResponse struct {
	Receipt *model.EvmReceipt `json:"receipt"`
	Logs    []model.EvmLog    `json:"logs"`
	Trace   *client.Call      `json:"trace"`
}
