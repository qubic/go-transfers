//go:build !ci
// +build !ci

package client

import (
	"context"
	"flag"
	"github.com/ardanlabs/conf"
	"github.com/gookit/slog"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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
	const envPrefix = "QUBIC_TRANSFERS"
	err := godotenv.Load("../.env.local")
	if err != nil {
		slog.Info("Using no env file")
	}
	var config struct {
		Client struct {
			EventApiUrl string `conf:"required"`
			CoreApiUrl  string `conf:"required"`
		}
	}
	err = conf.Parse(os.Args[1:], envPrefix, &config)
	if err != nil {
		slog.Error("error getting config", "err", err)
		os.Exit(-1)
	}
	eventClient, err = NewIntegrationEventClient(config.Client.EventApiUrl, config.Client.CoreApiUrl)
	if err != nil {
		slog.Error("error creating event client", "err", err)
		os.Exit(-1)
	}
}
