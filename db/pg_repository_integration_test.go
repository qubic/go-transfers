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
	tickId, err := repository.GetOrCreateTick(42)
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

	transactionId, err := repository.GetOrCreateTransaction("test-hash", tickId)
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
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)
	transactionId, err := repository.GetOrCreateTransaction("test-hash", tickId)
	assert.Nil(t, err)

	eventId, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	assert.Greater(t, eventId, 0)

	// clean up
	deleteEvent(eventId, t)
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}

func TestPgRepository_GetOrCreateEvent_GivenEvent_ThenGet(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)
	transactionId, err := repository.GetOrCreateTransaction("test-hash", tickId)
	assert.Nil(t, err)
	eventId, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)

	reloaded, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	assert.Equal(t, eventId, reloaded)

	// clean up
	deleteEvent(eventId, t)
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}

// qu transfer event

func TestPgRepository_GetOrCreateQuTransferEvent(t *testing.T) {
	tickId, err := repository.GetOrCreateTick(42)
	assert.Nil(t, err)
	transactionId, err := repository.GetOrCreateTransaction("test-hash", tickId)
	assert.Nil(t, err)
	eventId, err := repository.GetOrCreateEvent(transactionId, 1, 2, "foo")
	assert.Nil(t, err)
	sourceEntityId, err := repository.GetOrCreateEntity("SOURCE_ID")
	assert.Nil(t, err)
	destinationEntityId, err := repository.GetOrCreateEntity("DESTINATION_ID")
	assert.Nil(t, err)

	transferId, err := repository.GetOrCreateQuTransferEvent(eventId, sourceEntityId, destinationEntityId, 123456789)
	assert.Nil(t, err)
	assert.Greater(t, transferId, 0)

	// clean up
	deleteTransferQuEvent(transferId, t)
	deleteEvent(eventId, t)
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
	deleteEntity(sourceEntityId, t)
	deleteEntity(destinationEntityId, t)
}

func setup() {
	c, err := config.GetConfig("..")
	if err != nil {
		slog.Error("error getting config")
		os.Exit(-1)
	}

	repository, err = NewRepository(&c.Database)
	if err != nil {
		slog.Error("error creating repository")
		os.Exit(-1)
	}
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

func (r *PgRepository) delete(statement string, args ...interface{}) (int64, error) {
	res := r.db.MustExec(statement, args...)
	return res.RowsAffected()
}

func teardown() {
	repository.Close()
}
