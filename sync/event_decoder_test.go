package sync

import (
	"encoding/base64"
	"testing"
)

//goland:noinspection SpellCheckingInspection
func TestEventDecoder_Decode_QuTransferEvent(t *testing.T) {
	// asset transfer via qx contract (1000000 paid to BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARMID)
	eventData, err := base64.StdEncoding.DecodeString("sMmo18V9WMO9LstUtxvWC2ZfJc2/FZWKEUdAKOqNKDIBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBCDwAAAAAA")
	if err != nil {
		t.Error(err)
	}
	decoded, err := DecodeQuTransferEvent(eventData)
	if err != nil {
		t.Error(err)
	}

	if decoded.GetQuTransferEvent().GetSourceId() != "AKJDFZYITPCNRFJEBDFRNBDUJYIAALOAFGPDFGSQAEHRQYBWQHVYSWLBXHQE" {
		t.Error(decoded.GetQuTransferEvent().GetSourceId())
	}

	if decoded.GetQuTransferEvent().GetDestId() != "BAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAARMID" {
		t.Error(decoded.GetQuTransferEvent().GetDestId())
	}

	if decoded.GetQuTransferEvent().Amount != 1_000_000 {
		t.Error(decoded.GetQuTransferEvent().Amount)
	}
}

//goland:noinspection SpellCheckingInspection
func TestEventDecoder_Decode_AssetOwnershipChangeEvent(t *testing.T) {

	eventData, err := base64.StdEncoding.DecodeString("sMmo18V9WMO9LstUtxvWC2ZfJc2/FZWKEUdAKOqNKDIvyKKaekppac06VyRMSMUCe1tpQO0R9znQUrQOndNX+ggwu2O/fV4WSsjL04aAYw/3Zwoevzn3IQtAvNyiU9Bf2XE+AAAAAABDRkIAAAAAAADQANAjGBU=")
	if err != nil {
		t.Error(err)
	}
	decoded, err := DecodeAssetOwnershipChangeEvent(eventData)
	if err != nil {
		t.Error(err)
	}

	if decoded.GetAssetOwnershipChangeEvent().GetSourceId() != "AKJDFZYITPCNRFJEBDFRNBDUJYIAALOAFGPDFGSQAEHRQYBWQHVYSWLBXHQE" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetSourceId())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetDestId() != "VFWIEWBYSIMPBDHBXYFJVMLGKCCABZKRYFLQJVZTRBUOYSUHOODPVAHHKXPJ" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetDestId())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetAssetName() != "CFB" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetAssetName())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetIssuerId() != "CFBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJL" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetIssuerId())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetNumberOfShares() != 4092377 {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetNumberOfShares())
	}
}

//goland:noinspection SpellCheckingInspection
func TestEventDecoder_Decode_AssetPossessionChangeEvent(t *testing.T) {

	eventData, err := base64.StdEncoding.DecodeString("sMmo18V9WMO9LstUtxvWC2ZfJc2/FZWKEUdAKOqNKDIvyKKaekppac06VyRMSMUCe1tpQO0R9znQUrQOndNX+ggwu2O/fV4WSsjL04aAYw/3Zwoevzn3IQtAvNyiU9Bf2XE+AAAAAABDRkIAAAAAAADQANAjGBU=")
	if err != nil {
		t.Error(err)
	}
	decoded, err := DecodeAssetPossessionChangeEvent(eventData)
	if err != nil {
		t.Error(err)
	}

	if decoded.GetAssetPossessionChangeEvent().GetSourceId() != "AKJDFZYITPCNRFJEBDFRNBDUJYIAALOAFGPDFGSQAEHRQYBWQHVYSWLBXHQE" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetSourceId())
	}

	if decoded.GetAssetPossessionChangeEvent().GetDestId() != "VFWIEWBYSIMPBDHBXYFJVMLGKCCABZKRYFLQJVZTRBUOYSUHOODPVAHHKXPJ" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetDestId())
	}

	if decoded.GetAssetPossessionChangeEvent().GetAssetName() != "CFB" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetAssetName())
	}

	if decoded.GetAssetPossessionChangeEvent().GetIssuerId() != "CFBMEMZOIDEXQAUXYYSZIURADQLAPWPMNJXQSNVQZAHYVOPYUKKJBJUCTVJL" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetIssuerId())
	}

	if decoded.GetAssetPossessionChangeEvent().GetNumberOfShares() != 4092377 {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetNumberOfShares())
	}
}
