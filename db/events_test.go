package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPgRepository_GetOrCreateEvent_GivenNoEvent_ThenInsert(t *testing.T) {
	tickId, transactionId := setupTransactionTestData(t)

	eventId, err := repository.GetOrCreateEvent(context.Background(), transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	assert.Greater(t, eventId, 0)

	// clean up
	deleteEvent(eventId, t)
	cleanUpTransactionTestData(t, transactionId, tickId)
}

func TestPgRepository_GetOrCreateEvent_GivenEvent_ThenGet(t *testing.T) {
	tickId, transactionId := setupTransactionTestData(t)
	eventId, err := repository.insertEvent(context.Background(), transactionId, 1, 2, "foo")
	assert.Nil(t, err)

	reloaded, err := repository.GetOrCreateEvent(context.Background(), transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	assert.Equal(t, eventId, reloaded)

	// clean up
	deleteEvent(eventId, t)
	cleanUpTransactionTestData(t, transactionId, tickId)
}

// qu transfer event

func TestPgRepository_GetOrCreateQuTransferEvent_GivenNoTransferEvent_ThenCreate(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 0)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)

	transferId, err := repository.GetOrCreateQuTransferEvent(context.Background(), eventId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Greater(t, transferId, 0)

	// clean up
	deleteTransferQuEvent(transferId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetOrCreateQuTransferEvent_GivenTransferEvent_ThenGet(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 0)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	transferId, err := repository.insertQuTransferEvent(context.Background(), eventId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)

	reloaded, err := repository.GetOrCreateQuTransferEvent(context.Background(), eventId, sourceEntityId, destinationEntityId, 123)
	assert.Nil(t, err)
	assert.Equal(t, transferId, reloaded)

	// clean up
	deleteTransferQuEvent(transferId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

// asset change event (ownership or possession change)

func TestPgRepository_GetOrCreateAssetChangeEvent_GivenNone_ThenCreate(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 2)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId(context.Background(), AAA, "QX") // don't clean up
	assert.Nil(t, err)

	assetEventId, err := repository.GetOrCreateAssetChangeEvent(context.Background(), eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Greater(t, assetEventId, 0)

	// clean up
	deleteAssetChangeEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetOrCreateAssetChangeEvent_GivenEntry_ThenGet(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 3)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId(context.Background(), AAA, "QX") // don't clean up
	assert.Nil(t, err)

	assetEventId, err := repository.insertAssetChangeEvent(context.Background(), eventId, assetId, sourceEntityId, destinationEntityId, 123456789)

	reloaded, err := repository.GetOrCreateAssetChangeEvent(context.Background(), eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Equal(t, assetEventId, reloaded)

	// clean up
	deleteAssetChangeEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}
