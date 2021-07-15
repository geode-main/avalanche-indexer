package model

import (
	"encoding/base64"
	"time"
)

type RawMessageTopic struct {
	ID    int    `json:"id"`
	Chain string `json:"chain"`
	Name  string `json:"name"`
	VM    string `json:"vm"`
}

type RawMessage struct {
	ID          int        `json:"id"`
	TopicID     int        `json:"topic_id"`
	IndexID     int        `json:"index_id"`
	Data        string     `json:"data"`
	Hash        string     `json:"hash"`
	CreatedAt   time.Time  `json:"created_at"`
	ProcessedAt *time.Time `json:"-"`
}

func (msg RawMessage) DataBytes() ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		result, err = base64.RawStdEncoding.DecodeString(msg.Data)
	}
	return result, err
}
