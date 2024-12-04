package service

import (
	"encoding/base64"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"github.com/qubic/go-qubic/sdk/events"
	"log/slog"
)

type EventProcessor struct {
	repository Repository
}

func NewEventProcessor(repository Repository) *EventProcessor {
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

			count += len(relevantEvents)
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
					slog.Error("Could not base64 decode event data.", "eventData", event.GetEventData(), "error", err)
					return -1, errors.Wrap(err, "decoding event data.")
				}
				eventType := uint8(event.EventType)
				if eventType == events.EventTypeQuTransfer {
					err = ep.storeQuTransferEvent(eventData, eventId)
				} else if eventType == events.EventTypeAssetOwnershipChange {
					err = ep.storeAssetOwnershipChangeEvent(eventData, eventId)
				} else if eventType == events.EventTypeAssetPossessionChange {
					err = ep.storeAssetPossessionChangeEvent(eventData, eventId)
				} else if eventType == events.EventTypeAssetIssuance {
					err = ep.storeAssetIssuanceEvent(eventData, eventId)
				} else {
					err = errors.New("unexpected unhandled event type.")
				}
				if err != nil {
					slog.Error("Could not process event.", "eventType", event.EventType, "eventData", eventData, "error", err)
					return -1, errors.Wrap(err, "storing event details")
				}
			}
		}
	}
	return count, nil
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

func (ep *EventProcessor) storeAssetIssuanceEvent(eventData []byte, eventId int) error {
	decodedEvent, err := DecodeAssetIssuanceEvent(eventData)
	if err != nil {
		return errors.Wrap(err, "decoding event")
	}
	assetIssuanceEvent := decodedEvent.GetAssetIssuanceEvent()
	assetId, err := ep.repository.GetOrCreateAsset(assetIssuanceEvent.GetSourceId(), assetIssuanceEvent.GetAssetName())
	if err != nil {
		return err
	}
	assetIssuanceEventId, err := ep.repository.GetOrCreateAssetIssuanceEvent(eventId, assetId,
		assetIssuanceEvent.GetNumberOfShares(),
		assetIssuanceEvent.GetMeasurementUnit(),
		assetIssuanceEvent.GetNumberOfDecimals())
	if err != nil {
		return err
	} else {
		slog.Debug("Stored asset issuance event.", "id", assetIssuanceEventId)
	}
	return nil
}

func (ep *EventProcessor) storeAssetPossessionChangeEvent(eventData []byte, eventId int) error {
	decodedEvent, err := DecodeAssetPossessionChangeEvent(eventData)
	if err != nil {
		return errors.Wrap(err, "decoding event")
	}
	assetChangeEvent := decodedEvent.GetAssetPossessionChangeEvent()
	sourceId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetSourceId())
	if err != nil {
		return err
	}
	destinationId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetDestId())
	if err != nil {
		return err
	}
	assetId, err := ep.repository.GetOrCreateAsset(assetChangeEvent.GetIssuerId(), assetChangeEvent.GetAssetName())
	if err != nil {
		return err
	}
	assetChangeEventId, err := ep.repository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceId, destinationId, assetChangeEvent.GetNumberOfShares())
	if err != nil {
		return err
	} else {
		slog.Debug("Stored asset possession change event.", "id", assetChangeEventId)
	}
	return nil
}

func (ep *EventProcessor) storeAssetOwnershipChangeEvent(eventData []byte, eventId int) error {
	decodedEvent, err := DecodeAssetOwnershipChangeEvent(eventData)
	if err != nil {
		return errors.Wrap(err, "decoding event")
	}
	assetChangeEvent := decodedEvent.GetAssetOwnershipChangeEvent()
	sourceId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetSourceId())
	if err != nil {
		return err
	}
	destinationId, err := ep.repository.GetOrCreateEntity(assetChangeEvent.GetDestId())
	if err != nil {
		return err
	}
	assetId, err := ep.repository.GetOrCreateAsset(assetChangeEvent.GetIssuerId(), assetChangeEvent.GetAssetName())
	if err != nil {
		return err
	}
	assetChangeEventId, err := ep.repository.GetOrCreateAssetChangeEvent(eventId, assetId, sourceId, destinationId, assetChangeEvent.GetNumberOfShares())
	if err != nil {
		return err
	} else {
		slog.Debug("Stored asset ownership change event.", "id", assetChangeEventId)
	}
	return nil
}

func (ep *EventProcessor) storeQuTransferEvent(eventData []byte, eventId int) error {
	decodedEvent, err := DecodeQuTransferEvent(eventData)
	if err != nil {
		return errors.Wrap(err, "decoding event")
	}
	transferEvent := decodedEvent.GetQuTransferEvent()
	sourceId, err := ep.repository.GetOrCreateEntity(transferEvent.GetSourceId())
	if err != nil {
		return errors.Wrap(err, "getting source entity")
	}
	destinationId, err := ep.repository.GetOrCreateEntity(transferEvent.GetDestId())
	if err != nil {
		return errors.Wrap(err, "getting destination entity")
	}
	transferId, err := ep.repository.GetOrCreateQuTransferEvent(eventId, sourceId, destinationId, transferEvent.GetAmount())
	if err != nil {
		return errors.Wrap(err, "creating qu transfer event")
	} else {
		slog.Debug("Stored qu transfer event.", "id", transferId)
	}
	return nil
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
