package store

import (
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
)

type AddressesStore struct {
	*gorm.DB
}

func (s AddressesStore) Import(records []model.Address) error {
	return bulkImport(s.DB, queries.AddressesImport, len(records), func(i int) Row {
		return Row{
			records[i].Value,
			records[i].Balance,
			records[i].UnlockedBalance,
			records[i].LockedNotStakeable,
			records[i].CreatedAt,
			records[i].UpdatedAt,
		}
	})
}
