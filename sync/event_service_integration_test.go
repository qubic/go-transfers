package sync

import (
	"flag"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
	"go-transfers/client"
	"go-transfers/config"
	"go-transfers/db"
	"os"
	"testing"
)

var (
	eventClient  EventClient
	eventService *EventService
	repository   *db.PgRepository
)

func TestEventService_GetEventRange(t *testing.T) {

	tick, err := eventService.ProcessTickEvents(17603769, 17603770)
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

	dbc := c.Database
	pgDb, err := db.CreateDatabase(dbc.User, dbc.Pass, dbc.Name, dbc.Host, dbc.Port, dbc.MaxOpen, dbc.MaxIdle)
	if err != nil {
		slog.Error("error creating database")
		os.Exit(-1)
	}

	repository = db.NewRepository(pgDb)
	eventProcessor := NewEventProcessor(repository)
	eventService, err = NewEventService(eventClient, eventProcessor, repository)
	if err != nil {
		slog.Error("error creating event service")
		os.Exit(-1)
	}
}
