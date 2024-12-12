//go:build !ci
// +build !ci

package client

import (
	"context"
	"flag"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
	"go-transfers/config"
	"os"
	"testing"
)

var (
	eventClient *IntegrationEventClient
)

func TestEventClient_GetEvents(t *testing.T) {
	const tickNumber uint32 = 17302596
	tickEvents, err := eventClient.GetEvents(context.Background(), tickNumber)
	assert.Nil(t, err)
	slog.Info("Received tick events.", "tick", tickNumber, "events", tickEvents)
}

func TestEventClient_GetStatus(t *testing.T) {
	status, err := eventClient.GetStatus(context.Background())
	assert.Nil(t, err)
	assert.NotNil(t, status.AvailableTick, "last processed tick is nil")
	slog.Info("Received event status", "event status", status)
}

func TestEventClient_GetTickInfo(t *testing.T) {
	info, err := eventClient.GetTickInfo(context.Background())
	assert.Nil(t, err)
	slog.Info("Received tick info", "tick info", info)
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

	eventClient, err = NewIntegrationEventClient(c.Client.EventApiUrl, c.Client.CoreApiUrl)
	if err != nil {
		slog.Error("error creating event client")
		os.Exit(-1)
	}
}
