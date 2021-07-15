package store

import (
	"encoding/base64"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/figment-networks/avalanche-indexer/model"
	"github.com/figment-networks/avalanche-indexer/store/queries"
	"github.com/figment-networks/avalanche-indexer/util"
)

type RawMessagesStore struct {
	*gorm.DB
}

func (s RawMessagesStore) GetTopics() ([]model.RawMessageTopic, error) {
	topics := []model.RawMessageTopic{}
	err := s.Model(&model.RawMessageTopic{}).Find(&topics).Error
	return topics, err
}

func (s RawMessagesStore) GetTopicByChain(chain string) (*model.RawMessageTopic, error) {
	topic := &model.RawMessageTopic{}
	err := s.Find(topic, "chain = ?", chain).Error
	return topic, checkErr(err)
}

func (s RawMessagesStore) CreateOrFindTopic(chain, name, vm string) (*model.RawMessageTopic, error) {
	topic := &model.RawMessageTopic{
		Chain: chain,
		Name:  name,
		VM:    vm,
	}

	err := s.
		Model(topic).
		FirstOrCreate(topic, &model.RawMessageTopic{Chain: chain}).
		Error

	return topic, err
}

func (s RawMessagesStore) GetMessage(id int) (*model.RawMessage, error) {
	result := &model.RawMessage{}
	err := s.Model(result).First(result, "id = ?", id).Error
	return result, checkErr(err)
}

func (s RawMessagesStore) CreateMessage(topic int, data []byte, timestamp time.Time) error {
	encoded := base64.StdEncoding.EncodeToString(data)

	hash, err := util.AvalancheID(data)
	if err != nil {
		return err
	}

	return s.Exec(queries.RawMessagesCreate, topic, encoded, hash, timestamp).Error
}

type GetRawMessagesInput struct {
	Chain     string `form:"chain"`
	Topic     int    `form:"topic"`
	StartID   int    `form:"start_id"`
	Limit     int    `form:"limit"`
	Processed *bool  `form:"-"`
}

func (s RawMessagesStore) LastProcessedID(topic int) (int, error) {
	msg := &model.RawMessage{}

	err := s.
		Model(msg).
		Where("topic_id = ?", topic).
		Where("processed_at IS NULL").
		Order("id ASC").
		First(msg).
		Error

	if err != nil {
		if err == ErrNotFound {
			return 0, nil
		}
		return -1, err
	}

	return msg.ID, nil
}

func (s RawMessagesStore) GetMessages(input GetRawMessagesInput) ([]model.RawMessage, error) {
	if input.Limit == 0 {
		input.Limit = 100
	}
	if input.Limit > 1000 {
		return nil, errors.New("max limit is 1000")
	}

	if input.Topic == 0 {
		if input.Chain == "" {
			return nil, errors.New("topic or chain is not provided")
		}

		topic, err := s.GetTopicByChain(input.Chain)
		if err != nil {
			return nil, err
		}
		input.Topic = topic.ID
	}

	messages := []model.RawMessage{}

	scope := s.
		Model(&model.RawMessage{}).
		Where("topic_id = ?", input.Topic).
		Where("id >= ?", input.StartID)

	if processed := input.Processed; processed != nil {
		if *processed {
			scope = scope.Where("processed_at IS NOT NULL")
		} else {
			scope = scope.Where("processed_at IS NULL")
		}
	}

	err := scope.
		Order("id ASC").
		Limit(input.Limit).
		Find(&messages).
		Error

	return messages, checkErr(err)
}

func (s RawMessagesStore) MarkMessageProcessed(message *model.RawMessage) error {
	return s.Exec("UPDATE raw_messages SET processed_at = ? WHERE id = ?", time.Now().UTC(), message.ID).Error
}
