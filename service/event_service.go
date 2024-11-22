package service

import (
	"bytes"
	"encoding/base64"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/sdk/events"
	"go-transfers/client"
	"log/slog"
)

type EventService struct {
	client *client.EventClient
}

func NewEventService(client *client.EventClient) *EventService {
	return &EventService{client: client}
}

func (es *EventService) ReadEvents(from uint32, to uint32) error {

	for i := from; i < to; i++ {
		tickEvents, err := es.client.GetEvents(i)
		if err != nil {
			slog.Error("Error getting events for tick.", "Tick", i)
			return err
		}
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
					decodedEvent, decodeErr := decodeEvent(eventType, eventData)
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

func decodeEvent(eventType uint8, eventData []byte) (*eventspb.DecodedEvent, error) {
	switch eventType {
	case events.EventTypeQuTransfer:
		var event events.QuTransferEvent
		err := event.UnmarshalBinary(eventData)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling qu transfer event")
		}

		sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting source identity pubkey")
		}

		destID, err := common.PubKeyToIdentity(event.DestinationIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting destination identity pubkey")
		}

		pbEvent := eventspb.DecodedEvent_QuTransferEvent_{
			QuTransferEvent: &eventspb.DecodedEvent_QuTransferEvent{
				SourceId: sourceID.String(),
				DestId:   destID.String(),
				Amount:   event.Amount,
			},
		}
		return &eventspb.DecodedEvent{Event: &pbEvent}, nil
	case events.EventTypeAssetIssuance:
		var event events.AssetIssuanceEvent
		err := event.UnmarshalBinary(eventData)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling asset issuance event")
		}

		sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting source identity pubkey")
		}

		pbEvent := eventspb.DecodedEvent_AssetIssuanceEvent_{
			AssetIssuanceEvent: &eventspb.DecodedEvent_AssetIssuanceEvent{
				SourceId:         sourceID.String(),
				AssetName:        string(bytes.TrimRight(event.AssetName[:], "\x00")),
				NumberOfDecimals: uint32(event.NumberOfDecimals),
				MeasurementUnit:  event.MeasurementUnit[:],
				NumberOfShares:   event.NumberOfShares,
			},
		}

		return &eventspb.DecodedEvent{Event: &pbEvent}, nil
	case events.EventTypeAssetOwnershipChange:
		var event events.AssetOwnershipChangeEvent
		err := event.UnmarshalBinary(eventData)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling asset ownership change event")
		}

		sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting source identity pubkey")
		}

		destID, err := common.PubKeyToIdentity(event.DestinationIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting destination identity pubkey")
		}

		issuerID, err := common.PubKeyToIdentity(event.IssuerIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting issuer identity pubkey")
		}

		pbEvent := eventspb.DecodedEvent_AssetOwnershipChangeEvent_{
			AssetOwnershipChangeEvent: &eventspb.DecodedEvent_AssetOwnershipChangeEvent{
				SourceId:         sourceID.String(),
				DestId:           destID.String(),
				IssuerId:         issuerID.String(),
				AssetName:        string(bytes.TrimRight(event.AssetName[:], "\x00")),
				NumberOfDecimals: uint32(event.NumberOfDecimals),
				MeasurementUnit:  event.MeasurementUnit[:],
				NumberOfShares:   event.NumberOfShares,
			},
		}

		return &eventspb.DecodedEvent{Event: &pbEvent}, nil
	case events.EventTypeAssetPossessionChange:
		var event events.AssetPossessionChangeEvent
		err := event.UnmarshalBinary(eventData)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshalling asset possession change event")
		}

		sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting source identity pubkey")
		}

		destID, err := common.PubKeyToIdentity(event.DestinationIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting destination identity pubkey")
		}

		issuerID, err := common.PubKeyToIdentity(event.IssuerIdentityPubKey)
		if err != nil {
			return nil, errors.Wrap(err, "converting issuer identity pubkey")
		}

		pbEvent := eventspb.DecodedEvent_AssetPossessionChangeEvent_{
			AssetPossessionChangeEvent: &eventspb.DecodedEvent_AssetPossessionChangeEvent{
				SourceId:         sourceID.String(),
				DestId:           destID.String(),
				IssuerId:         issuerID.String(),
				AssetName:        convertAssetName(event),
				NumberOfDecimals: uint32(event.NumberOfDecimals),
				MeasurementUnit:  event.MeasurementUnit[:],
				NumberOfShares:   event.NumberOfShares,
			},
		}

		return &eventspb.DecodedEvent{Event: &pbEvent}, nil
	default:
		return nil, errors.Errorf("not supported event type: %d", eventType)
	}
}

func convertAssetName(event events.AssetPossessionChangeEvent) string {
	return string(bytes.TrimRight(event.AssetName[:], "\x00"))
}
