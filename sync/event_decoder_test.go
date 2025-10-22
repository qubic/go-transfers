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

	eventData, err := base64.StdEncoding.DecodeString("QvMt7n7vPwdDhVUbxbRVOxMpx/7trku3V9udvL77Hfm0XNyWnewpiwi3DPqGYe9p1T1ee0dgKChsGN91xWt9RAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAE1MTQAAAAAAAAAAAAAAAA==")
	if err != nil {
		t.Error(err)
	}
	decoded, err := DecodeAssetOwnershipChangeEvent(eventData)
	if err != nil {
		t.Error(err)
	}

	if decoded.GetAssetOwnershipChangeEvent().GetSourceId() != "OEJPBLJQYOIMFADQMWCTAZYYWUSBNNJSMPQIXTNKIFBXAGTTYUYUTCGHUDDL" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetSourceId())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetDestId() != "SDXHYCDNHADCBEUUVBZWKYESTZBDLHDRAKDCSEKIEBCPFIZLCHGDQSZBZMSL" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetDestId())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetAssetName() != "MLM" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetAssetName())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetIssuerId() != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB" {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetIssuerId())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetManagingContractIndex() != 10 {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetManagingContractIndex())
	}

	if decoded.GetAssetOwnershipChangeEvent().GetNumberOfShares() != 1 {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetNumberOfShares())
	}
}

//goland:noinspection SpellCheckingInspection
func TestEventDecoder_Decode_AssetPossessionChangeEvent(t *testing.T) {

	eventData, err := base64.StdEncoding.DecodeString("QvMt7n7vPwdDhVUbxbRVOxMpx/7trku3V9udvL77Hfm0XNyWnewpiwi3DPqGYe9p1T1ee0dgKChsGN91xWt9RAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAKAAAAAAAAAE1MTQAAAAAAAAAAAAAAAA==")
	if err != nil {
		t.Error(err)
	}
	decoded, err := DecodeAssetPossessionChangeEvent(eventData)
	if err != nil {
		t.Error(err)
	}

	if decoded.GetAssetPossessionChangeEvent().GetSourceId() != "OEJPBLJQYOIMFADQMWCTAZYYWUSBNNJSMPQIXTNKIFBXAGTTYUYUTCGHUDDL" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetSourceId())
	}

	if decoded.GetAssetPossessionChangeEvent().GetDestId() != "SDXHYCDNHADCBEUUVBZWKYESTZBDLHDRAKDCSEKIEBCPFIZLCHGDQSZBZMSL" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetDestId())
	}

	if decoded.GetAssetPossessionChangeEvent().GetAssetName() != "MLM" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetAssetName())
	}

	if decoded.GetAssetPossessionChangeEvent().GetIssuerId() != "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFXIB" {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetIssuerId())
	}

	if decoded.GetAssetPossessionChangeEvent().GetManagingContractIndex() != 10 {
		t.Error(decoded.GetAssetOwnershipChangeEvent().GetManagingContractIndex())
	}

	if decoded.GetAssetPossessionChangeEvent().GetNumberOfShares() != 1 {
		t.Error(decoded.GetAssetPossessionChangeEvent().GetNumberOfShares())
	}
}
