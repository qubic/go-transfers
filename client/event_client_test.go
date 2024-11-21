package client

import (
	"log"
	"testing"
)

const (
	targetUrl string = "95.216.243.140:8003"
)

func Test_GetEvents(t *testing.T) {

	eventClient, newErr := NewEventClient(targetUrl)
	if newErr != nil {
		t.Error(newErr)
	}
	const tickNumber uint32 = 17302596
	tickEvents, err := eventClient.GetEvents(tickNumber)
	if err != nil {
		t.Error(err)
	}
	log.Printf("Tick events for tick %d: %+v", tickNumber, tickEvents)
}
