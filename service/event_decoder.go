package service

import (
	"bytes"
	"github.com/pkg/errors"
	eventspb "github.com/qubic/go-events/proto"
	"github.com/qubic/go-qubic/common"
	"github.com/qubic/go-qubic/sdk/events"
)

func DecodeQuTransferEvent(eventData []byte) (*eventspb.DecodedEvent, error) {
	var event events.QuTransferEvent
	err := event.UnmarshalBinary(eventData)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling qu transfer event")
	}

	sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting source identity public key")
	}

	destID, err := common.PubKeyToIdentity(event.DestinationIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting destination identity public key")
	}

	pbEvent := eventspb.DecodedEvent_QuTransferEvent_{
		QuTransferEvent: &eventspb.DecodedEvent_QuTransferEvent{
			SourceId: sourceID.String(),
			DestId:   destID.String(),
			Amount:   event.Amount,
		},
	}
	return &eventspb.DecodedEvent{Event: &pbEvent}, nil
}

func DecodeAssetIssuanceEvent(eventData []byte) (*eventspb.DecodedEvent, error) {
	var event events.AssetIssuanceEvent
	err := event.UnmarshalBinary(eventData)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling asset issuance event")
	}

	sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting source identity public key")
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
}

func DecodeAssetOwnershipChangeEvent(eventData []byte) (*eventspb.DecodedEvent, error) {
	var event events.AssetOwnershipChangeEvent
	err := event.UnmarshalBinary(eventData)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling asset ownership change event")
	}

	sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting source identity public key")
	}

	destID, err := common.PubKeyToIdentity(event.DestinationIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting destination identity public key")
	}

	issuerID, err := common.PubKeyToIdentity(event.IssuerIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting issuer identity public key")
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
}

func DecodeAssetPossessionChangeEvent(eventData []byte) (*eventspb.DecodedEvent, error) {
	var event events.AssetPossessionChangeEvent
	err := event.UnmarshalBinary(eventData)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling asset possession change event")
	}

	sourceID, err := common.PubKeyToIdentity(event.SourceIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting source identity public key")
	}

	destID, err := common.PubKeyToIdentity(event.DestinationIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting destination identity public key")
	}

	issuerID, err := common.PubKeyToIdentity(event.IssuerIdentityPubKey)
	if err != nil {
		return nil, errors.Wrap(err, "converting issuer identity public key")
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
}

func convertAssetName(event events.AssetPossessionChangeEvent) string {
	return string(bytes.TrimRight(event.AssetName[:], "\x00"))
}
