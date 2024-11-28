package service

import (
	eventspb "github.com/qubic/go-events/proto"
	"testing"
)

type FakeEventClient struct {
	events map[uint32]*eventspb.TickEvents
}

func NewFakeEventClient(tickEvents map[uint32]*eventspb.TickEvents) (*FakeEventClient, error) {
	return &FakeEventClient{events: tickEvents}, nil
}

func (eventClient *FakeEventClient) GetEvents(tickNumber uint32) (*eventspb.TickEvents, error) {
	return eventClient.events[tickNumber], nil
}

func TestEventService_ProcessTickEvents(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	event := event(0, "sMmo18V9WMO9LstUtxvWC2ZfJc2/FZWKEUdAKOqNKDIBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBCDwAAAAAA")

	tx1Events := transactionEvents("tx-id-1", &event)
	tx2Events := transactionEvents("tx-id-2", &event)
	tx3Events := transactionEvents("tx-id-3", &event, &event)

	tickEvents1 := tickEvents(123, &tx1Events, &tx2Events)
	tickEvents2 := tickEvents(124)
	tickEvents3 := tickEvents(125, &tx3Events)

	eventMap := map[uint32]*eventspb.TickEvents{
		123: &tickEvents1,
		124: &tickEvents2,
		125: &tickEvents3,
	}

	fakeEventClient, err := NewFakeEventClient(eventMap)
	if err != nil {
		t.Error(err)
	}

	eventService := NewEventService(fakeEventClient)
	err = eventService.ProcessTickEvents(123, 126)
	if err != nil {
		t.Error(err)
	}

	// TODO verify processing

}

func event(eventType uint32, eventData string) eventspb.Event {
	return eventspb.Event{
		EventType: eventType,
		EventData: eventData,
	}
}

func transactionEvents(txId string, events ...*eventspb.Event) eventspb.TransactionEvents {
	return eventspb.TransactionEvents{
		TxId:   txId,
		Events: events,
	}
}

func tickEvents(tick uint32, transactionEvents ...*eventspb.TransactionEvents) eventspb.TickEvents {
	return eventspb.TickEvents{
		Tick:     tick,
		TxEvents: transactionEvents,
	}
}
