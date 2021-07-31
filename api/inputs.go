package api

import (
	"github.com/figment-networks/avalanche-indexer/store"
	"github.com/gin-gonic/gin"
)

func eventsSearchInput(c *gin.Context) *store.EventSearchInput {
	input := &store.EventSearchInput{}

	if err := c.Bind(input); err != nil {
		badRequest(c, err)
		return nil
	}

	if err := input.Validate(); err != nil {
		badRequest(c, err)
		return nil
	}

	return input
}
