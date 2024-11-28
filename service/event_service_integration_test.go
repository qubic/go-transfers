package service

import (
	"go-transfers/client"
	"log/slog"
	"testing"
)

const (
	targetUrl string = "95.216.243.140:8003"
)

func TestEventService_GetEventRange(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	eventClient, err := client.NewIntegrationEventClient(targetUrl)
	if err != nil {
		t.Error(err)
	}

	eventService := NewEventService(eventClient)
	err = eventService.ProcessTickEvents(17313985, 17313986)
	if err != nil {
		t.Error(err)
	}
}
