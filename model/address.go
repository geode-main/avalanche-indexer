package model

import "time"

type Address struct {
	ID                 int       `json:"id"`
	Chain              string    `json:"chain"`
	Value              string    `json:"value"`
	Balance            uint64    `json:"balance"`
	UnlockedBalance    uint64    `json:"unlocked_balance"`
	LockedStakeable    uint64    `json:"locked_stakeable"`
	LockedNotStakeable uint64    `json:"locked_not_stakeable"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}
