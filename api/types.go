package api

import (
	"time"

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
