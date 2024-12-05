package sync

import (
	"flag"
	"go-transfers/client"
	"go-transfers/config"
	"go-transfers/db"
	"log/slog"
	"os"
	"testing"
)

var (
	eventClient  EventClient
	eventService *EventService
)

func TestEventService_GetEventRange(t *testing.T) {

	err := eventService.ProcessTickEvents(17396000, 17396010)
	if err != nil {
		t.Error(err)
	}
	err = eventService.ProcessTickEvents(16660843, 16660845)
	if err != nil {
		t.Error(err)
	}
}

//func TestEventService_Loop(t *testing.T) {
//	go eventService.SyncInLoop()
//	time.Sleep(time.Second * 30)
//}

// test setup

func TestMain(m *testing.M) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	setup()
	// Parse args and run
	flag.Parse()
	exitCode := m.Run()
	// Exit
	os.Exit(exitCode)
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

	pgDb, err := db.CreateDatabaseWithConfig(&c.Database)
	// defer pgDb.Close()
	if err != nil {
		slog.Error("error creating database")
		os.Exit(-1)
	}

	repository := db.NewRepository(pgDb)
	eventProcessor := NewEventProcessor(repository)
	eventService, err = NewEventService(eventClient, eventProcessor, repository)
	if err != nil {
		slog.Error("error creating event service")
		os.Exit(-1)
	}
}
