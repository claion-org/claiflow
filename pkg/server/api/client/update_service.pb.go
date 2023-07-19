// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: pkg/server/api/client/update_service.proto

package client

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpdateServiceStatusRequestV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Uuid     string                 `protobuf:"bytes,1,opt,name=uuid,proto3" json:"uuid,omitempty"`
	Sequence int32                  `protobuf:"varint,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Status   int32                  `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
	Result   string                 `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
	Started  *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=started,proto3" json:"started,omitempty"`
	Ended    *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=ended,proto3" json:"ended,omitempty"`
	Error    string                 `protobuf:"bytes,7,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateServiceStatusRequestV1) Reset() {
	*x = UpdateServiceStatusRequestV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_server_api_client_update_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateServiceStatusRequestV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateServiceStatusRequestV1) ProtoMessage() {}

func (x *UpdateServiceStatusRequestV1) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_server_api_client_update_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateServiceStatusRequestV1.ProtoReflect.Descriptor instead.
func (*UpdateServiceStatusRequestV1) Descriptor() ([]byte, []int) {
	return file_pkg_server_api_client_update_service_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateServiceStatusRequestV1) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *UpdateServiceStatusRequestV1) GetSequence() int32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *UpdateServiceStatusRequestV1) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *UpdateServiceStatusRequestV1) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *UpdateServiceStatusRequestV1) GetStarted() *timestamppb.Timestamp {
	if x != nil {
		return x.Started
	}
	return nil
}

func (x *UpdateServiceStatusRequestV1) GetEnded() *timestamppb.Timestamp {
	if x != nil {
		return x.Ended
	}
	return nil
}

func (x *UpdateServiceStatusRequestV1) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type UpdateServiceStatusResponseV1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateServiceStatusResponseV1) Reset() {
	*x = UpdateServiceStatusResponseV1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_server_api_client_update_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateServiceStatusResponseV1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateServiceStatusResponseV1) ProtoMessage() {}

func (x *UpdateServiceStatusResponseV1) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_server_api_client_update_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateServiceStatusResponseV1.ProtoReflect.Descriptor instead.
func (*UpdateServiceStatusResponseV1) Descriptor() ([]byte, []int) {
	return file_pkg_server_api_client_update_service_proto_rawDescGZIP(), []int{1}
}

var File_pkg_server_api_client_update_service_proto protoreflect.FileDescriptor

var file_pkg_server_api_client_update_service_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x70, 0x6b,
	0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfd, 0x01, 0x0a, 0x1d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x5f, 0x76, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x34, 0x0a, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x07, 0x73, 0x74, 0x61, 0x72, 0x74, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x05,
	0x65, 0x6e, 0x64, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x65, 0x6e, 0x64, 0x65, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x22, 0x20, 0x0a, 0x1e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x5f, 0x76, 0x31, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x61, 0x69, 0x6f, 0x6e, 0x2d, 0x6f, 0x72, 0x67, 0x2f,
	0x63, 0x6c, 0x61, 0x69, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_server_api_client_update_service_proto_rawDescOnce sync.Once
	file_pkg_server_api_client_update_service_proto_rawDescData = file_pkg_server_api_client_update_service_proto_rawDesc
)

func file_pkg_server_api_client_update_service_proto_rawDescGZIP() []byte {
	file_pkg_server_api_client_update_service_proto_rawDescOnce.Do(func() {
		file_pkg_server_api_client_update_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_server_api_client_update_service_proto_rawDescData)
	})
	return file_pkg_server_api_client_update_service_proto_rawDescData
}

var file_pkg_server_api_client_update_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_pkg_server_api_client_update_service_proto_goTypes = []interface{}{
	(*UpdateServiceStatusRequestV1)(nil),  // 0: pkg.server.api.client.UpdateServiceStatusRequest_v1
	(*UpdateServiceStatusResponseV1)(nil), // 1: pkg.server.api.client.UpdateServiceStatusResponse_v1
	(*timestamppb.Timestamp)(nil),         // 2: google.protobuf.Timestamp
}
var file_pkg_server_api_client_update_service_proto_depIdxs = []int32{
	2, // 0: pkg.server.api.client.UpdateServiceStatusRequest_v1.started:type_name -> google.protobuf.Timestamp
	2, // 1: pkg.server.api.client.UpdateServiceStatusRequest_v1.ended:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_server_api_client_update_service_proto_init() }
func file_pkg_server_api_client_update_service_proto_init() {
	if File_pkg_server_api_client_update_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_server_api_client_update_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateServiceStatusRequestV1); i {
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
		file_pkg_server_api_client_update_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateServiceStatusResponseV1); i {
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
			RawDescriptor: file_pkg_server_api_client_update_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_server_api_client_update_service_proto_goTypes,
		DependencyIndexes: file_pkg_server_api_client_update_service_proto_depIdxs,
		MessageInfos:      file_pkg_server_api_client_update_service_proto_msgTypes,
	}.Build()
	File_pkg_server_api_client_update_service_proto = out.File
	file_pkg_server_api_client_update_service_proto_rawDesc = nil
	file_pkg_server_api_client_update_service_proto_goTypes = nil
	file_pkg_server_api_client_update_service_proto_depIdxs = nil
}
