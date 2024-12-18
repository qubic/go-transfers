package db

import (
	"context"
	"flag"
	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

var (
	repository        *PgRepository
	postgresContainer testcontainers.Container
)

const (
	testTickNumber        = 42
	AAA                   = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB"
	testSourceIdentity    = "SOURCE_IDENTITY"
	testDestinationEntity = "TARGET_IDENTITY"
	testTransactionHash   = "test-hash"
)

// test data set-ups and clean-ups

func setupTransactionTestData(t *testing.T) (int, int) {
	tickId, err := repository.GetOrCreateTick(context.Background(), testTickNumber)
	assert.Nil(t, err)
	transactionId, err := repository.GetOrCreateTransaction(context.Background(), testTransactionHash, tickId)
	assert.Nil(t, err)
	return tickId, transactionId
}

func cleanUpTransactionTestData(t *testing.T, transactionId int, tickId int) {
	deleteTransaction(transactionId, t)
	deleteTick(tickId, t)
}

func setupEventTestData(t *testing.T, eventType uint32) (int, int, int) {
	tickId, transactionId := setupTransactionTestData(t)
	eventId, err := repository.GetOrCreateEvent(context.Background(), transactionId, 1, eventType, "foo")
	assert.Nil(t, err)
	return tickId, transactionId, eventId
}

func cleanupEventTestData(t *testing.T, transactionId int, tickId int, eventId int) {
	deleteEvent(eventId, t)
	cleanUpTransactionTestData(t, transactionId, tickId)
}

func setupSourceAndDestinationEntity(t *testing.T) (int, int) {
	entityId1, err := repository.GetOrCreateEntity(context.Background(), testSourceIdentity)
	assert.Nil(t, err)
	entityId2, err := repository.GetOrCreateEntity(context.Background(), testDestinationEntity)
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
	slog.SetLogLevel(slog.DebugLevel)
	setup()
	// Parse args and run
	flag.Parse()
	exitCode := m.Run()
	teardown()
	// Exit
	os.Exit(exitCode)
}

func setup() {
	connectionString, err := setupDatabase()
	if err != nil {
		slog.Error("setting up test database", "error", err)
		os.Exit(-1)
	}
	slog.Info("DB", "connection-string", connectionString)
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		slog.Error("connecting to database", "error", err)
		os.Exit(-1)
	}
	err = Migrate(connectionString)

	if err != nil {
		slog.Error("migrating database", "error", err)
		os.Exit(-1)
	}
	repository = NewRepository(db)
}

func teardown() {
	repository.Close()
	err := testcontainers.TerminateContainer(postgresContainer)
	if err != nil {
		slog.Error("terminating the postgres container", "error", err)
	}
}

func setupDatabase() (string, error) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		testcontainers.WithLogger(slog.NewStdLogger()),
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForExposedPort(),
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		slog.Error("starting postgres container", "error", err)
		return "", err
	}
	return postgresContainer.ConnectionString(ctx, "sslmode=disable")
}
