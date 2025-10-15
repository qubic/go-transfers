package sync

import (
	"context"
	"go-transfers/client"
	"math"
	"time"

	"github.com/gookit/slog"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
)

type EventClient interface {
	GetEvents(ctx context.Context, tickNumber uint32) (*eventspb.TickEvents, error)
	GetStatus(ctx context.Context) (*client.EventStatus, error)
	GetTickInfo(ctx context.Context) (*client.TickInfo, error)
}

type TickNumberRepository interface {
	GetLatestTick(ctx context.Context) (int, error)
	UpdateLatestTick(ctx context.Context, tickNumber int) error
}

type Metrics interface {
	SetLatestProcessedTick(tick uint32)
	SetLatestEventTick(tick uint32)
	SetLatestLiveTick(tick uint32)
}

type EventService struct {
	client         EventClient
	eventProcessor *EventProcessor
	repository     TickNumberRepository
	metrics        Metrics
}

func NewEventService(c EventClient, ep *EventProcessor, r TickNumberRepository, m Metrics) (*EventService, error) {
	es := EventService{
		client:         c,
		eventProcessor: ep,
		repository:     r,
		metrics:        m,
	}
	return &es, nil
}

func (es *EventService) SyncInLoop() {
	var count uint64
	loopTick := time.Tick(time.Second)
	for range loopTick {
		err := es.sync(count)
		count++
		if err != nil {
			slog.Error("processing tick events", "err", err.Error())
		}
		time.Sleep(time.Second)
	}
}

func (es *EventService) sync(count uint64) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	startTick, currentTick, err := es.calculateStartTick(ctx)
	if err != nil {
		return errors.Wrap(err, "calculating start tick")
	}
	es.metrics.SetLatestLiveTick(uint32(currentTick))

	status, err := es.client.GetStatus(ctx)
	if err != nil {
		return errors.Wrap(err, "getting status from event service")
	}
	es.metrics.SetLatestEventTick(status.AvailableTick)
	endTick := int(math.Min(float64(status.AvailableTick), float64(currentTick)))
	endTick = int(math.Min(float64(endTick), float64(startTick+100))) // max batch process 100 ticks per run

	if count%500 == 0 { // log status in regular intervals
		slog.Info("Status:", "next", startTick, "current", currentTick, "available", status.AvailableTick)
	}

	if startTick > endTick {
		return nil
	}

	slog.Debug("Syncing:", "from", startTick, "to", endTick)
	err = es.processTickEventsRange(ctx, startTick, endTick+1) // end tick exclusive
	if err != nil {
		return errors.Wrap(err, "processing tick events")
	}
	return nil
}

func (es *EventService) calculateStartTick(ctx context.Context) (int, int, error) {
	processedTick, err := es.repository.GetLatestTick(ctx)
	if err != nil {
		slog.Error(err.Error())
		return -1, -1, errors.Wrap(err, "getting processed tick")
	}

	tickInfo, err := es.client.GetTickInfo(ctx)
	if err != nil {
		return -1, -1, errors.Wrap(err, "getting tick info")
	}

	if int(tickInfo.InitialTick) > processedTick {
		slog.Info("initial tick > processed tick", "initial", tickInfo.InitialTick, "processed", processedTick)
	}
	return int(math.Max(float64(processedTick+1), float64(tickInfo.InitialTick))), int(tickInfo.CurrentTick), nil
}

func (es *EventService) processTickEventsRange(ctx context.Context, from, toExcl int) error {
	for tick := from; tick < toExcl; tick++ {
		err := es.processTickEvents(ctx, tick)
		if err != nil {
			return errors.Wrapf(err, "processing tick events from [%d] to [%d]", from, toExcl)
		}
	}
	return nil
}

func (es *EventService) processTickEvents(ctx context.Context, tick int) error {

	if tick > math.MaxInt32 {
		return errors.New("uint32 overflow")
	}

	tickEvents, err := es.client.GetEvents(ctx, uint32(tick)) // attention. need to cast here.
	if err != nil {
		slog.Warn("Error getting events.", "tick", tick)
		return errors.Wrapf(err, "getting events for tick [%d]", tick)
	}

	eventCount, err := es.eventProcessor.ProcessTickEvents(ctx, tickEvents)
	if err != nil {
		return errors.Wrapf(err, "processing events for tick [%d]", tick)
	}

	err = es.repository.UpdateLatestTick(ctx, tick)
	if err != nil {
		return errors.Wrapf(err, "updating latest tick to [%d]", tick)
	}
	es.metrics.SetLatestProcessedTick(uint32(tick))

	var numberOfTransactionEvents, numberOfTotalEvents int
	for _, txEv := range tickEvents.TxEvents {
		numberOfTotalEvents += len(txEv.Events)
		numberOfTransactionEvents++
	}

	slog.Info("Processed:", "tick", tick, "stored", eventCount, "transactions", numberOfTransactionEvents, "events", numberOfTotalEvents)
	return nil
}
