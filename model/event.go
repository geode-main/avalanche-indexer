package model

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Scope     string    `json:"scope"`
	Actor     string    `json:"actor"`
	ActorType string    `json:"actor_type"`
	Item      string    `json:"item"`
	ItemType  string    `json:"item_type"`
	Meta      Metadata  `json:"meta"`
	CreatedAt time.Time `json:"created_at"`
}

func (Event) TableName() string {
	return "events"
}

type Metadata string

func (m Metadata) MarshalJSON() ([]byte, error) {
	val := map[string]interface{}{}
	if err := json.Unmarshal([]byte(m), &val); err != nil {
		return nil, err
	}
	return json.Marshal(val)
}
