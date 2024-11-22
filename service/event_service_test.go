package service

import (
	"go-transfers/client"
	"testing"
)

const (
	targetUrl string = "0.0.0.0:8003" // TODO mock
)

func Test_GetEventRange(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	eventClient, newClientErr := client.NewEventClient(targetUrl)
	if newClientErr != nil {
		t.Error(newClientErr)
	}

	eventService := NewEventService(eventClient)
	err := eventService.ReadEvents(17313985, 17313986)
	if err != nil {
		t.Error(err)
	}
}
