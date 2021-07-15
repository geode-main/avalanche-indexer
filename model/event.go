package model

import (
	"time"

	"github.com/figment-networks/avalanche-indexer/model/types"
)

type Event struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	Chain       string    `json:"chain"`
	Block       string    `json:"block"`
	Transaction string    `json:"transaction"`
	ItemID      string    `json:"item_id"`
	ItemType    string    `json:"item_type"`
	Meta        types.Map `json:"meta"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Event) TableName() string {
	return "events"
}
