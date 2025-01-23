package db

import (
	"context"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"go-transfers/proto"
	"testing"
)

func TestPgRepository_GetAssetChangeEventsForTick(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 2)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId(context.Background(), AAA, "QX") // don't clean up
	assert.Nil(t, err)
	assetEventId, err := repository.insertAssetChangeEvent(context.Background(), eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
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

func TestPgRepository_GetAssetIssuanceEventsForTick(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 1)
	assetId, err := repository.getAssetId(context.Background(), AAA, "QX") // don't clean up
	assert.Nil(t, err)
	assetEventId, err := repository.insertAssetIssuanceEvent(context.Background(), eventId, assetId, 123456789, base64.StdEncoding.EncodeToString([]byte{0, 0, 0, 0, 0, 0, 0}), 2)
	assert.Nil(t, err)

	events, err := repository.GetAssetIssuanceEventsForTick(context.Background(), testTickNumber)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(events))
	assert.Equal(t, &proto.AssetIssuanceEvent{
		IssuerId:              AAA,
		Name:                  "QX",
		NumberOfShares:        123456789,
		UnitOfMeasurement:     "AAAAAAAAAA==",
		NumberOfDecimalPlaces: 2,
		TransactionHash:       testTransactionHash,
		Tick:                  testTickNumber,
		EventType:             1,
	}, events[0])

	deleteAssetIssuanceEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
}

func TestPgRepository_GetQuTransferEventsForTick(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t, 0)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)

	transferId, err := repository.GetOrCreateQuTransferEvent(context.Background(), eventId, sourceEntityId, destinationEntityId, 123_456_789_012_345)
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

	transferId, err := repository.GetOrCreateQuTransferEvent(context.Background(), eventId, sourceEntityId, destinationEntityId, 123_456_789_012_345)
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
	assetId, err := repository.getAssetId(context.Background(), AAA, "QX") // don't clean up
	assert.Nil(t, err)
	assetEventId, err := repository.insertAssetChangeEvent(context.Background(), eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
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
