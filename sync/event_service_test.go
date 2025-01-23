package sync

import (
	"context"
	"github.com/gookit/slog"
	eventspb "github.com/qubic/go-events/proto"
	"github.com/stretchr/testify/assert"
	"go-transfers/client"
	"math/rand/v2"
	"testing"
)

var (
	processedTestTick      = 0
	availableTestTick      = 0
	storedQuTransferEvents = 0
)

type FakeEventClient struct {
	events map[uint32]*eventspb.TickEvents
}

func NewFakeEventClient(tickEvents map[uint32]*eventspb.TickEvents) (*FakeEventClient, error) {
	return &FakeEventClient{events: tickEvents}, nil
}

func (eventClient *FakeEventClient) GetStatus(_ context.Context) (*client.EventStatus, error) {
	return &client.EventStatus{AvailableTick: uint32(availableTestTick)}, nil
}

func (eventClient *FakeEventClient) GetEvents(_ context.Context, tickNumber uint32) (*eventspb.TickEvents, error) {
	return eventClient.events[tickNumber], nil
}

func (eventClient *FakeEventClient) GetTickInfo(_ context.Context) (*client.TickInfo, error) {
	return &client.TickInfo{CurrentTick: uint32(availableTestTick)}, nil
}

type FakeRepository struct {
}

func (f FakeRepository) GetLatestTick(_ context.Context) (int, error) {
	return processedTestTick, nil
}

func (f FakeRepository) UpdateLatestTick(_ context.Context, tickNumber int) error {
	processedTestTick = tickNumber
	return nil
}

func (f FakeRepository) GetOrCreateAssetIssuanceEvent(_ context.Context, _ int, _ int, _ int64, _ string, _ uint32) (int, error) {
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateAsset(_ context.Context, _, _ string) (int, error) {
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateAssetChangeEvent(_ context.Context, _, _, _, _ int, _ int64) (int, error) {
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateEntity(_ context.Context, _ string) (int, error) {
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateQuTransferEvent(_ context.Context, _ int, _ int, _ int, _ uint64) (int, error) {
	storedQuTransferEvents++
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateEvent(_ context.Context, _ int, _ uint64, _ uint32, _ string) (int, error) {
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateTransaction(_ context.Context, _ string, _ int) (int, error) {
	return rand.IntN(1000), nil
}

func (f FakeRepository) GetOrCreateTick(_ context.Context, _ uint32) (int, error) {
	return rand.IntN(1000), nil
}

//goland:noinspection SpellCheckingInspection
func TestEventService_ProcessTickEvents(t *testing.T) {
	slog.SetLogLevel(slog.DebugLevel)

	event := event(0, "sMmo18V9WMO9LstUtxvWC2ZfJc2/FZWKEUdAKOqNKDIBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBCDwAAAAAA", &eventspb.Event_Header{EventId: rand.Uint64N(1000000)})

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

	fakeRepo := &FakeRepository{}

	eventProcessor := EventProcessor{
		repository: fakeRepo,
	}

	processedTestTick = 122
	availableTestTick = 125
	eventService, err := NewEventService(fakeEventClient, &eventProcessor, fakeRepo)
	if err != nil {
		t.Error(err)
	}

	err = eventService.sync(42)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 4, storedQuTransferEvents)
	assert.Equal(t, 125, processedTestTick)

}

func event(eventType uint32, eventData string, header *eventspb.Event_Header) eventspb.Event {
	return eventspb.Event{
		Header:    header,
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
