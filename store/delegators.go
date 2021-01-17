package store

import (
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
)

type DelegatorsStore struct {
	*gorm.DB
}

type DelegationsSearch struct {
	NodeID        string `form:"node_id"`
	RewardAddress string `form:"reward_address"`
}

// Search performs a seach on delegations
func (s DelegatorsStore) Search(search DelegationsSearch) ([]model.Delegation, error) {
	result := []model.Delegation{}

	scope := s.
		Model(&model.Delegation{}).
		Where("active = ?", true).
		Order("id DESC")

	if search.NodeID != "" {
		scope = scope.Where("node_id = ?", search.NodeID)
	}
	if search.RewardAddress != "" {
		scope = scope.Where("reward_address = ?", search.RewardAddress)
	}

	err := scope.Find(&result).Error
	return result, checkErr(err)
}

// Import imports delegations records in bulk
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
