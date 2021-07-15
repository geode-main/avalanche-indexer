package store

import "errors"

type BlocksSearch struct {
	Chain       string `form:"chain"`
	Type        string `form:"type"`
	StartHeight int    `form:"start_height"`
	EndHeight   int    `form:"end_height"`
	Order       string `form:"order"`
	Limit       int    `form:"limit"`
	Offset      int    `form:"offset"`
	Page        int    `form:"page"`
}

func (s *BlocksSearch) Validate() error {
	if s.Chain == "" {
		return errors.New("chain ID is required")
	}

	if s.Limit < 0 {
		return errors.New("invalid limit")
	}
	if s.Limit == 0 {
		s.Limit = 100
	}
	if s.Limit > 100 {
		return errors.New("max limit is 100")
	}

	if s.StartHeight < 0 {
		return errors.New("invalid start height")
	}
	if s.EndHeight < 0 {
		return errors.New("invalid end height")
	}

	if s.Order == "" {
		if s.StartHeight > 0 {
			s.Order = "height_asc"
		} else {
			s.Order = "height_desc"
		}
	}

	switch s.Order {
	case "height_asc":
	case "height_desc":
	default:
		return errors.New("invalid order")
	}

	return nil
}
