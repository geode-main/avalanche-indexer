package store

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
)

type NetworksStore struct {
	*gorm.DB
}

// CreateMetric creates a new network metric record
func (s NetworksStore) CreateMetric(record *model.NetworkMetric) error {
	err := s.Model(record).Create(record).Error
	return checkErr(err)
}

// CreateStats creates a stat for a given time bucket
func (s NetworksStore) CreateStats(t time.Time, bucket string) error {
	startTime, endTime := getTimeRange(t, bucket)
	query := prepareBucket(queries.NetworkStatsCreate, bucket)

	return s.Exec(query, startTime, endTime).Error
}

// GetStats returns a set of stats for a given time bucket
func (s NetworksStore) GetStats(bucket string, limit int) ([]model.NetworkStat, error) {
	result := []model.NetworkStat{}

	err := s.
		Model(&model.NetworkStat{}).
		Where("bucket = ?", bucket).
		Order("time ASC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}
