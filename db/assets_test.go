package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

// asset

func TestPgRepository_insertAsset(t *testing.T) {

	entityId, err := repository.insertEntity(context.Background(), "TEST-ISSUER")
	assert.Nil(t, err)

	assetId, err := repository.insertAsset(context.Background(), entityId, "TEST-ASSET")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
	deleteEntity(entityId, t)

}

func TestPgRepository_getAssetId_GivenUnknown_ThenErrNoRows(t *testing.T) {
	_, err := repository.getAssetId(context.Background(), "FOO", "QX")
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestPgRepository_GetOrCreateAsset_GivenAsset_ThenGet(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset(context.Background(), AAA, "QX")
	assert.Nil(t, err)
	assert.Equal(t, 1, assetId) // in seed data
}

func TestPgRepository_GetOrCreateAsset_GivenNoEntity_ThenCreateEntityAndAsset(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset(context.Background(), "FOO", "BAR")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
	id, err := repository.getEntityId(context.Background(), "FOO")
	assert.Nil(t, err)
	deleteEntity(id, t)
}

func TestPgRepository_GetOrCreateAsset_GivenNoAsset_ThenCreate(t *testing.T) {
	assetId, err := repository.GetOrCreateAsset(context.Background(), AAA, "UNKNOWN")
	assert.Nil(t, err)
	assert.Greater(t, assetId, 0)

	// clean up
	deleteAsset(assetId, t)
}
