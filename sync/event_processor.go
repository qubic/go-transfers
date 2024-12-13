package sync

import (
	"encoding/base64"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"github.com/qubic/go-qubic/sdk/events"
	"strings"
)

const AAA = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

type EventRepository interface {
	GetOrCreateEntity(identity string) (int, error)
	GetOrCreateAsset(issuer, name string) (int, error)
	GetOrCreateTick(tickNumber uint32) (int, error)
	GetOrCreateTransaction(hash string, tickId int) (int, error)
	GetOrCreateEvent(transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error)
	GetOrCreateQuTransferEvent(eventId int, sourceEntityId int, destinationEntityId int, amount uint64) (int, error)
	GetOrCreateAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error)
	GetOrCreateAssetIssuanceEvent(eventId int, assetId int, numberOfShares int64, unitOfMeasurement []byte, numberOfDecimalPlaces uint32) (int, error)
}

type EventProcessor struct {
	repository EventRepository
}

func NewEventProcessor(repository EventRepository) *EventProcessor {
	ep := EventProcessor{
		repository: repository,
	}
	return &ep
}

func (ep *EventProcessor) ProcessTickEvents(tickEvents *eventspb.TickEvents) (int, error) {

	var count int
	for _, transactionEvents := range tickEvents.TxEvents {

		slog.Debug("Processing transaction events.", "transaction_events", transactionEvents)
		relevantEvents := filterRelevantEvents(transactionEvents.Events)
		if len(relevantEvents) > 0 {

			slog.Debug("Processing events of transaction.", "hash", transactionEvents.TxId, "count", len(relevantEvents))

			transactionId, err := ep.storeTransaction(tickEvents.GetTick(), transactionEvents.GetTxId())
			if err != nil {
				return -1, errors.Wrap(err, "storing transaction")
			}

			for _, event := range relevantEvents {
				slog.Debug("Processing event.", "event", event)
				eventId, err := ep.getOrCreateEvent(event, transactionId)
				if err != nil {
					return -1, errors.Wrap(err, "storing event")
				}
				eventData, err := base64.StdEncoding.DecodeString(event.EventData)
				if err != nil {
					return -1, errors.Wrap(err, "base64 decoding event data.")
				}
				var dbId int
				eventType := uint8(event.EventType)
				if eventType == events.EventTypeQuTransfer {
					dbId, err = ep.storeQuTransferEvent(eventData, eventId)
				} else if eventType == events.EventTypeAssetOwnershipChange {
					dbId, err = ep.storeAssetOwnershipChangeEvent(eventData, eventId)
				} else if eventType == events.EventTypeAssetPossessionChange {
					dbId, err = ep.storeAssetPossessionChangeEvent(eventData, eventId)
				} else if eventType == events.EventTypeAssetIssuance {
					dbId, err = ep.storeAssetIssuanceEvent(eventData, eventId)
				} else {
					err = errors.New("unexpected unhandled event type.")
				}
				if err != nil {
					slog.Error("Could not process event.", "eventType", event.EventType, "eventData", eventData, "error", err)
					return -1, errors.Wrap(err, "storing event details")
				} else {
					slog.Info("Stored event:", "id", dbId, "type", eventType, "transaction", transactionId)
					count++
				}
			}
		}
	}
	return count, nil
}

func (ep *EventProcessor) getTransactionId(tickNumber uint32, hash string) (int, error) {
	transactionId, err := ep.storeTransaction(tickNumber, hash)
	if err != nil {
		return -1, errors.Wrap(err, "storing transaction")
	}
	return transactionId, nil
}

func (ep *EventProcessor) storeTransaction(tick uint32, transactionHash string) (int, error) {
	tickId, err := ep.repository.GetOrCreateTick(tick)
	if err != nil {
		return -1, errors.Wrap(err, "storing tick")
	}
	transactionId, err := ep.repository.GetOrCreateTransaction(transactionHash, tickId)
	if err != nil {
		return -1, errors.Wrap(err, "storing transaction")
	}
	return transactionId, nil
}

func (ep *EventProcessor) getOrCreateEvent(event *eventspb.Event, transactionId int) (int, error) {
	eventEventId, err := ep.getEventId(event)
	if err != nil {
		return -1, errors.Wrap(err, "extracting event id")
	}
	eventId, err := ep.repository.GetOrCreateEvent(transactionId, eventEventId, event.EventType, event.EventData)
	if err != nil {
		return -1, errors.Wrap(err, "get  or creating event")
	}
	return eventId, nil
}

func (ep *EventProcessor) getEventId(event *eventspb.Event) (uint64, error) {
	if event.GetHeader() != nil {
		return event.Header.EventId, nil
	} else {
		slog.Error("Event header not found.", "event", event)
		return 0, errors.New("No event header found.")
	}

}

func (ep *EventProcessor) storeAssetIssuanceEvent(eventData []byte, eventId int) (int, error) {
	decodedEvent, err := DecodeAssetIssuanceEvent(eventData)
	if err != nil {
		return -1, errors.Wrap(err, "decoding asset issuance")
	}
	assetIssuanceEvent := decodedEvent.GetAssetIssuanceEvent()
	assetId, err := ep.repository.GetOrCreateAsset(assetIssuanceEvent.GetSourceId(), assetIssuanceEvent.GetAssetName())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset issuance")
	}
	assetIssuanceEventId, err := ep.repository.GetOrCreateAssetIssuanceEvent(eventId, assetId,
		assetIssuanceEvent.GetNumberOfShares(),
		assetIssuanceEvent.GetMeasurementUnit(),
		assetIssuanceEvent.GetNumberOfDecimals())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset issuance")
	} else {
		slog.Debug("Stored asset issuance event.", "id", assetIssuanceEventId)
	}
	return assetIssuanceEventId, nil
}

func (ep *EventProcessor) storeAssetPossessionChangeEvent(eventData []byte, eventId int) (int, error) {
	decodedEvent, err := DecodeAssetPossessionChangeEvent(eventData)
	if err != nil {
		return -1, errors.Wrap(err, "decoding asset possession change")
	}
	assetChangeEvent := decodedEvent.GetAssetPossessionChangeEvent()
	sourceId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetSourceId())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset possession change")
	}
	destinationId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetDestId())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset possession change")
	}
	assetId, err := ep.repository.GetOrCreateAsset(assetChangeEvent.GetIssuerId(), assetChangeEvent.GetAssetName())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset possession change")
	}
	assetChangeEventId, err := ep.repository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceId, destinationId, assetChangeEvent.GetNumberOfShares())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset possession change")
	} else {
		slog.Debug("Stored asset possession change event.", "id", assetChangeEventId)
	}
	return assetChangeEventId, nil
}

func (ep *EventProcessor) storeAssetOwnershipChangeEvent(eventData []byte, eventId int) (int, error) {
	decodedEvent, err := DecodeAssetOwnershipChangeEvent(eventData)
	if err != nil {
		return -1, errors.Wrap(err, "decoding asset ownership change")
	}
	assetChangeEvent := decodedEvent.GetAssetOwnershipChangeEvent()
	sourceId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetSourceId())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset ownership change")
	}
	destinationId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetDestId())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset ownership change")
	}
	assetId, err := ep.repository.GetOrCreateAsset(assetChangeEvent.GetIssuerId(), assetChangeEvent.GetAssetName())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset ownership change")
	}
	assetChangeEventId, err := ep.repository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceId, destinationId, assetChangeEvent.GetNumberOfShares())
	if err != nil {
		return -1, errors.Wrap(err, "storing asset ownership change")
	} else {
		slog.Debug("Stored asset ownership change event.", "id", assetChangeEventId)
	}
	return assetChangeEventId, nil
}

func (ep *EventProcessor) storeQuTransferEvent(eventData []byte, eventId int) (int, error) {
	decodedEvent, err := DecodeQuTransferEvent(eventData)
	if err != nil {
		return -1, errors.Wrap(err, "decoding qu transfer")
	}
	transferEvent := decodedEvent.GetQuTransferEvent()

	sourceId, err := ep.repository.GetOrCreateEntity(transferEvent.GetSourceId())
	if err != nil {
		return -1, errors.Wrap(err, "storing qu transfer")
	}
	destinationId, err := ep.repository.GetOrCreateEntity(transferEvent.GetDestId())
	if err != nil {
		return -1, errors.Wrap(err, "storing qu transfer")
	}

	transferId, err := ep.repository.GetOrCreateQuTransferEvent(eventId, sourceId, destinationId, transferEvent.GetAmount())
	if err != nil {
		return -1, errors.Wrap(err, "storing qu transfer")
	} else {
		slog.Debug("Stored qu transfer event.", "id", transferId)
	}
	return transferId, nil
}

func filterRelevantEvents(events []*eventspb.Event) []*eventspb.Event {
	var filtered []*eventspb.Event
	for _, ev := range events {
		if isRelevantEvent(ev) {
			filtered = append(filtered, ev)
		}
	}
	return filtered
}

func isRelevantEvent(ev *eventspb.Event) bool {
	eventType := uint8(ev.EventType)
	// this is a bit awkward. As we don't have the transaction data we need to look into the event data
	// for checking, if it is relevant. Same decoding will happen once more later. This seems to be easier
	// than to change the transaction creation logic (we need to persist the transaction first).
	if eventType == events.EventTypeQuTransfer {
		eventData, err := base64.StdEncoding.DecodeString(ev.EventData)
		if err != nil {
			slog.Error("Error decoding event data", "data", ev.EventData, "error", err)
			return false
		}
		decodedEvent, err := DecodeQuTransferEvent(eventData)
		if err != nil {
			slog.Error("Error decoding qu transfer event", "data", ev.EventData, "error", err)
			return false
		}
		transferEvent := decodedEvent.GetQuTransferEvent()
		// ignore qu transfers that go to AAA or come from AAA (mainly mining deposits but could be burning, too)
		return !(strings.HasPrefix(transferEvent.GetSourceId(), AAA) || strings.HasPrefix(transferEvent.GetDestId(), AAA))
	}
	return eventType == events.EventTypeAssetPossessionChange ||
		eventType == events.EventTypeAssetOwnershipChange ||
		eventType == events.EventTypeAssetIssuance
}
