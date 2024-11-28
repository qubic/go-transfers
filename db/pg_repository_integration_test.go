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

func TestPgRepository_GetOrCreateEntity_GivenEntity_ThenDoNotInsert(t *testing.T) {
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

func TestPgRepository_getAssetId_ThenReturnId(t *testing.T) {
	assetId, err := repository.getAssetId("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB", "QX")
	assert.Nil(t, err)
	assert.Equal(t, 1, assetId)
}

func TestPgRepository_getAssetId_GivenUnknown_ThenErrNoRows(t *testing.T) {
	_, err := repository.getAssetId("FOO", "QX")
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestPgRepository_GetOrCreateAsset_WhenNewEntity_ThenCreateEntityAndAsset(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset("FOO", "BAR")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
	id, err := repository.getEntityId("FOO")
	assert.Nil(t, err)
	deleteEntity(id, t)
}

func TestPgRepository_GetOrCreateAsset_WhenNewAsset_ThenCreateAsset(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB", "BAR")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
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

func (r *PgRepository) delete(statement string, args ...interface{}) (int64, error) {
	res := r.db.MustExec(statement, args...)
	return res.RowsAffected()
}

func teardown() {
	repository.Close()
}
