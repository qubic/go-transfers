package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// transaction
func TestPgRepository_GetOrCreateTransaction_GivenNoTransaction_ThenInsert(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(context.Background(), 42)
	assert.Nil(t, err)

	transactionId, err := repository.GetOrCreateTransaction(context.Background(), "test-hash", tickId)
	assert.Nil(t, err)
	assert.Greater(t, transactionId, 0)

	// clean up
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}

func TestPgRepository_GetOrCreateTransaction_GivenTransaction_ThenGet(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(context.Background(), 42)
	assert.Nil(t, err)

	transactionId, err := repository.insertTransaction(context.Background(), "test-hash", tickId)
	assert.Nil(t, err)
	assert.Greater(t, transactionId, 0)

	reloaded, err := repository.GetOrCreateTransaction(context.Background(), "test-hash", tickId)
	assert.Nil(t, err)
	assert.Equal(t, transactionId, reloaded)

	// clean up
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}
