package model

import "time"

type Reward struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	Rewarded      bool      `json:"rewarded"`
	RewardedAt    time.Time `json:"rewarded_at"`
	ProcessedAt   time.Time `json:"processed_at"`
}

type RewardsOwner struct {
	ID        string                `json:"id"`
	Locktime  uint64                `json:"locktime"`
	Threshold uint32                `json:"threshold"`
	Addresses []RewardsOwnerAddress `json:"-" gorm:"-"`
	Outputs   []RewardsOwnerOutput  `json:"-" gorm:"-"`
}

type RewardsOwnerAddress struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Index   uint32 `json:"index"`
}

type RewardsOwnerOutput struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Index         uint32 `json:"index"`
}
