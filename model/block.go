package model

import (
	"time"
)

type Block struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Parent    string    `json:"parent"`
	Chain     string    `json:"chain"`
	Height    uint64    `json:"height"`
	Timestamp time.Time `json:"timestamp"`
}

func (Block) TableName() string {
	return "blocks"
}
