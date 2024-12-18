package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPgRepository_GetLatestTick(t *testing.T) {
	value, err := repository.GetLatestTick(context.Background())
	assert.Nil(t, err)
	assert.True(t, value >= 0)
}

func TestPgRepository_UpdatedNumericValue(t *testing.T) {
	original, err := repository.GetLatestTick(context.Background())
	assert.Nil(t, err)

	err = repository.UpdateLatestTick(context.Background(), 42)
	assert.Nil(t, err)
	updated, err := repository.GetLatestTick(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 42, updated)

	_ = repository.UpdateLatestTick(context.Background(), original) // clean up
}
