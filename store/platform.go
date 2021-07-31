package store

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
)

const (
	outputsImportBatchSize = 100
)

type PlatformStore struct {
	*gorm.DB
}

func NewPlatformStore(db *gorm.DB) PlatformStore {
	return PlatformStore{
		DB: db,
	}
}

// Chains returns all existing chain records
func (s *PlatformStore) Chains() ([]model.Chain, error) {
	result := []model.Chain{}
	err := s.Model(&model.Chain{}).Find(&result).Error
	return result, err
}

// GetChain returns a single chain record by ID
func (s *PlatformStore) GetChain(chainID string) (*model.Chain, error) {
	chain := &model.Chain{}
	err := s.Model(chain).First(chain, "chain_id = ?", chainID).Error
	return chain, checkErr(err)
}

// GetChain returns a single chain record by name
func (s *PlatformStore) GetChainByName(name string) (*model.Chain, error) {
	chain := &model.Chain{}
	err := s.Model(chain).First(chain, "name = ?", name).Error
	return chain, checkErr(err)
}

// CreateChain creates a new chain record
func (s *PlatformStore) CreateChain(chain *model.Chain) error {
	err := s.
		Model(chain).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(chain).
		Error

	return err
}

// GetBlock returns a single block by hash
func (s *PlatformStore) GetBlock(hash string) (*model.Block, error) {
	block := &model.Block{}
	err := s.Model(block).First(block, "id = ?", hash).Error
	return block, checkErr(err)
}

// GetBlocks returns blocks matching the search query
func (s *PlatformStore) GetBlocks(search *BlocksSearch) ([]model.Block, error) {
	if err := search.Validate(); err != nil {
		return nil, err
	}

	scope := s.
		Model(&model.Block{}).
		Where("chain = ?", search.Chain)

	if search.StartHeight > 0 {
		scope = scope.Where("height >= ?", search.StartHeight)
	}

	if search.EndHeight > 0 {
		scope = scope.Where("height <= ?", search.EndHeight)
	}

	if search.Type != "" {
		types := strings.Split(search.Type, ",")
		scope = scope.Where("type IN (?)", types)
	}

	switch search.Order {
	case "height_asc":
		scope = scope.Order("height ASC")
	case "height_desc":
		scope = scope.Order("height DESC")
	}

	result := []model.Block{}
	err := scope.Limit(search.Limit).Find(&result).Error

	return result, err
}

// LastBlock returns a block for the latest height
func (s *PlatformStore) LastBlock(chain string) (*model.Block, error) {
	result := &model.Block{}

	err := s.
		Model(result).
		Where("chain = ?", chain).
		Order("height DESC").
		Take(result).
		Error

	return result, checkErr(err)
}

// CreateBlock creates a new block
func (s *PlatformStore) CreateBlock(block *model.Block) error {
	return s.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(block).Error
}

// CreateTransaction creates a new transaction
func (s *PlatformStore) CreateTransaction(tx *model.Transaction) error {
	return s.Exec(queries.PlatformCreateTransaction,
		tx.ID,
		tx.ReferenceTxID,
		tx.Status,
		tx.Type,
		tx.Block,
		tx.BlockHeight,
		tx.Chain,
		tx.Memo,
		tx.MemoText,
		tx.Fee,
		tx.Nonce,
		tx.SourceChain,
		tx.DestinationChain,
		tx.Timestamp,
		tx.Metadata,
	).Error
}

// CreateTxInputs creates transaction input records
func (s *PlatformStore) CreateTxInputs(ids []string, txID string) error {
	return bulkImport(s.DB, queries.PlatformCreateInputs, len(ids), func(i int) Row {
		return Row{
			ids[i],
			txID,
		}
	})
}

// CreateTxOutputs create transaction output records
func (s *PlatformStore) CreateTxOutputs(outputs []model.Output) error {
	n := len(outputs)

	for idx := 0; idx < n; idx += outputsImportBatchSize {
		endIdx := idx + outputsImportBatchSize
		if endIdx > n {
			endIdx = n
		}

		batch := outputs[idx:endIdx]

		err := bulkImport(s.DB, queries.PlatformCreateOutputs, len(batch), func(i int) Row {
			r := batch[i]

			return Row{
				r.ID,
				r.TxID,
				r.Chain,
				r.Asset,
				r.Type,
				r.Index,
				r.Locktime,
				r.Threshold,
				r.Amount,
				r.Group,
				r.Stake,
				r.Reward,
				r.Spent,
				r.SpentTxID,
				r.Addresses,
				r.Payload,
			}
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetTransactionOutput returns a single output record
func (s *PlatformStore) GetTransactionOutput(id string) (*model.Output, error) {
	result := &model.Output{}
	err := s.Model(result).First(result, "id = ?", id).Error
	return result, checkErr(err)
}

// MarkOutputsSpent updates the spent transaction reference on outputs
func (s *PlatformStore) MarkOutputsSpent(ids []string, txID string, txTime time.Time) error {
	return s.Exec(queries.PlatformMarkOutputsSpent, txID, ids).Error
}

// CreateRewardsOwner creates rewards owner records
func (s *PlatformStore) CreateRewardsOwner(owner *model.RewardsOwner) error {
	err := s.Exec(queries.PlatformCreateRewardOwners,
		owner.ID,
		owner.Locktime,
		owner.Threshold,
	).Error

	if err != nil {
		return err
	}

	for _, addr := range owner.Addresses {
		if err := s.Exec(queries.PlatformCreateRewardAddresses,
			addr.ID,
			addr.Address,
			addr.Index,
		).Error; err != nil {
			return err
		}
	}

	for _, output := range owner.Outputs {
		if err := s.Exec(queries.PlatformCreateRewardsOwnerOutputs,
			output.ID,
			output.TransactionID,
			output.Index,
		).Error; err != nil {
			return err
		}
	}

	return nil
}

// UpdateSyncStatus creates or updates an sync status record
func (s *PlatformStore) UpdateSyncStatus(status *model.SyncStatus) error {
	return s.Exec(
		queries.PlatformUpdateChainIndexStatus,
		status.ID,
		status.IndexID,
		status.IndexTime,
		status.TipID,
		status.TipTime,
	).Error
}

// GetSyncStatus returns the chain sync status record
func (s *PlatformStore) GetSyncStatus(chain string) (*model.SyncStatus, error) {
	result := &model.SyncStatus{}
	err := s.Model(result).First(result, "id = ?", chain).Error
	return result, checkErr(err)
}

// GetSyncStatuses returns all sync status records
func (s *PlatformStore) GetSyncStatuses() ([]model.SyncStatus, error) {
	result := []model.SyncStatus{}
	err := s.Model(&model.SyncStatus{}).Order("id ASC").Find(&result).Error
	return result, err
}
