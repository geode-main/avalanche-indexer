package model

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model/types"
	"github.com/figment-networks/avalanche-indexer/util"
)

type Event struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Scope       string    `json:"scope,omitempty"`
	Chain       string    `json:"chain,omitempty"`
	BlockHash   string    `json:"block"`
	BlockHeight uint64    `json:"block_height"`
	TxHash      string    `json:"tx_hash"`
	ItemID      string    `json:"item_id"`
	ItemType    string    `json:"item_type"`
	Timestamp   time.Time `json:"timestamp"`
	Data        types.Map `json:"data,omitempty" gorm:"type:text"`
}

func (Event) TableName() string {
	return "events"
}

func (e *Event) AssignID() error {
	idStr := e.Scope + e.Type + e.TxHash + e.ItemID + e.ItemType
	eventID, err := util.AvalancheIDFromString(idStr)
	if err != nil {
		return err
	}
	e.ID = eventID
	return nil
}
