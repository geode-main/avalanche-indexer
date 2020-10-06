package store

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
)

type ValidatorsStore struct {
	*gorm.DB
}

func (s ValidatorsStore) FindByNodeID(id string) (*model.Validator, error) {
	result := &model.Validator{}

	err := s.
		Model(result).
		First(result, "node_id = ?", id).
		Error

	return result, checkErr(err)
}

func (s ValidatorsStore) FindAll() ([]model.Validator, error) {
	result := []model.Validator{}

	err := s.
		Model(&model.Validator{}).
		Order("stake_amount").
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s ValidatorsStore) Import(records []model.Validator) error {
	if err := s.Exec("UPDATE validators SET active = FALSE").Error; err != nil {
		return err
	}

	return bulkImport(s.DB, queries.ValidatorsImport, len(records), func(i int) Row {
		r := records[i]

		return Row{
			r.NodeID,
			r.StakeAmount,
			r.StakePercent,
			r.PotentialReward,
			r.RewardAddress,
			r.Active,
			r.ActiveStartTime,
			r.ActiveEndTime,
			r.ActiveProgressPercent,
			r.Uptime,
			r.DelegationsCount,
			r.DelegationsPercent,
			r.DelegatedAmount,
			r.DelegatedAmountPercent,
			r.DelegationFee,
			r.Capacity,
			r.CapacityPercent,
			r.FirstHeight,
			r.LastHeight,
			r.CreatedAt,
			r.UpdatedAt,
		}
	})
}

func (s ValidatorsStore) ImportSeq(records []model.ValidatorSeq) error {
	return bulkImport(s.DB, queries.ValidatorSeqImport, len(records), func(i int) Row {
		r := records[i]

		return Row{
			r.Time,
			r.Height,
			r.NodeID,
			r.StakeAmount,
			r.StakePercent,
			r.PotentialReward,
			r.RewardAddress,
			r.Active,
			r.ActiveStartTime,
			r.ActiveEndTime,
			r.ActiveProgressPercent,
			r.DelegationsCount,
			r.DelegationsPercent,
			r.DelegatedAmount,
			r.DelegatedAmountPercent,
			r.DelegationFee,
			r.Uptime,
		}
	})
}

func (s ValidatorsStore) PurgeSeq(before time.Time) (int64, error) {
	result := s.Exec(queries.ValidatorSeqPurge, before)
	return result.RowsAffected, result.Error
}

func (s ValidatorsStore) CreateStats(t time.Time, bucket string) error {
	startTime, endTime := getTimeRange(t, bucket)
	query := prepareBucket(queries.ValidatorStatsCreate, bucket)

	return s.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(queries.ValidatorStatsDelete, startTime, bucket).Error; err != nil {
			return err
		}
		return tx.Exec(query, startTime, endTime).Error
	})
}

func (s ValidatorsStore) GetStats(id string, bucket string, limit int) ([]model.ValidatorStat, error) {
	result := []model.ValidatorStat{}

	err := s.
		Model(&model.ValidatorStat{}).
		Where("node_id = ? AND bucket = ?", id, bucket).
		Order("time ASC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}
