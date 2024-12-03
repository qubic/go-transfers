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
	slog.Info("Processing tick events.", "Tick", tickNumber)

	txEvs := tickEvents.TxEvents

	for _, transactionEvents := range txEvs {

		slog.Debug("Processing transaction events.", "transaction_events", transactionEvents)
		relevantEvents := filterRelevantEvents(transactionEvents.Events)
		if len(relevantEvents) > 0 {

			slog.Info("Storing events.", "hash", transactionEvents.TxId, "count", len(relevantEvents))

			tickId, err := es.eventRepository.GetOrCreateTick(tickEvents.GetTick())
			if err != nil {
				return err
			}
			transactionId, err := es.eventRepository.GetOrCreateTransaction(transactionEvents.GetTxId(), tickId)
			if err != nil {
				return err
			}
			for _, event := range relevantEvents {
				slog.Debug("Processing event.", "event", event)
				eventEventId, err := getEventId(event)
				if err != nil {
					return err
				}

				eventId, err := es.eventRepository.GetOrCreateEvent(transactionId, eventEventId, event.EventType, event.EventData)
				if err != nil {
					return err
				}

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
						slog.Error("Could not decode event.", "eventType", event.EventType, "eventData", eventData, "error", err)
						return err
					}
					transferEvent := decodedEvent.GetQuTransferEvent()
					sourceId, err := es.eventRepository.GetOrCreateEntity(transferEvent.GetSourceId())
					if err != nil {
						return err
					}
					destinationId, err := es.eventRepository.GetOrCreateEntity(transferEvent.GetDestId())
					if err != nil {
						return err
					}
					transferId, err := es.eventRepository.GetOrCreateQuTransferEvent(eventId, sourceId, destinationId, transferEvent.GetAmount())
					if err != nil {
						return err
					} else {
						slog.Debug("Stored qu transfer event.", "id", transferId)
					}
				case events.EventTypeAssetOwnershipChange:
					decodedEvent, err := DecodeAssetOwnershipChangeEvent(eventData)
					if err != nil {
						slog.Error("Could not decode event.", "eventType", event.EventType, "eventData", eventData, "error", err)
						return err
					}
					assetChangeEvent := decodedEvent.GetAssetOwnershipChangeEvent()
					sourceId, err := es.eventRepository.GetOrCreateEntity(assetChangeEvent.GetSourceId())
					if err != nil {
						return err
					}
					destinationId, err := es.eventRepository.GetOrCreateEntity(assetChangeEvent.GetDestId())
					if err != nil {
						return err
					}
					assetId, err := es.eventRepository.GetOrCreateAsset(assetChangeEvent.GetIssuerId(), assetChangeEvent.GetAssetName())
					if err != nil {
						return err
					}
					assetChangeEventId, err := es.eventRepository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceId, destinationId, assetChangeEvent.GetNumberOfShares())
					if err != nil {
						return err
					} else {
						slog.Debug("Stored asset ownership change event.", "id", assetChangeEventId, "eventType", event.EventType)
					}
				case events.EventTypeAssetPossessionChange:
					decodedEvent, err := DecodeAssetPossessionChangeEvent(eventData)
					if err != nil {
						slog.Error("Could not decode event.", "eventType", event.EventType, "eventData", eventData, "error", err)
						return err
					}
					assetChangeEvent := decodedEvent.GetAssetPossessionChangeEvent()
					sourceId, err := es.eventRepository.GetOrCreateEntity(assetChangeEvent.GetSourceId())
					if err != nil {
						return err
					}
					destinationId, err := es.eventRepository.GetOrCreateEntity(assetChangeEvent.GetDestId())
					if err != nil {
						return err
					}
					assetId, err := es.eventRepository.GetOrCreateAsset(assetChangeEvent.GetIssuerId(), assetChangeEvent.GetAssetName())
					if err != nil {
						return err
					}
					assetChangeEventId, err := es.eventRepository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceId, destinationId, assetChangeEvent.GetNumberOfShares())
					if err != nil {
						return err
					} else {
						slog.Debug("Stored asset possession change event.", "id", assetChangeEventId, "eventType", event.EventType)
					}
				case events.EventTypeAssetIssuance:
					decodedEvent, err := DecodeAssetIssuanceEvent(eventData)
					if err != nil {
						slog.Error("Could not decode event.", "eventType", event.EventType, "eventData", eventData, "error", err)
						return err
					}
					assetIssuanceEvent := decodedEvent.GetAssetIssuanceEvent()
					assetId, err := es.eventRepository.GetOrCreateAsset(assetIssuanceEvent.GetSourceId(), assetIssuanceEvent.GetAssetName())
					if err != nil {
						return err
					}
					assetIssuanceEventId, err := es.eventRepository.GetOrCreateAssetIssuanceEvent(eventId, assetId,
						assetIssuanceEvent.GetNumberOfShares(),
						assetIssuanceEvent.GetMeasurementUnit(),
						assetIssuanceEvent.GetNumberOfDecimals())
					if err != nil {
						return err
					} else {
						slog.Debug("Stored asset issuance event.", "id", assetIssuanceEventId, "eventType", event.EventType)
					}
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
