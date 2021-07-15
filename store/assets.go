package store

import (
	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AssetsStore struct {
	*gorm.DB
}

func (s AssetsStore) GetAll() ([]model.Asset, error) {
	result := []model.Asset{}
	err := s.Model(&model.Asset{}).Order("name ASC").Find(&result).Error
	return result, err
}

func (s AssetsStore) GetByType(assetType string) ([]model.Asset, error) {
	result := []model.Asset{}
	err := s.Model(&model.Asset{}).Where("type = ?", assetType).Order("name ASC").Find(&result).Error
	return result, checkErr(err)
}

func (s AssetsStore) Get(assetID string) (*model.Asset, error) {
	asset := &model.Asset{}
	err := s.Model(asset).First(asset, "asset_id = ?", assetID).Error
	return asset, checkErr(err)
}

func (s AssetsStore) GetTransactionsCount(assetID string) (*int, error) {
	rows, err := s.Raw(queries.PlatformAssetTransactionsCount, assetID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var num int
	return &num, rows.Scan(&num)
}

func (s AssetsStore) Create(asset *model.Asset) error {
	return s.
		Model(asset).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(asset).
		Error
}
