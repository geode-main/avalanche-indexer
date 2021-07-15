package model

import (
	"github.com/lib/pq"
)

type Output struct {
	ID        string         `json:"id"`
	TxID      string         `json:"tx_id"`
	Chain     string         `json:"chain"`
	Asset     string         `json:"asset"`
	Type      string         `json:"type"`
	Index     uint64         `json:"index"`
	Locktime  uint64         `json:"locktime"`
	Threshold uint32         `json:"threshold"`
	Amount    uint64         `json:"amount"`
	Group     uint32         `json:"group"`
	Addresses pq.StringArray `json:"addresses" gorm:"type:text[]"`
	Stake     bool           `json:"stake"`
	Reward    bool           `json:"reward"`
	Spent     bool           `json:"spent"`
	SpentTxID *string        `json:"spent_in_tx"`
	Payload   *string        `json:"payload,omitempty"`
}

func (Output) TableName() string {
	return "transaction_outputs"
}
