package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-transfers/proto"
	"testing"
)

func TestPgRepository_GetAssetChangeEventsForTick(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 2)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId(AAA, "QX") // don't clean up
	assert.Nil(t, err)
	assetEventId, err := repository.insertAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)

	events, err := repository.GetAssetChangeEventsForTick(context.Background(), testTickNumber)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.AssetChangeEvent{
		SourceId:        testSourceIdentity,
		DestinationId:   testDestinationEntity,
		IssuerId:        AAA,
		Name:            "QX",
		NumberOfShares:  123456789,
		TransactionHash: testTransactionHash,
		Tick:            testTickNumber,
		EventType:       2,
	}, events[0])

	deleteAssetChangeEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetQuTransferEventsForTick(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 0)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)

	transferId, err := repository.GetOrCreateQuTransferEvent(eventId, sourceEntityId, destinationEntityId, 123_456_789_012_345)
	assert.Nil(t, err)

	events, err := repository.GetQuTransferEventsForTick(context.Background(), testTickNumber)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.QuTransferEvent{
		SourceId:        testSourceIdentity,
		DestinationId:   testDestinationEntity,
		Amount:          123_456_789_012_345,
		TransactionHash: testTransactionHash,
		Tick:            testTickNumber,
		EventType:       0,
	}, events[0])

	// clean up
	deleteTransferQuEvent(transferId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetQuTransferEventsForEntity(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 0)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)

	transferId, err := repository.GetOrCreateQuTransferEvent(eventId, sourceEntityId, destinationEntityId, 123_456_789_012_345)
	assert.Nil(t, err)

	events, err := repository.GetQuTransferEventsForEntity(context.Background(), testSourceIdentity)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.QuTransferEvent{
		SourceId:        testSourceIdentity,
		DestinationId:   testDestinationEntity,
		Amount:          123_456_789_012_345,
		TransactionHash: testTransactionHash,
		Tick:            testTickNumber,
		EventType:       0,
	}, events[0])

	events, err = repository.GetQuTransferEventsForEntity(context.Background(), testDestinationEntity)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.QuTransferEvent{
		SourceId:        testSourceIdentity,
		DestinationId:   testDestinationEntity,
		Amount:          123_456_789_012_345,
		TransactionHash: testTransactionHash,
		Tick:            testTickNumber,
		EventType:       0,
	}, events[0])

	// clean up
	deleteTransferQuEvent(transferId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetAssetChangeEventsForEntity(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 2)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId(AAA, "QX") // don't clean up
	assert.Nil(t, err)
	assetEventId, err := repository.insertAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)

	events, err := repository.GetAssetChangeEventsForEntity(context.Background(), testSourceIdentity)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.AssetChangeEvent{
		SourceId:        testSourceIdentity,
		DestinationId:   testDestinationEntity,
		IssuerId:        AAA,
		Name:            "QX",
		NumberOfShares:  123456789,
		TransactionHash: testTransactionHash,
		Tick:            testTickNumber,
		EventType:       2,
	}, events[0])

	events, err = repository.GetAssetChangeEventsForEntity(context.Background(), testDestinationEntity)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.AssetChangeEvent{
		SourceId:        testSourceIdentity,
		DestinationId:   testDestinationEntity,
		IssuerId:        AAA,
		Name:            "QX",
		NumberOfShares:  123456789,
		TransactionHash: testTransactionHash,
		Tick:            testTickNumber,
		EventType:       2,
	}, events[0])

	deleteAssetChangeEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}
