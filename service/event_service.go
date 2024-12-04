package service

import (
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"go-transfers/client"
	"log/slog"
	"math"
	"time"
)

type EventClient interface {
	GetEvents(tickNumber uint32) (*eventspb.TickEvents, error)
	GetStatus() (*client.EventStatus, error)
	GetTickInfo() (*client.TickInfo, error)
}

type Repository interface {
	GetOrCreateEntity(identity string) (int, error)
	GetOrCreateAsset(issuer, name string) (int, error)
	GetOrCreateTick(tickNumber uint32) (int, error)
	GetOrCreateTransaction(hash string, tickId int) (int, error)
	GetOrCreateEvent(transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error)
	GetOrCreateQuTransferEvent(eventId int, sourceEntityId int, destinationEntityId int, amount uint64) (int, error)
	GetOrCreateAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error)
	GetOrCreateAssetIssuanceEvent(eventId int, assetId int, numberOfShares int64, unitOfMeasurement []byte, numberOfDecimalPlaces uint32) (int, error)
	Close()
}

type EventService struct {
	client         EventClient
	eventProcessor *EventProcessor
	repository     Repository
}

func NewEventService(client EventClient, eventProcessor *EventProcessor) *EventService {
	es := EventService{
		client:         client,
		eventProcessor: eventProcessor,
	}
	return &es
}

var processedTick uint64 = 17563100

func (es *EventService) SyncInLoop() {
	loopTick := time.Tick(time.Second * 3)
	for range loopTick {
		err := es.sync()
		time.Sleep(time.Second)
		if err != nil {
			slog.Error("processing tick events", "err", err.Error())
		}
	}
}

func (es *EventService) sync() error {
	tickInfo, err := es.client.GetTickInfo()
	if err != nil {
		return errors.Wrap(err, "getting tick info")
	}
	if uint64(tickInfo.InitialTick) > processedTick {
		slog.Info("initial tick > processed tick", "initial", tickInfo.InitialTick, "processed", processedTick)
	}
	startTick := uint64(math.Max(float64(processedTick+1), float64(tickInfo.InitialTick)))

	status, err := es.client.GetStatus()
	if err != nil {
		return errors.Wrap(err, "getting event status.")
	}
	endTick := uint64(math.Min(float64(status.AvailableTick), float64(tickInfo.CurrentTick)))
	endTick = uint64(math.Min(float64(endTick), float64(startTick+100))) // max batch process 100 ticks per run

	slog.Info("Status:", "processed", processedTick, "current", tickInfo.CurrentTick, "available", status.AvailableTick)
	if startTick <= endTick { // ok
		slog.Info("Syncing:", "from", startTick, "to", endTick)
		err := es.ProcessTickEvents(startTick, endTick+1)
		if err != nil {
			return errors.Wrap(err, "processing tick events")
		}
		processedTick = endTick
	}
	return nil
}

func (es *EventService) ProcessTickEvents(from uint64, toExcl uint64) error {
	for i := from; i < toExcl; i++ {

		if i > math.MaxInt32 {
			return errors.New("uint32 overflow")
		}

		tickEvents, err := es.client.GetEvents(uint32(i)) // attention. need to cast here.
		if err != nil {
			slog.Error("Error getting events for tick.", "Tick", i)
			return errors.Wrap(err, "Error getting events for tick.")
		}

		eventCount, err := es.eventProcessor.ProcessTickEvents(tickEvents)
		if err != nil {
			return errors.Wrap(err, "processing tick events.")
		}

		slog.Info("Processed tick.", "tick", i, "events", eventCount)

	}
	return nil

}
