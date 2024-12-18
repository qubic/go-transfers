package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

// entity

func TestPgRepository_InsertEntity(t *testing.T) {
	entityId, err := repository.insertEntity(context.Background(), "INSERTED")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)

	// clean up
	deleteEntity(entityId, t)
}

func TestPgRepository_GetEntityId_ThenReturnId(t *testing.T) {
	entityId, err := repository.getEntityId(context.Background(), AAA)
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)
}

func TestPgRepository_GetEntityId_GivenUnknown_ThenErrNoRows(t *testing.T) {
	_, err := repository.getEntityId(context.Background(), "UNKNOWN-IDENTITY")
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestPgRepository_GetOrCreateEntity_GivenNoneThenCreate(t *testing.T) {
	entityId, err := repository.GetOrCreateEntity(context.Background(), "TEST-IDENTITY")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)

	// clean up
	deleteEntity(entityId, t)
}

func TestPgRepository_GetOrCreateEntity_GivenEntity_ThenGet(t *testing.T) {
	entityId, err := repository.insertEntity(context.Background(), "MANUALLY-INSERTED")
	assert.Nil(t, err)
	assert.Greater(t, entityId, 0)

	result, err := repository.GetOrCreateEntity(context.Background(), "MANUALLY-INSERTED")
	assert.Nil(t, err)
	assert.Equal(t, entityId, result) // same entity found

	// clean up
	deleteEntity(entityId, t)
}
