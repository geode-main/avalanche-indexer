package model

import "time"

type SyncStatus struct {
	ID        string    `json:"id"`
	IndexID   int64     `json:"index_id"`
	IndexTime time.Time `json:"index_time"`
	TipID     int64     `json:"tip_id"`
	TipTime   time.Time `json:"tip_time"`
}

func (SyncStatus) TableName() string {
	return "sync_statuses"
}

func (s SyncStatus) AtTip() bool {
	return s.IndexID >= s.TipID
}

func (s SyncStatus) NextID() int64 {
	return s.IndexID + 1
}

func (s SyncStatus) Lag() int64 {
	return s.TipID - s.IndexID
}
