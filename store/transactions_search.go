package store

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/figment-networks/avalanche-indexer/model"
)

var transactionTypes = map[string]bool{}

func init() {
	for _, name := range model.TransactionTypes {
		transactionTypes[name] = true
	}
}

type TxSearchInput struct {
	Chain       string `form:"chain"`
	Type        string `form:"type"`
	Address     string `form:"address"`
	Asset       string `form:"asset"`
	Memo        string `form:"memo"`
	Order       string `form:"order"`
	Limit       int    `form:"limit"`
	Offset      int    `form:"offset"`
	Page        int    `form:"page"`
	StartTime   string `form:"start_time"`
	EndTime     string `form:"end_time"`
	StartHeight int    `form:"start_height"`
	EndHeight   int    `form:"end_height"`
	BlockHash   string `form:"block_hash"`

	startTime *time.Time
	endTime   *time.Time
	types     []string
}

func (input *TxSearchInput) Validate() error {
	if input.Memo != "" && len(input.Memo) < 3 {
		return errors.New("memo field is too short")
	}

	if input.Limit == 0 {
		input.Limit = 25
	}
	if input.Limit > 100 {
		return errors.New("maximum limit value is 100")
	}

	if input.Page > 0 && input.Offset == 0 {
		input.Offset = (input.Page - 1) * input.Limit
	}

	if input.StartTime != "" {
		ts, err := parseTimeFilter(input.StartTime, "bod")
		if err != nil {
			return errors.New("invalid start time")
		}
		input.startTime = ts
	}

	if input.EndTime != "" {
		ts, err := parseTimeFilter(input.EndTime, "eod")
		if err != nil {
			return errors.New("invalid end time")
		}
		if input.startTime != nil && ts.Before(*input.startTime) {
			return errors.New("end time must be greater than start time")
		}
		input.endTime = ts
	}

	if input.Type != "" {
		input.types = strings.Split(input.Type, ",")

		for _, t := range input.types {
			if !transactionTypes[t] {
				return fmt.Errorf("invalid transaction type: %s", t)
			}
		}
	}

	return nil
}

func parseTimeFilter(input string, mode string) (*time.Time, error) {
	if input == "" {
		return nil, nil
	}

	var t time.Time
	var err error

	if reDate.MatchString(input) {
		t, err = time.Parse("2006-01-02", input)
		if err == nil {
			switch mode {
			case "eod":
				t = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
			}
		}
	} else if reUnixTime.MatchString(input) {
		unixTime, err := strconv.Atoi(input)
		if err != nil {
			return nil, err
		}
		t = time.Unix(int64(unixTime), 0)
	} else {
		t, err = time.Parse(time.RFC3339, input)
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

type TxSearchOutput struct {
	Transactions []model.Transaction `json:"data"`
}
