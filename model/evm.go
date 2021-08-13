package model

import (
	"time"

	"github.com/lib/pq"
)

type EvmTrace struct {
	ID        string    `json:"id"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func (EvmTrace) TableName() string {
	return "evm_traces"
}

type EvmReceipt struct {
	ID              string `json:"id"`
	Type            int    `json:"type"`
	Status          int    `json:"status"`
	ContractAddress string `json:"contract_address"`
	Logs            string `json:"-"`
}

func (EvmReceipt) TableName() string {
	return "evm_receipts"
}

type EvmLog struct {
	Idx     int            `json:"index"`
	TxIdx   int            `json:"tx_index"`
	Address string         `json:"address"`
	Removed bool           `json:"removed"`
	Topics  pq.StringArray `gorm:"type:text[]" json:"topics"`
	Data    string         `json:"data"`
}
