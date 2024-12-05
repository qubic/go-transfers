package sync

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

type TickNumberRepository interface {
	GetNumericValue(key string) (int, error)
	UpdateNumericValue(key string, value int) error
}

type EventService struct {
	client         EventClient
	eventProcessor *EventProcessor
	repository     TickNumberRepository
}

var processedTick = 0

func NewEventService(client EventClient, eventProcessor *EventProcessor, repository TickNumberRepository) (*EventService, error) {
	var err error
	processedTick, err = repository.GetNumericValue("tick")
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.Wrap(err, "getting processed tick value")
	}

	es := EventService{
		client:         client,
		eventProcessor: eventProcessor,
		repository:     repository,
	}
	return &es, nil
}

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
	if int(tickInfo.InitialTick) > processedTick {
		slog.Info("initial tick > processed tick", "initial", tickInfo.InitialTick, "processed", processedTick)
	}
	startTick := int(math.Max(float64(processedTick+1), float64(tickInfo.InitialTick)))

	status, err := es.client.GetStatus()
	if err != nil {
		return errors.Wrap(err, "getting event status.")
	}
	endTick := int(math.Min(float64(status.AvailableTick), float64(tickInfo.CurrentTick)))
	endTick = int(math.Min(float64(endTick), float64(startTick+100))) // max batch process 100 ticks per run

	slog.Info("Status:", "processed", processedTick, "current", tickInfo.CurrentTick, "available", status.AvailableTick)
	if startTick <= endTick { // ok
		slog.Debug("Syncing:", "from", startTick, "to", endTick)
		err := es.ProcessTickEvents(startTick, endTick+1)
		if err != nil {
			return errors.Wrap(err, "processing tick events")
		}
		processedTick = endTick
		err = es.repository.UpdateNumericValue("tick", endTick)
		if err != nil {
			return errors.Wrap(err, "updating processed tick")
		}
	}
	return nil
}

func (es *EventService) ProcessTickEvents(from int, toExcl int) error {
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
