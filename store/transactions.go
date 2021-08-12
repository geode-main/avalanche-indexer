package store

import (
	"regexp"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/figment-networks/avalanche-indexer/model"
)

var (
	reDate     = regexp.MustCompile(`^[\d]{4}-[\d]{2}-[\d]{2}$`)
	reUnixTime = regexp.MustCompile(`^[\d]{10}$`)
)

type TransactionsStore struct {
	*gorm.DB
}

// Search returns transactions matching search input
func (store TransactionsStore) Search(input TxSearchInput) (*TxSearchOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	scope := store.
		Model(&model.Transaction{}).
		Order("transactions.timestamp DESC")

	if input.startTime != nil {
		scope = scope.Where("transactions.timestamp >= ?", input.startTime)
	}
	if input.endTime != nil {
		scope = scope.Where("transactions.timestamp <= ?", input.endTime)
	}

	if input.StartHeight > 0 {
		scope = scope.Where("transactions.block_height >= ?", input.StartHeight)
	}

	if input.Chain != "" {
		scope = scope.Where("transactions.chain = ?", input.Chain)
	}

	if len(input.types) > 0 {
		scope = scope.Where("transactions.type IN (?)", input.types)
	}

	if input.Memo != "" {
		words := strings.Split(strings.TrimSpace(input.Memo), " ")
		for _, word := range words {
			scope = scope.Where("memo_tsv @@ to_tsquery('english', ?::text)", word)
		}
	}

	if input.BlockHash != "" {
		scope = scope.Where("block = ?", input.BlockHash)
	}

	if input.Address != "" || input.Asset != "" {
		scope = scope.
			Select("transactions.*").
			Joins("INNER JOIN transaction_outputs ON transaction_outputs.tx_id = transactions.id")

		if input.Asset != "" {
			scope = scope.Where("transaction_outputs.asset = ?", input.Asset)
		}

		if input.Address != "" {
			addresses := pq.StringArray(strings.Split(input.Address, ","))
			scope = scope.Where("transaction_outputs.addresses && ?", addresses)
		}
	}

	if input.BeforeID != "" {
		beforeTx, err := store.GetShortByID(input.BeforeID)
		if err != nil {
			return nil, err
		}
		scope = scope.Where("transactions.timestamp < ?", beforeTx.Timestamp)
	}

	if input.AfterID != "" {
		afterTx, err := store.GetShortByID(input.AfterID)
		if err != nil {
			return nil, err
		}
		scope = scope.Where("transactions.timestamp > ?", afterTx.Timestamp)
	}

	transactions := []model.Transaction{}

	err := scope.
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&transactions).
		Error

	if err != nil {
		return nil, err
	}

	txIDs := make([]string, len(transactions))
	for idx, tx := range transactions {
		txIDs[idx] = tx.ID
	}

	inputs := []model.Output{}
	outputs := []model.Output{}

	if len(txIDs) > 0 {
		if err := store.Model(&model.Output{}).Where("spent_tx_id IN (?)", txIDs).Find(&inputs).Error; err != nil {
			return nil, err
		}

		if err := store.Model(&model.Output{}).Where("tx_id IN (?)", txIDs).Find(&outputs).Error; err != nil {
			return nil, err
		}
	}

	for idx, tx := range transactions {
		for _, input := range inputs {
			if *input.SpentTxID == tx.ID {
				transactions[idx].Inputs = append(transactions[idx].Inputs, input)
			}
		}

		for _, output := range outputs {
			if output.TxID == tx.ID {
				transactions[idx].Outputs = append(transactions[idx].Outputs, output)
			}
		}

		transactions[idx].UpdateAmounts()
	}

	return &TxSearchOutput{Transactions: transactions}, nil
}

// GetShortByID returns just the transaction record without any extra data
func (store TransactionsStore) GetShortByID(id string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	err := store.Model(tx).Take(tx, "id = ?", id).Error
	return tx, err
}

// GetByID returns transaction by ID
func (store TransactionsStore) GetByID(id string) (*model.Transaction, error) {
	tx := &model.Transaction{}

	err := store.Model(tx).Take(tx, "id = ?", id).Error
	if err != nil {
		return nil, checkErr(err)
	}

	if tx.UsesUTXOs() {
		err = store.Model(&model.Output{}).Where("spent_tx_id = ?", id).Order("index ASC").Find(&tx.Inputs).Error
		if err != nil {
			return nil, checkErr(err)
		}

		err = store.Model(&model.Output{}).Where("tx_id = ?", id).Order("index ASC").Find(&tx.Outputs).Error
		if err != nil {
			return nil, checkErr(err)
		}
	}

	tx.UpdateAmounts()

	return tx, nil
}

// GetTypeCounts returns transaction types with counts
func (s TransactionsStore) GetTypeCounts(chain string) ([]model.TransactionTypeCount, error) {
	result := []model.TransactionTypeCount{}

	scope := s.Model(&model.Transaction{}).Select("type, COUNT(1) AS total_count")
	if chain != "" {
		scope = scope.Where("chain = ?", chain)
	}

	err := scope.
		Group("type").
		Order("total_count DESC").
		Find(&result).
		Error

	return result, err
}
