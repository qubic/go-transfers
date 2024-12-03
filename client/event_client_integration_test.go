package client

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"go-transfers/config"
	"log/slog"
	"os"
	"testing"
)

var (
	eventClient *IntegrationEventClient
)

func TestEventClient_GetEvents(t *testing.T) {
	const tickNumber uint32 = 17302596
	tickEvents, err := eventClient.GetEvents(tickNumber)
	assert.Nil(t, err)
	slog.Info("Received tick events.", "tick", tickNumber, "events", tickEvents)
}

func TestEventClient_GetStatue(t *testing.T) {
	status, err := eventClient.GetStatus()
	assert.Nil(t, err)
	assert.NotNil(t, status.LastProcessedTick, "last processed tick is nil")
	slog.Info("Received event status", "event status", status)
}

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

	eventClient, err = NewIntegrationEventClient(c.EventClient.TargetUrl)
	if err != nil {
		slog.Error("error creating event client")
		os.Exit(-1)
	}
}
