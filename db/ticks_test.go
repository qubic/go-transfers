package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// tick
func TestPgRepository_GetOrCreateTick_GivenNewTick_ThenCreate(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(context.Background(), 42)
	assert.Nil(t, err)
	assert.Greater(t, tickId, 0)

	// clean up
	deleteTick(tickId, t)
}

func TestPgRepository_GetOrCreateTick_GivenTick_ThenGet(t *testing.T) {
	tickId, err := repository.insertTick(context.Background(), 42)
	assert.Nil(t, err)
	assert.Greater(t, tickId, 0)

	reloaded, err := repository.GetOrCreateTick(context.Background(), 42)
	assert.Nil(t, err)
	assert.Equal(t, tickId, reloaded)

	// clean up
	deleteTick(tickId, t)
}
