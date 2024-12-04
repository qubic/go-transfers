package service

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
	repository   Repository
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

	eventClient, err = client.NewIntegrationEventClient(c.Client.EventApiUrl, c.Client.CoreApiUrl)
	if err != nil {
		slog.Error("error creating event client")
		os.Exit(-1)
	}

	repository, err = db.NewRepository(&c.Database)
	if err != nil {
		slog.Error("error creating repository")
		os.Exit(-1)
	}

	eventProcessor := NewEventProcessor(repository)
	eventService = NewEventService(eventClient, eventProcessor)
}

func teardown() {
	repository.Close()
}
