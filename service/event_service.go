package service

import (
	"encoding/base64"
	"github.com/qubic/go-qubic/sdk/events"
	"go-transfers/client"
	"log/slog"
)

type EventService struct {
	client       client.EventClient
	eventDecoder EventDecoder
}

func NewEventService(client client.EventClient) *EventService {
	return &EventService{client: client, eventDecoder: EventDecoder{}}
}

func (es *EventService) ProcessTickEvents(from uint32, toExcl uint32) error {

	slog.Info("Processing tick range.", "From", from, "To", toExcl)
	for i := from; i < toExcl; i++ {
		tickEvents, err := es.client.GetEvents(i)
		if err != nil {
			slog.Error("Error getting events for tick.", "Tick", i)
			return err
		}
		slog.Debug("Processing events.", "TickEvents", tickEvents)
		txEvs := tickEvents.TxEvents
		for _, transactionEvents := range txEvs {
			for _, event := range transactionEvents.Events {
				eventData, base64Err := base64.StdEncoding.DecodeString(event.EventData)
				if base64Err != nil {
					slog.Error("Could not base64 decode event data.", "eventData", event.GetEventData(), "error", base64Err)
					return base64Err
				}

				eventType := uint8(event.GetEventType())
				switch eventType {
				case events.EventTypeQuTransfer, events.EventTypeAssetOwnershipChange, events.EventTypeAssetPossessionChange:
					slog.Info("Processing event data.", "transaction", transactionEvents.TxId, "eventData", event.EventData)
					decodedEvent, decodeErr := es.eventDecoder.DecodeEvent(eventType, eventData)
					if decodeErr != nil {
						return decodeErr
					}
					slog.Debug("Decoded event.", "event", decodedEvent)
				default:
					slog.Info("Ignoring unhandled event type.", "Transaction", transactionEvents.TxId, "eventType", eventType, "eventData", event.EventData)
				}
			}
		}
	}

	return nil

}
