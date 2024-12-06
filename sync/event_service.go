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
	GetLatestTick() (int, error)
	UpdateLatestTick(tickNumber int) error
}

type EventService struct {
	client         EventClient
	eventProcessor *EventProcessor
	repository     TickNumberRepository
}

func NewEventService(client EventClient, eventProcessor *EventProcessor, repository TickNumberRepository) (*EventService, error) {
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
	processedTick, err := es.repository.GetLatestTick()
	if err != nil {
		slog.Error(err.Error())
		return errors.Wrap(err, "getting processed tick")
	}

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

	slog.Debug("Status:", "processed", processedTick, "current", tickInfo.CurrentTick, "available", status.AvailableTick)
	if startTick <= endTick { // ok
		slog.Debug("Syncing:", "from", startTick, "to", endTick)
		tick, err := es.ProcessTickEvents(startTick, endTick+1)
		if err != nil {
			return errors.Wrap(err, "processing tick events")
		}
		err = es.repository.UpdateLatestTick(tick)
		if err != nil {
			return errors.Wrap(err, "updating processed tick")
		}
	}
	return nil
}

func (es *EventService) ProcessTickEvents(from int, toExcl int) (int, error) {
	tick := from
	for ; tick < toExcl; tick++ {

		if tick > math.MaxInt32 {
			return -1, errors.New("uint32 overflow")
		}

		tickEvents, err := es.client.GetEvents(uint32(tick)) // attention. need to cast here.
		if err != nil {
			slog.Error("Error getting events.", "tick", tick)
			return -1, errors.Wrap(err, "Error getting events for tick.")
		}

		eventCount, err := es.eventProcessor.ProcessTickEvents(tickEvents)
		if err != nil {
			return -1, errors.Wrap(err, "processing tick events.")
		}

		slog.Info("Processed:", "tick", tick, "events", eventCount)
	}
	return tick, nil
}
