package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountPercentOf(t *testing.T) {
	assert.Equal(t, 100.0, NewInt64Amount(1000).PercentOf(NewInt64Amount(1000)))
	assert.Equal(t, 25.0, NewInt64Amount(250).PercentOf(NewInt64Amount(1000)))
	assert.Equal(t, 0.0, NewInt64Amount(0).PercentOf(NewInt64Amount(0)))
}
