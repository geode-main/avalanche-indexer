package store

import (
	"errors"
	"time"
)

type EventSearchInput struct {
	Type        string `form:"type"`
	Scope       string `form:"scope"`
	Chain       string `form:"chain"`
	ItemID      string `form:"item_id"`
	ItemType    string `form:"item_type"`
	StartTime   string `form:"start_time"`
	EndTime     string `form:"end_time"`
	StartHeight int    `form:"start_height"`
	EndHeight   int    `form:"end_height"`
	Limit       int    `form:"limit"`
	Offset      int    `form:"offset"`
	Page        int    `form:"page"`

	startTime *time.Time
	endTime   *time.Time
}

func (input *EventSearchInput) Validate() error {
	if input.ItemID != "" && input.ItemType == "" {
		return errors.New("item_type parameter is required")
	}

	if input.ItemType != "" && input.ItemID == "" {
		return errors.New("item_id parameter is required")
	}

	if input.StartTime != "" {
		ts, err := parseTimeFilter(input.StartTime)
		if err != nil {
			return errors.New("invalid start time")
		}
		input.startTime = ts
	}

	if input.EndTime != "" {
		ts, err := parseTimeFilter(input.EndTime)
		if err != nil {
			return errors.New("invalid end time")
		}
		input.endTime = ts
	}

	if input.Limit < 0 {
		return errors.New("invalid limit value")
	}
	if input.Limit == 0 {
		input.Limit = 100
	}
	if input.Limit > 1000 {
		return errors.New("limit param max value is 1000")
	}

	if input.Offset < 0 {
		return errors.New("invalid offset value")
	}

	if input.Page < 0 {
		return errors.New("invalid page value")
	}
	if input.Page > 0 {
		input.Offset = input.Limit * input.Page
	}

	return nil
}
