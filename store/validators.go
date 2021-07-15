package store

import (
	"errors"
	"time"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
)

type ValidatorsStore struct {
	*gorm.DB
}

func (s ValidatorsStore) LastHeight() (int64, error) {
	rows, err := s.DB.Raw(queries.ValidatorLastHeight).Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}

	var result int64
	return result, rows.Scan(&result)
}

func (s ValidatorsStore) LastTime() (*time.Time, error) {
	rows, err := s.DB.Raw(queries.ValidatorSeqLastTime).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var result time.Time
	return &result, rows.Scan(&result)
}

func (s ValidatorsStore) FindByNodeID(id string) (*model.Validator, error) {
	result := &model.Validator{}

	err := s.
		Model(result).
		Where("node_id = ?", id).
		Take(result).
		Error

	return result, checkErr(err)
}

func (s ValidatorsStore) Search(search ValidatorsSearch) ([]model.Validator, error) {
	result := []model.Validator{}

	scope := s.
		Model(&model.Validator{}).
		Order("stake_amount").
		Where("active = ?", true)

	if search.RewardAddress != "" {
		scope = scope.Where("reward_address = ?", search.RewardAddress)
	}
	if search.CapacityPercentMin > 0 {
		scope = scope.Where("capacity_percent >= ?", search.CapacityPercentMin-1)
	}
	if search.CapacityPercentMax > 0 {
		scope = scope.Where("capacity_percent <= ?", search.CapacityPercentMax)
	}

	err := scope.Find(&result).Error

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
		Order("time DESC").
		Limit(limit).
		Find(&result).
		Error

	return result, checkErr(err)
}

type ValidatorsSearch struct {
	RewardAddress      string `form:"reward_address"`
	CapacityPercentMin uint   `form:"capacity_percent_min"`
	CapacityPercentMax uint   `form:"capacity_percent_max"`
}

func (s ValidatorsSearch) Validate() error {
	if s.CapacityPercentMin > 100 {
		return errors.New("capacity_percent_min must be below 100")
	}
	if s.CapacityPercentMax > 100 {
		return errors.New("capacity_percent_max must be below 100")
	}
	return nil
}
