// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.1
// 	protoc        v3.21.12
// source: pkg/server/api/client/service.proto

package client

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_pkg_server_api_client_service_proto protoreflect.FileDescriptor

var file_pkg_server_api_client_service_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x70, 0x6b, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x1a, 0x20, 0x70, 0x6b,
	0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b,
	0x70, 0x6b, 0x67, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x70, 0x6f,
	0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2a, 0x70, 0x6b, 0x67,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0xef, 0x02, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x07, 0x41, 0x75, 0x74,
	0x68, 0x5f, 0x76, 0x31, 0x12, 0x25, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x75, 0x74,
	0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x76, 0x31, 0x1a, 0x26, 0x2e, 0x70, 0x6b,
	0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x5f, 0x76, 0x31, 0x22, 0x00, 0x12, 0x78, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x50, 0x6f, 0x6c, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x76, 0x31, 0x12, 0x2f, 0x2e, 0x70, 0x6b, 0x67,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x76, 0x31, 0x1a, 0x30, 0x2e, 0x70, 0x6b,
	0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x6f, 0x6c, 0x6c, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x76, 0x31, 0x22, 0x00, 0x12,
	0x87, 0x01, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x76, 0x31, 0x12, 0x34, 0x2e, 0x70, 0x6b, 0x67,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x76, 0x31,
	0x1a, 0x35, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x76, 0x31, 0x22, 0x00, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x61, 0x69, 0x6f, 0x6e, 0x2d, 0x6f,
	0x72, 0x67, 0x2f, 0x63, 0x6c, 0x61, 0x69, 0x66, 0x6c, 0x6f, 0x77, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_pkg_server_api_client_service_proto_goTypes = []interface{}{
	(*AuthRequestV1)(nil),                 // 0: pkg.server.api.client.AuthRequest_v1
	(*ServicePollingRequestV1)(nil),       // 1: pkg.server.api.client.ServicePollingRequest_v1
	(*UpdateServiceStatusRequestV1)(nil),  // 2: pkg.server.api.client.UpdateServiceStatusRequest_v1
	(*AuthResponseV1)(nil),                // 3: pkg.server.api.client.AuthResponse_v1
	(*ServicePollingResponseV1)(nil),      // 4: pkg.server.api.client.ServicePollingResponse_v1
	(*UpdateServiceStatusResponseV1)(nil), // 5: pkg.server.api.client.UpdateServiceStatusResponse_v1
}
var file_pkg_server_api_client_service_proto_depIdxs = []int32{
	0, // 0: pkg.server.api.client.ClientService.Auth_v1:input_type -> pkg.server.api.client.AuthRequest_v1
	1, // 1: pkg.server.api.client.ClientService.ServicePolling_v1:input_type -> pkg.server.api.client.ServicePollingRequest_v1
	2, // 2: pkg.server.api.client.ClientService.UpdateServiceStatus_v1:input_type -> pkg.server.api.client.UpdateServiceStatusRequest_v1
	3, // 3: pkg.server.api.client.ClientService.Auth_v1:output_type -> pkg.server.api.client.AuthResponse_v1
	4, // 4: pkg.server.api.client.ClientService.ServicePolling_v1:output_type -> pkg.server.api.client.ServicePollingResponse_v1
	5, // 5: pkg.server.api.client.ClientService.UpdateServiceStatus_v1:output_type -> pkg.server.api.client.UpdateServiceStatusResponse_v1
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_server_api_client_service_proto_init() }
func file_pkg_server_api_client_service_proto_init() {
	if File_pkg_server_api_client_service_proto != nil {
		return
	}
	file_pkg_server_api_client_auth_proto_init()
	file_pkg_server_api_client_service_polling_proto_init()
	file_pkg_server_api_client_update_service_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_server_api_client_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_server_api_client_service_proto_goTypes,
		DependencyIndexes: file_pkg_server_api_client_service_proto_depIdxs,
	}.Build()
	File_pkg_server_api_client_service_proto = out.File
	file_pkg_server_api_client_service_proto_rawDesc = nil
	file_pkg_server_api_client_service_proto_goTypes = nil
	file_pkg_server_api_client_service_proto_depIdxs = nil
}
