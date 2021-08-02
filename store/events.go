package store

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/figment-networks/avalanche-indexer/model"
)

type EventsStore struct {
	*gorm.DB
}

func NewEventsStore(db *gorm.DB) EventsStore {
	return EventsStore{
		DB: db,
	}
}

// FindByID returns an event by ID
func (s EventsStore) FindByID(id string) (*model.Event, error) {
	result := &model.Event{}
	err := s.First(result, "id = ?", id).Error

	return result, checkErr(err)
}

// Create creates a new event record or ignores if it already exists
func (s EventsStore) Create(event *model.Event) error {
	if event.ID == "" {
		if err := event.AssignID(); err != nil {
			return err
		}
	}

	return s.
		Model(event).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(event).
		Error
}

// Search returns event records matching the search input
func (s EventsStore) Search(input *EventSearchInput) ([]model.Event, error) {
	result := []model.Event{}
	scope := s.Model(&model.Event{})

	if input.Chain != "" {
		scope = scope.Where("chain = ?", input.Chain)
	}

	if input.Scope != "" {
		scope = scope.Where("scope = ?", input.Scope)
	}

	if input.Type != "" {
		types := strings.Split(input.Type, ",")
		scope = scope.Where("type IN (?)", types)
	}

	if input.ItemID != "" && input.ItemType != "" {
		scope = scope.Where("item_id = ? AND item_type = ?", input.ItemID, input.ItemType)
	}

	if input.startTime != nil {
		scope = scope.Where("timestamp >= ?", *input.startTime)
	}

	if input.endTime != nil {
		scope = scope.Where("timestamp <= ?", *input.endTime)
	}

	if input.StartHeight > 0 {
		scope = scope.Where("block_height >= ?", input.StartHeight)
	}

	if input.EndHeight > 0 {
		scope = scope.Where("block_height <= ?", input.EndHeight)
	}

	err := scope.
		Order("timestamp DESC").
		Limit(input.Limit).
		Offset(input.Offset).
		Find(&result).
		Error

	return result, err
}
