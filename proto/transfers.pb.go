// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.21.12
// source: transfers.proto

package proto

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HealthResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transfers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transfers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HealthResponse.ProtoReflect.Descriptor instead.
func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_transfers_proto_rawDescGZIP(), []int{0}
}

func (x *HealthResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type TickRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tick uint32 `protobuf:"varint,1,opt,name=tick,proto3" json:"tick,omitempty"`
}

func (x *TickRequest) Reset() {
	*x = TickRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transfers_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TickRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TickRequest) ProtoMessage() {}

func (x *TickRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transfers_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TickRequest.ProtoReflect.Descriptor instead.
func (*TickRequest) Descriptor() ([]byte, []int) {
	return file_transfers_proto_rawDescGZIP(), []int{1}
}

func (x *TickRequest) GetTick() uint32 {
	if x != nil {
		return x.Tick
	}
	return 0
}

type EntityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Identity string `protobuf:"bytes,1,opt,name=identity,proto3" json:"identity,omitempty"`
}

func (x *EntityRequest) Reset() {
	*x = EntityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transfers_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EntityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EntityRequest) ProtoMessage() {}

func (x *EntityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_transfers_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EntityRequest.ProtoReflect.Descriptor instead.
func (*EntityRequest) Descriptor() ([]byte, []int) {
	return file_transfers_proto_rawDescGZIP(), []int{2}
}

func (x *EntityRequest) GetIdentity() string {
	if x != nil {
		return x.Identity
	}
	return ""
}

type AssetTransferResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceId       string `protobuf:"bytes,1,opt,name=sourceId,proto3" json:"sourceId,omitempty"`
	DestinationId  string `protobuf:"bytes,2,opt,name=destinationId,proto3" json:"destinationId,omitempty"`
	IssuerId       string `protobuf:"bytes,3,opt,name=issuerId,proto3" json:"issuerId,omitempty"`
	Name           string `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	NumberOfShares uint64 `protobuf:"varint,5,opt,name=numberOfShares,proto3" json:"numberOfShares,omitempty"`
}

func (x *AssetTransferResponse) Reset() {
	*x = AssetTransferResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_transfers_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssetTransferResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssetTransferResponse) ProtoMessage() {}

func (x *AssetTransferResponse) ProtoReflect() protoreflect.Message {
	mi := &file_transfers_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssetTransferResponse.ProtoReflect.Descriptor instead.
func (*AssetTransferResponse) Descriptor() ([]byte, []int) {
	return file_transfers_proto_rawDescGZIP(), []int{3}
}

func (x *AssetTransferResponse) GetSourceId() string {
	if x != nil {
		return x.SourceId
	}
	return ""
}

func (x *AssetTransferResponse) GetDestinationId() string {
	if x != nil {
		return x.DestinationId
	}
	return ""
}

func (x *AssetTransferResponse) GetIssuerId() string {
	if x != nil {
		return x.IssuerId
	}
	return ""
}

func (x *AssetTransferResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AssetTransferResponse) GetNumberOfShares() uint64 {
	if x != nil {
		return x.NumberOfShares
	}
	return 0
}

var File_transfers_proto protoreflect.FileDescriptor

var file_transfers_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x15, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x28, 0x0a, 0x0e, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x21, 0x0a,
	0x0b, 0x54, 0x69, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x69, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x69, 0x63, 0x6b,
	0x22, 0x2b, 0x0a, 0x0d, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0xb1, 0x01,
	0x0a, 0x15, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x0d, 0x64, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x64, 0x65, 0x73, 0x74,
	0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x72, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x73, 0x73,
	0x75, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x6e, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x4f, 0x66, 0x53, 0x68, 0x61, 0x72, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x0e, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x53, 0x68, 0x61, 0x72, 0x65,
	0x73, 0x32, 0xac, 0x03, 0x0a, 0x0f, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5f, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x25, 0x2e, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x16,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2f,
	0x68, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x12, 0x95, 0x01, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x41, 0x73,
	0x73, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x46, 0x6f, 0x72, 0x54,
	0x69, 0x63, 0x6b, 0x12, 0x22, 0x2e, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x27, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x21, 0x12, 0x1f, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x69, 0x63, 0x6b, 0x2f, 0x7b, 0x74, 0x69, 0x63, 0x6b, 0x7d, 0x2f, 0x61,
	0x73, 0x73, 0x65, 0x74, 0x2d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x12, 0x9f,
	0x01, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x41, 0x73, 0x73, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x65, 0x72, 0x73, 0x46, 0x6f, 0x72, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x24, 0x2e,
	0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x71, 0x75, 0x62, 0x69, 0x63, 0x2e, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x73, 0x73, 0x65,
	0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x2d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x12, 0x25, 0x2f, 0x76, 0x31, 0x2f, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x2f, 0x7b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x7d,
	0x2f, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73,
	0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x71,
	0x75, 0x62, 0x69, 0x63, 0x2f, 0x67, 0x6f, 0x2d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72,
	0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_transfers_proto_rawDescOnce sync.Once
	file_transfers_proto_rawDescData = file_transfers_proto_rawDesc
)

func file_transfers_proto_rawDescGZIP() []byte {
	file_transfers_proto_rawDescOnce.Do(func() {
		file_transfers_proto_rawDescData = protoimpl.X.CompressGZIP(file_transfers_proto_rawDescData)
	})
	return file_transfers_proto_rawDescData
}

var file_transfers_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_transfers_proto_goTypes = []any{
	(*HealthResponse)(nil),        // 0: qubic.transfers.proto.HealthResponse
	(*TickRequest)(nil),           // 1: qubic.transfers.proto.TickRequest
	(*EntityRequest)(nil),         // 2: qubic.transfers.proto.EntityRequest
	(*AssetTransferResponse)(nil), // 3: qubic.transfers.proto.AssetTransferResponse
	(*emptypb.Empty)(nil),         // 4: google.protobuf.Empty
}
var file_transfers_proto_depIdxs = []int32{
	4, // 0: qubic.transfers.proto.TransferService.Health:input_type -> google.protobuf.Empty
	1, // 1: qubic.transfers.proto.TransferService.GetAssetTransfersForTick:input_type -> qubic.transfers.proto.TickRequest
	2, // 2: qubic.transfers.proto.TransferService.GetAssetTransfersForEntity:input_type -> qubic.transfers.proto.EntityRequest
	0, // 3: qubic.transfers.proto.TransferService.Health:output_type -> qubic.transfers.proto.HealthResponse
	3, // 4: qubic.transfers.proto.TransferService.GetAssetTransfersForTick:output_type -> qubic.transfers.proto.AssetTransferResponse
	3, // 5: qubic.transfers.proto.TransferService.GetAssetTransfersForEntity:output_type -> qubic.transfers.proto.AssetTransferResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_transfers_proto_init() }
func file_transfers_proto_init() {
	if File_transfers_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_transfers_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*HealthResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_transfers_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TickRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_transfers_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EntityRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_transfers_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*AssetTransferResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_transfers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_transfers_proto_goTypes,
		DependencyIndexes: file_transfers_proto_depIdxs,
		MessageInfos:      file_transfers_proto_msgTypes,
	}.Build()
	File_transfers_proto = out.File
	file_transfers_proto_rawDesc = nil
	file_transfers_proto_goTypes = nil
	file_transfers_proto_depIdxs = nil
}
