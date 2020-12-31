package store

import (
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
)

type DelegatorsStore struct {
	*gorm.DB
}

func (s DelegatorsStore) FindActive() ([]model.Delegation, error) {
	result := []model.Delegation{}

	err := s.
		Model(&model.Delegation{}).
		Where("active = ?", true).
		Order("stake_amount DESC").
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s DelegatorsStore) FindByRewardAddress(address string) ([]model.Delegation, error) {
	result := []model.Delegation{}

	err := s.
		Model(&model.Delegation{}).
		Where("reward_address = ?", address).
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s DelegatorsStore) FindByNodeID(id string) ([]model.Delegation, error) {
	result := []model.Delegation{}

	err := s.
		Model(&model.Delegation{}).
		Where("node_id = ? AND active = ?", id, true).
		Find(&result).
		Error

	return result, checkErr(err)
}

func (s DelegatorsStore) Import(records []model.Delegation) error {
	if err := s.Exec("UPDATE delegations SET active = FALSE").Error; err != nil {
		return err
	}

	return bulkImport(s.DB, queries.DelegatorsImport, len(records), func(idx int) Row {
		r := records[idx]

		return Row{
			r.ReferenceID,
			r.NodeID,
			r.StakeAmount,
			r.PotentialReward,
			r.RewardAddress,
			r.Active,
			r.ActiveStartTime,
			r.ActiveEndTime,
			r.FirstHeight,
			r.LastHeight,
			r.CreatedAt,
			r.UpdatedAt,
		}
	})
}
