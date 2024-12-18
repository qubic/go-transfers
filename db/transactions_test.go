package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// transaction
func TestPgRepository_GetOrCreateTransaction_GivenNoTransaction_ThenInsert(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)

	transactionId, err := repository.GetOrCreateTransaction("test-hash", tickId)
	assert.Nil(t, err)
	assert.Greater(t, transactionId, 0)

	// clean up
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}

func TestPgRepository_GetOrCreateTransaction_GivenTransaction_ThenGet(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)

	transactionId, err := repository.insertTransaction("test-hash", tickId)
	assert.Nil(t, err)
	assert.Greater(t, transactionId, 0)

	reloaded, err := repository.GetOrCreateTransaction("test-hash", tickId)
	assert.Nil(t, err)
	assert.Equal(t, transactionId, reloaded)

	// clean up
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}
