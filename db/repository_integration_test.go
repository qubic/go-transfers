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
	setup()
	// Parse args and run
	flag.Parse()
	exitCode := m.Run()
	teardown()
	// Exit
	os.Exit(exitCode)
}

func Test_GetOrCreateEntity(t *testing.T) {
	entityId, err := repository.GetOrCreateEntity("TEST-IDENTITY")
	assert.Nil(t, err)
	assert.Greater(t, entityId, int64(0))

	// clean up
	count, err := repository.DeleteEntity(entityId)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func Test_GetEntityId_ThenReturnId(t *testing.T) {
	entityId, err := repository.GetEntityId("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB")
	assert.Nil(t, err)
	assert.Greater(t, entityId, int64(0))
}

func Test_GetEntityId_GivenUnknown_ThenReturnMinusOne(t *testing.T) {
	entityId, err := repository.GetEntityId("TEST-IDENTITY")
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Equal(t, int64(0), entityId)
}

func (r *PgRepository) DeleteEntity(id int64) (int64, error) {
	deleteSql := `delete from entities where id = $1;`
	res, err := r.db.Exec(deleteSql, id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
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

func teardown() {
	repository.Close()
}
