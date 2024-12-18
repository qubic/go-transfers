//go:build !ci
// +build !ci

package sync

import (
	"context"
	"flag"
	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go-transfers/client"
	"go-transfers/config"
	"go-transfers/db"
	"os"
	"testing"
	"time"
)

var (
	eventClient       EventClient
	eventService      *EventService
	postgresContainer *postgres.PostgresContainer
	repository        *db.PgRepository
)

func TestEventService_GetEventRange(t *testing.T) {

	tick, err := eventService.ProcessTickEvents(context.Background(), 17603769, 17603770)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 17603769, tick)

}

// test setup

func TestMain(m *testing.M) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	setup()
	// Parse args and run
	flag.Parse()
	exitCode := m.Run()
	// Exit
	tearDown()
	os.Exit(exitCode)
}

func tearDown() {
	repository.Close()
	err := testcontainers.TerminateContainer(postgresContainer)
	if err != nil {
		slog.Error("terminating the postgres container", "error", err)
	}
}

func setup() {
	c, err := config.GetConfig("..")
	if err != nil {
		slog.Error("error getting config")
		os.Exit(-1)
	}
	eventClient, err = client.NewIntegrationEventClient(c.Client.EventApiUrl, c.Client.CoreApiUrl)
	if err != nil {
		slog.Error("error creating event client")
		os.Exit(-1)
	}

	repository = db.NewRepository(setupDatabase(context.Background()))
	eventProcessor := NewEventProcessor(repository)
	eventService, err = NewEventService(eventClient, eventProcessor, repository)
	if err != nil {
		slog.Error("error creating event service")
		os.Exit(-1)
	}
}

func createPgContainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx,
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
		return nil, err
	}
	return pgContainer, nil
}

func setupDatabase(ctx context.Context) *sqlx.DB {
	var err error
	postgresContainer, err = createPgContainer(ctx)
	if err != nil {
		slog.Error("setting up test database", "error", err)
		os.Exit(-1)
	}
	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		slog.Error("getting connection string", "error", err)
		os.Exit(-1)
	}
	slog.Info("DB", "connection-string", connectionString)
	pgDb, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		slog.Error("connecting to database", "error", err)
		os.Exit(-1)
	}
	err = db.Migrate(connectionString)
	if err != nil {
		slog.Error("migrating database", "error", err)
		os.Exit(-1)
	}
	return pgDb
}
