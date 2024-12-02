package service

import (
	"encoding/base64"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"github.com/qubic/go-qubic/sdk/events"
	"log/slog"
)

type EventClient interface {
	GetEvents(tickNumber uint32) (*eventspb.TickEvents, error)
}

type Repository interface {
	GetOrCreateTick(tickNumber uint32) (int, error)
	GetOrCreateTransaction(hash string, tickId int) (int, error)
	GetOrCreateEvent(transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error)
	Close()
}

type EventService struct {
	client          EventClient
	eventRepository Repository
}

func NewEventService(client EventClient, eventRepository Repository) *EventService {
	return &EventService{
		client:          client,
		eventRepository: eventRepository,
	}
}

func (es *EventService) ProcessTickEvents(from uint32, toExcl uint32) error {

	slog.Info("Processing tick range.", "From", from, "To", toExcl)
	for i := from; i < toExcl; i++ {
		err := es.processTickEvents(i)
		if err != nil {
			return err
		}
	}
	return nil

}

func (es *EventService) processTickEvents(tickNumber uint32) error {
	tickEvents, err := es.client.GetEvents(tickNumber)
	if err != nil {
		slog.Error("Error getting events for tick.", "Tick", tickNumber)
		return err
	}
	slog.Debug("Processing events.", "TickEvents", tickEvents)
	txEvs := tickEvents.TxEvents

	for _, transactionEvents := range txEvs {

		relevantEvents := filterRelevantEvents(transactionEvents.Events)
		if len(relevantEvents) > 0 {

			tickId, err := es.eventRepository.GetOrCreateTick(tickEvents.GetTick())
			if err != nil {
				return err
			}
			transactionId, err := es.eventRepository.GetOrCreateTransaction(transactionEvents.GetTxId(), tickId)
			if err != nil {
				return err
			}
			for _, event := range relevantEvents {
				eventEventId, err := getEventId(event)
				if err != nil {
					return err
				}
				eventId, err := es.eventRepository.GetOrCreateEvent(transactionId, eventEventId, event.EventType, event.EventData)
				if err != nil {
					return err
				}
				_ = eventId // FIXME

				eventData, err := base64.StdEncoding.DecodeString(event.EventData)
				if err != nil {
					slog.Error("Could not base64 decode event data.", "eventData", event.GetEventData(), "error", err)
					return err
				}
				eventType := uint8(event.EventType)
				switch eventType {
				case events.EventTypeQuTransfer:
					decodedEvent, err := DecodeQuTransferEvent(eventData)
					if err != nil {
						slog.Error("Could not decode qu transfer event.", "eventType", event.EventType, "eventData", eventData, "error", err)
						return err
					}
					transferEvent := decodedEvent.GetQuTransferEvent()
					transferEvent.GetSourceId()
					transferEvent.GetDestId()
					transferEvent.GetAmount()
					// TODO store qu transfer event

				case events.EventTypeAssetOwnershipChange:
					// TODO store asset ownership change event
				case events.EventTypeAssetPossessionChange:
					// TODO store asset possession change event
				case events.EventTypeAssetIssuance:
					// TODO store asset issuance event
				default:
					slog.Error("unexpected unhandled event type.",
						"Transaction", transactionEvents.TxId,
						"eventType", eventType,
						"eventData", event.EventData)
				}

			}
		}
	}
	return nil
}

func getEventId(event *eventspb.Event) (uint64, error) {
	if event.GetHeader() != nil {
		return event.Header.EventId, nil
	} else {
		slog.Error("Event header not found.", "event", event)
		return 0, errors.New("No event header found.")
	}

}

func filterRelevantEvents(events []*eventspb.Event) []*eventspb.Event {
	var filtered []*eventspb.Event
	for _, ev := range events {
		if isRelevantType(ev) {
			filtered = append(filtered, ev)
		}
	}
	return filtered
}

func isRelevantType(ev *eventspb.Event) bool {
	eventType := uint8(ev.EventType)
	return eventType == events.EventTypeQuTransfer ||
		eventType == events.EventTypeAssetPossessionChange ||
		eventType == events.EventTypeAssetOwnershipChange ||
		eventType == events.EventTypeAssetIssuance
}
