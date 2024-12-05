package db

import (
	"database/sql"
	"flag"
	"github.com/stretchr/testify/assert"
	"go-transfers/config"
	"log/slog"
	"os"
	"testing"
)

var (
	repository *PgRepository
)

// entity

func TestPgRepository_InsertEntity(t *testing.T) {
	entityId, err := repository.insertEntity("INSERTED")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)

	// clean up
	deleteEntity(entityId, t)
}

func TestPgRepository_GetEntityId_ThenReturnId(t *testing.T) {
	entityId, err := repository.getEntityId("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)
}

func TestPgRepository_GetEntityId_GivenUnknown_ThenErrNoRows(t *testing.T) {
	_, err := repository.getEntityId("UNKNOWN-IDENTITY")
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestPgRepository_GetOrCreateEntity_GivenNoneThenCreate(t *testing.T) {
	entityId, err := repository.GetOrCreateEntity("TEST-IDENTITY")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)

	// clean up
	deleteEntity(entityId, t)
}

func TestPgRepository_GetOrCreateEntity_GivenEntity_ThenGet(t *testing.T) {
	entityId, err := repository.insertEntity("MANUALLY-INSERTED")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)

	result, err := repository.GetOrCreateEntity("MANUALLY-INSERTED")
	assert.Nil(t, err)
	assert.Equal(t, entityId, result) // same entity found

	// clean up
	deleteEntity(entityId, t)
}

// asset

func TestPgRepository_insertAsset(t *testing.T) {

	entityId, err := repository.insertEntity("TEST-ISSUER")
	assert.Nil(t, err)

	assetId, err := repository.insertAsset(entityId, "TEST-ASSET")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
	deleteEntity(entityId, t)

}

func TestPgRepository_getAssetId_GivenUnknown_ThenErrNoRows(t *testing.T) {
	_, err := repository.getAssetId("FOO", "QX")
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestPgRepository_GetOrCreateAsset_GivenAsset_ThenGet(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB", "QX")
	assert.Nil(t, err)
	assert.Equal(t, 1, assetId) // in seed data
}

func TestPgRepository_GetOrCreateAsset_GivenNoEntity_ThenCreateEntityAndAsset(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset("FOO", "BAR")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
	id, err := repository.getEntityId("FOO")
	assert.Nil(t, err)
	deleteEntity(id, t)
}

func TestPgRepository_GetOrCreateAsset_GivenNoAsset_ThenCreate(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB", "UNKNOWN")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
}

// tick
func TestPgRepository_GetOrCreateTick_GivenNewTick_ThenCreate(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)
	assert.Greater(t, tickId, 0)

	// clean up
	deleteTick(tickId, t)
}

func TestPgRepository_GetOrCreateTick_GivenTick_ThenGet(t *testing.T) {
	tickId, err := repository.insertTick(42)
	assert.Nil(t, err)
	assert.Greater(t, tickId, 0)

	reloaded, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)
	assert.Equal(t, tickId, reloaded)

	// clean up
	deleteTick(tickId, t)
}

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

// event

func TestPgRepository_GetOrCreateEvent_GivenNoEvent_ThenInsert(t *testing.T) {
	tickId, transactionId := setupTransactionTestData(t)

	eventId, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	assert.Greater(t, eventId, 0)

	// clean up
	deleteEvent(eventId, t)
	cleanUpTransactionTestData(t, transactionId, tickId)
}

func TestPgRepository_GetOrCreateEvent_GivenEvent_ThenGet(t *testing.T) {
	tickId, transactionId := setupTransactionTestData(t)
	eventId, err := repository.insertEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)

	reloaded, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	assert.Equal(t, eventId, reloaded)

	// clean up
	deleteEvent(eventId, t)
	cleanUpTransactionTestData(t, transactionId, tickId)
}

// qu transfer event

func TestPgRepository_GetOrCreateQuTransferEvent_GivenNoTransferEvent_ThenCreate(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)

	transferId, err := repository.GetOrCreateQuTransferEvent(eventId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Greater(t, transferId, 0)

	// clean up
	deleteTransferQuEvent(transferId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetOrCreateQuTransferEvent_GivenTransferEvent_ThenGet(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	transferId, err := repository.insertQuTransferEvent(eventId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)

	reloaded, err := repository.GetOrCreateQuTransferEvent(eventId, sourceEntityId, destinationEntityId, 123)
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
	tickId, transactionId, eventId := setupEventTestData(t)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB", "QX") // don't clean up
	assert.Nil(t, err)

	assetEventId, err := repository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Greater(t, assetEventId, 0)

	// clean up
	deleteAssetChangeEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func TestPgRepository_GetOrCreateAssetChangeEvent_GivenEntry_ThenGet(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t)
	sourceEntityId, destinationEntityId := setupSourceAndDestinationEntity(t)
	assetId, err := repository.getAssetId("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB", "QX") // don't clean up
	assert.Nil(t, err)

	assetEventId, err := repository.insertAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId, 123456789)

	reloaded, err := repository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Equal(t, assetEventId, reloaded)

	// clean up
	deleteAssetChangeEvent(assetEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

// asset issuance event

func TestPgRepository_GetOrCreateAssetIssuanceEvent_GivenNone_ThenCreate(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t)
	issuerId, err := repository.GetOrCreateEntity("TEST_ISSUER_ID")
	assert.Nil(t, err)
	assetId, err := repository.GetOrCreateAsset("TEST_ISSUER_ID", "A-NAME")
	assert.Nil(t, err)

	issuanceEventId, err := repository.GetOrCreateAssetIssuanceEvent(eventId, assetId, 1234567890, []byte{0, 0, 0, 0, 0, 0, 0}, 0)
	assert.Nil(t, err)
	assert.Greater(t, issuanceEventId, 0)

	deleteAssetIssuanceEvent(issuanceEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteAsset(assetId, t)
	deleteEntity(issuerId, t)
}

func TestPgRepository_GetOrCreateAssetIssuanceEvent_GivenEvent_ThenGet(t *testing.T) {
	tickId, transactionId, eventId := setupEventTestData(t)
	issuerId, err := repository.GetOrCreateEntity("TEST_ISSUER_ID")
	assert.Nil(t, err)
	assetId, err := repository.GetOrCreateAsset("TEST_ISSUER_ID", "A-NAME")
	assert.Nil(t, err)
	issuanceEventId, err := repository.insertAssetIssuanceEvent(eventId, assetId, 1234567890, []byte{0, 0, 1, 0, 0, 0, 0}, 2)
	assert.Nil(t, err)

	reloaded, err := repository.GetOrCreateAssetIssuanceEvent(eventId, assetId, 0, nil, 0)
	assert.Nil(t, err)
	assert.Equal(t, issuanceEventId, reloaded)

	deleteAssetIssuanceEvent(issuanceEventId, t)
	cleanupEventTestData(t, transactionId, tickId, eventId)
	deleteAsset(assetId, t)
	deleteEntity(issuerId, t)
}

func TestPgRepository_GetNumericValue(t *testing.T) {
	value, err := repository.GetNumericValue("tick")
	assert.Nil(t, err)
	assert.True(t, value >= 0)
}

func TestPgRepository_UpdatedNumericValue(t *testing.T) {
	original, err := repository.GetNumericValue("tick")
	assert.Nil(t, err)

	err = repository.UpdateNumericValue("tick", 42)
	assert.Nil(t, err)
	updated, err := repository.GetNumericValue("tick")
	assert.Nil(t, err)
	assert.Equal(t, 42, updated)

	_ = repository.UpdateNumericValue("tick", original) // clean up
}

// test data set ups and clean ups

func setupTransactionTestData(t *testing.T) (int, int) {
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)
	transactionId, err := repository.GetOrCreateTransaction("test-hash", tickId)
	assert.Nil(t, err)
	return tickId, transactionId
}

func cleanUpTransactionTestData(t *testing.T, transactionId int, tickId int) {
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}

func setupEventTestData(t *testing.T) (int, int, int) {
	tickId, transactionId := setupTransactionTestData(t)
	eventId, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	return tickId, transactionId, eventId
}

func cleanupEventTestData(t *testing.T, transactionId int, tickId int, eventId int) {
	deleteEvent(eventId, t)
	cleanUpTransactionTestData(t, transactionId, tickId)
}

func setupSourceAndDestinationEntity(t *testing.T) (int, int) {
	entityId1, err := repository.GetOrCreateEntity("FIRST_ENTITY_ID")
	assert.Nil(t, err)
	entityId2, err := repository.GetOrCreateEntity("SECOND_ENTITY_ID")
	assert.Nil(t, err)
	return entityId1, entityId2
}

func deleteEntity(id int, t *testing.T) {
	count, err := repository.delete(`delete from entities where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteAsset(id int, t *testing.T) {
	count, err := repository.delete(`delete from assets where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteTransaction(id int, t *testing.T) {
	count, err := repository.delete(`delete from transactions where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteTick(id int, t *testing.T) {
	count, err := repository.delete(`delete from ticks where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteEvent(id int, t *testing.T) {
	count, err := repository.delete(`delete from events where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteTransferQuEvent(id int, t *testing.T) {
	count, err := repository.delete(`delete from qu_transfer_events where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteAssetChangeEvent(id int, t *testing.T) {
	count, err := repository.delete(`delete from asset_change_events where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func deleteAssetIssuanceEvent(id int, t *testing.T) {
	count, err := repository.delete(`delete from asset_issuance_events where id = $1;`, id)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func (r *PgRepository) delete(statement string, args ...interface{}) (int64, error) {
	res := r.db.MustExec(statement, args...)
	return res.RowsAffected()
}

// test case infrastructure

func TestMain(m *testing.M) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	setup()
	// Parse args and run
	flag.Parse()
	exitCode := m.Run()
	teardown()
	// Exit
	os.Exit(exitCode)
}

func setup() {
	c, err := config.GetConfig("..")
	if err != nil {
		slog.Error("error getting config")
		os.Exit(-1)
	}

	db, err := CreateDatabaseWithConfig(&c.Database)
	if err != nil {
		slog.Error("error creating repository")
		os.Exit(-1)
	}
	repository = NewRepository(db)
}

func teardown() {
	repository.Close()
}
