// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.4
// source: querier.proto

package relayproto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

var File_querier_proto protoreflect.FileDescriptor

var file_querier_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x71, 0x75, 0x65, 0x72, 0x69, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x71, 0x75, 0x65,
	0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x9f, 0x07, 0x0a, 0x05, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x12, 0x51, 0x0a, 0x0a, 0x49, 0x73, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72,
	0x12, 0x21, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x49, 0x73, 0x41, 0x6e, 0x63, 0x65, 0x73,
	0x74, 0x6f, 0x72, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x49, 0x73, 0x41, 0x6e, 0x63, 0x65, 0x73,
	0x74, 0x6f, 0x72, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6c, 0x61,
	0x79, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x12, 0x17, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x1a, 0x23, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x47, 0x65, 0x74, 0x52, 0x65, 0x6c, 0x61, 0x79, 0x47,
	0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4c,
	0x61, 0x73, 0x74, 0x52, 0x65, 0x6f, 0x72, 0x67, 0x4c, 0x43, 0x41, 0x12, 0x17, 0x2e, 0x72, 0x65,
	0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x1a, 0x23, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x47, 0x65, 0x74, 0x4c, 0x61, 0x73,
	0x74, 0x52, 0x65, 0x6f, 0x72, 0x67, 0x4c, 0x43, 0x41, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x0d, 0x47,
	0x65, 0x74, 0x42, 0x65, 0x73, 0x74, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x12, 0x17, 0x2e, 0x72,
	0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x47, 0x65, 0x74, 0x42, 0x65,
	0x73, 0x74, 0x44, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0x00, 0x12, 0x57, 0x0a, 0x0c, 0x46, 0x69,
	0x6e, 0x64, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x12, 0x23, 0x2e, 0x72, 0x65, 0x6c,
	0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x1a,
	0x20, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x52, 0x65, 0x73, 0x46, 0x69, 0x6e, 0x64, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f,
	0x72, 0x22, 0x00, 0x12, 0x6f, 0x0a, 0x14, 0x48, 0x65, 0x61, 0x76, 0x69, 0x65, 0x73, 0x74, 0x46,
	0x72, 0x6f, 0x6d, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x12, 0x2b, 0x2e, 0x72, 0x65,
	0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x48, 0x65, 0x61, 0x76, 0x69, 0x65, 0x73, 0x74, 0x46, 0x72, 0x6f, 0x6d,
	0x41, 0x6e, 0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x1a, 0x28, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x48, 0x65,
	0x61, 0x76, 0x69, 0x65, 0x73, 0x74, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x6e, 0x63, 0x65, 0x73, 0x74,
	0x6f, 0x72, 0x22, 0x00, 0x12, 0x81, 0x01, 0x0a, 0x1a, 0x49, 0x73, 0x4d, 0x6f, 0x73, 0x74, 0x52,
	0x65, 0x63, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x6e, 0x63, 0x65, 0x73,
	0x74, 0x6f, 0x72, 0x12, 0x31, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x49, 0x73, 0x4d, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x6e,
	0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x1a, 0x2e, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x49, 0x73, 0x4d, 0x6f,
	0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x41, 0x6e,
	0x63, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x6c, 0x61,
	0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x00, 0x12, 0x5a, 0x0a, 0x0d, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x12, 0x24, 0x2e, 0x72,
	0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x73, 0x1a, 0x21, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x73, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0a, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x50, 0x72, 0x6f, 0x6f, 0x66, 0x12, 0x21, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x1a, 0x1e, 0x2e, 0x72, 0x65, 0x6c, 0x61, 0x79,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x50, 0x72, 0x6f, 0x6f, 0x66, 0x22, 0x00, 0x42, 0x27, 0x5a, 0x25, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x2d, 0x74,
	0x78, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x73, 0x2f, 0x72, 0x65, 0x6c, 0x61, 0x79, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_querier_proto_goTypes = []interface{}{
	(*QueryParamsIsAncestor)(nil),                 // 0: relayproto.QueryParamsIsAncestor
	(*EmptyParams)(nil),                           // 1: relayproto.EmptyParams
	(*QueryParamsFindAncestor)(nil),               // 2: relayproto.QueryParamsFindAncestor
	(*QueryParamsHeaviestFromAncestor)(nil),       // 3: relayproto.QueryParamsHeaviestFromAncestor
	(*QueryParamsIsMostRecentCommonAncestor)(nil), // 4: relayproto.QueryParamsIsMostRecentCommonAncestor
	(*QueryParamsGetRequest)(nil),                 // 5: relayproto.QueryParamsGetRequest
	(*QueryParamsCheckRequests)(nil),              // 6: relayproto.QueryParamsCheckRequests
	(*QueryParamsCheckProof)(nil),                 // 7: relayproto.QueryParamsCheckProof
	(*QueryResIsAncestor)(nil),                    // 8: relayproto.QueryResIsAncestor
	(*QueryResGetRelayGenesis)(nil),               // 9: relayproto.QueryResGetRelayGenesis
	(*QueryResGetLastReorgLCA)(nil),               // 10: relayproto.QueryResGetLastReorgLCA
	(*QueryResGetBestDigest)(nil),                 // 11: relayproto.QueryResGetBestDigest
	(*QueryResFindAncestor)(nil),                  // 12: relayproto.QueryResFindAncestor
	(*QueryResHeaviestFromAncestor)(nil),          // 13: relayproto.QueryResHeaviestFromAncestor
	(*QueryResIsMostRecentCommonAncestor)(nil),    // 14: relayproto.QueryResIsMostRecentCommonAncestor
	(*QueryResGetRequest)(nil),                    // 15: relayproto.QueryResGetRequest
	(*QueryResCheckRequests)(nil),                 // 16: relayproto.QueryResCheckRequests
	(*QueryResCheckProof)(nil),                    // 17: relayproto.QueryResCheckProof
}
var file_querier_proto_depIdxs = []int32{
	0,  // 0: relayproto.Query.IsAncestor:input_type -> relayproto.QueryParamsIsAncestor
	1,  // 1: relayproto.Query.GetRelayGenesis:input_type -> relayproto.EmptyParams
	1,  // 2: relayproto.Query.GetLastReorgLCA:input_type -> relayproto.EmptyParams
	1,  // 3: relayproto.Query.GetBestDigest:input_type -> relayproto.EmptyParams
	2,  // 4: relayproto.Query.FindAncestor:input_type -> relayproto.QueryParamsFindAncestor
	3,  // 5: relayproto.Query.HeaviestFromAncestor:input_type -> relayproto.QueryParamsHeaviestFromAncestor
	4,  // 6: relayproto.Query.IsMostRecentCommonAncestor:input_type -> relayproto.QueryParamsIsMostRecentCommonAncestor
	5,  // 7: relayproto.Query.GetRequest:input_type -> relayproto.QueryParamsGetRequest
	6,  // 8: relayproto.Query.CheckRequests:input_type -> relayproto.QueryParamsCheckRequests
	7,  // 9: relayproto.Query.CheckProof:input_type -> relayproto.QueryParamsCheckProof
	8,  // 10: relayproto.Query.IsAncestor:output_type -> relayproto.QueryResIsAncestor
	9,  // 11: relayproto.Query.GetRelayGenesis:output_type -> relayproto.QueryResGetRelayGenesis
	10, // 12: relayproto.Query.GetLastReorgLCA:output_type -> relayproto.QueryResGetLastReorgLCA
	11, // 13: relayproto.Query.GetBestDigest:output_type -> relayproto.QueryResGetBestDigest
	12, // 14: relayproto.Query.FindAncestor:output_type -> relayproto.QueryResFindAncestor
	13, // 15: relayproto.Query.HeaviestFromAncestor:output_type -> relayproto.QueryResHeaviestFromAncestor
	14, // 16: relayproto.Query.IsMostRecentCommonAncestor:output_type -> relayproto.QueryResIsMostRecentCommonAncestor
	15, // 17: relayproto.Query.GetRequest:output_type -> relayproto.QueryResGetRequest
	16, // 18: relayproto.Query.CheckRequests:output_type -> relayproto.QueryResCheckRequests
	17, // 19: relayproto.Query.CheckProof:output_type -> relayproto.QueryResCheckProof
	10, // [10:20] is the sub-list for method output_type
	0,  // [0:10] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_querier_proto_init() }
func file_querier_proto_init() {
	if File_querier_proto != nil {
		return
	}
	file_queries_proto_init()
	file_shared_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_querier_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_querier_proto_goTypes,
		DependencyIndexes: file_querier_proto_depIdxs,
	}.Build()
	File_querier_proto = out.File
	file_querier_proto_rawDesc = nil
	file_querier_proto_goTypes = nil
	file_querier_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	IsAncestor(ctx context.Context, in *QueryParamsIsAncestor, opts ...grpc.CallOption) (*QueryResIsAncestor, error)
	GetRelayGenesis(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*QueryResGetRelayGenesis, error)
	GetLastReorgLCA(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*QueryResGetLastReorgLCA, error)
	GetBestDigest(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*QueryResGetBestDigest, error)
	FindAncestor(ctx context.Context, in *QueryParamsFindAncestor, opts ...grpc.CallOption) (*QueryResFindAncestor, error)
	HeaviestFromAncestor(ctx context.Context, in *QueryParamsHeaviestFromAncestor, opts ...grpc.CallOption) (*QueryResHeaviestFromAncestor, error)
	IsMostRecentCommonAncestor(ctx context.Context, in *QueryParamsIsMostRecentCommonAncestor, opts ...grpc.CallOption) (*QueryResIsMostRecentCommonAncestor, error)
	GetRequest(ctx context.Context, in *QueryParamsGetRequest, opts ...grpc.CallOption) (*QueryResGetRequest, error)
	CheckRequests(ctx context.Context, in *QueryParamsCheckRequests, opts ...grpc.CallOption) (*QueryResCheckRequests, error)
	CheckProof(ctx context.Context, in *QueryParamsCheckProof, opts ...grpc.CallOption) (*QueryResCheckProof, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) IsAncestor(ctx context.Context, in *QueryParamsIsAncestor, opts ...grpc.CallOption) (*QueryResIsAncestor, error) {
	out := new(QueryResIsAncestor)
	err := c.cc.Invoke(ctx, "/relayproto.Query/IsAncestor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetRelayGenesis(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*QueryResGetRelayGenesis, error) {
	out := new(QueryResGetRelayGenesis)
	err := c.cc.Invoke(ctx, "/relayproto.Query/GetRelayGenesis", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetLastReorgLCA(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*QueryResGetLastReorgLCA, error) {
	out := new(QueryResGetLastReorgLCA)
	err := c.cc.Invoke(ctx, "/relayproto.Query/GetLastReorgLCA", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetBestDigest(ctx context.Context, in *EmptyParams, opts ...grpc.CallOption) (*QueryResGetBestDigest, error) {
	out := new(QueryResGetBestDigest)
	err := c.cc.Invoke(ctx, "/relayproto.Query/GetBestDigest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) FindAncestor(ctx context.Context, in *QueryParamsFindAncestor, opts ...grpc.CallOption) (*QueryResFindAncestor, error) {
	out := new(QueryResFindAncestor)
	err := c.cc.Invoke(ctx, "/relayproto.Query/FindAncestor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) HeaviestFromAncestor(ctx context.Context, in *QueryParamsHeaviestFromAncestor, opts ...grpc.CallOption) (*QueryResHeaviestFromAncestor, error) {
	out := new(QueryResHeaviestFromAncestor)
	err := c.cc.Invoke(ctx, "/relayproto.Query/HeaviestFromAncestor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) IsMostRecentCommonAncestor(ctx context.Context, in *QueryParamsIsMostRecentCommonAncestor, opts ...grpc.CallOption) (*QueryResIsMostRecentCommonAncestor, error) {
	out := new(QueryResIsMostRecentCommonAncestor)
	err := c.cc.Invoke(ctx, "/relayproto.Query/IsMostRecentCommonAncestor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) GetRequest(ctx context.Context, in *QueryParamsGetRequest, opts ...grpc.CallOption) (*QueryResGetRequest, error) {
	out := new(QueryResGetRequest)
	err := c.cc.Invoke(ctx, "/relayproto.Query/GetRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CheckRequests(ctx context.Context, in *QueryParamsCheckRequests, opts ...grpc.CallOption) (*QueryResCheckRequests, error) {
	out := new(QueryResCheckRequests)
	err := c.cc.Invoke(ctx, "/relayproto.Query/CheckRequests", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) CheckProof(ctx context.Context, in *QueryParamsCheckProof, opts ...grpc.CallOption) (*QueryResCheckProof, error) {
	out := new(QueryResCheckProof)
	err := c.cc.Invoke(ctx, "/relayproto.Query/CheckProof", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	IsAncestor(context.Context, *QueryParamsIsAncestor) (*QueryResIsAncestor, error)
	GetRelayGenesis(context.Context, *EmptyParams) (*QueryResGetRelayGenesis, error)
	GetLastReorgLCA(context.Context, *EmptyParams) (*QueryResGetLastReorgLCA, error)
	GetBestDigest(context.Context, *EmptyParams) (*QueryResGetBestDigest, error)
	FindAncestor(context.Context, *QueryParamsFindAncestor) (*QueryResFindAncestor, error)
	HeaviestFromAncestor(context.Context, *QueryParamsHeaviestFromAncestor) (*QueryResHeaviestFromAncestor, error)
	IsMostRecentCommonAncestor(context.Context, *QueryParamsIsMostRecentCommonAncestor) (*QueryResIsMostRecentCommonAncestor, error)
	GetRequest(context.Context, *QueryParamsGetRequest) (*QueryResGetRequest, error)
	CheckRequests(context.Context, *QueryParamsCheckRequests) (*QueryResCheckRequests, error)
	CheckProof(context.Context, *QueryParamsCheckProof) (*QueryResCheckProof, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) IsAncestor(context.Context, *QueryParamsIsAncestor) (*QueryResIsAncestor, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsAncestor not implemented")
}
func (*UnimplementedQueryServer) GetRelayGenesis(context.Context, *EmptyParams) (*QueryResGetRelayGenesis, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRelayGenesis not implemented")
}
func (*UnimplementedQueryServer) GetLastReorgLCA(context.Context, *EmptyParams) (*QueryResGetLastReorgLCA, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLastReorgLCA not implemented")
}
func (*UnimplementedQueryServer) GetBestDigest(context.Context, *EmptyParams) (*QueryResGetBestDigest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBestDigest not implemented")
}
func (*UnimplementedQueryServer) FindAncestor(context.Context, *QueryParamsFindAncestor) (*QueryResFindAncestor, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAncestor not implemented")
}
func (*UnimplementedQueryServer) HeaviestFromAncestor(context.Context, *QueryParamsHeaviestFromAncestor) (*QueryResHeaviestFromAncestor, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeaviestFromAncestor not implemented")
}
func (*UnimplementedQueryServer) IsMostRecentCommonAncestor(context.Context, *QueryParamsIsMostRecentCommonAncestor) (*QueryResIsMostRecentCommonAncestor, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsMostRecentCommonAncestor not implemented")
}
func (*UnimplementedQueryServer) GetRequest(context.Context, *QueryParamsGetRequest) (*QueryResGetRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRequest not implemented")
}
func (*UnimplementedQueryServer) CheckRequests(context.Context, *QueryParamsCheckRequests) (*QueryResCheckRequests, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckRequests not implemented")
}
func (*UnimplementedQueryServer) CheckProof(context.Context, *QueryParamsCheckProof) (*QueryResCheckProof, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckProof not implemented")
}

func RegisterQueryServer(s *grpc.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_IsAncestor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsIsAncestor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IsAncestor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/IsAncestor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IsAncestor(ctx, req.(*QueryParamsIsAncestor))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetRelayGenesis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetRelayGenesis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/GetRelayGenesis",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetRelayGenesis(ctx, req.(*EmptyParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetLastReorgLCA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetLastReorgLCA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/GetLastReorgLCA",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetLastReorgLCA(ctx, req.(*EmptyParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetBestDigest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetBestDigest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/GetBestDigest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetBestDigest(ctx, req.(*EmptyParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_FindAncestor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsFindAncestor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).FindAncestor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/FindAncestor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).FindAncestor(ctx, req.(*QueryParamsFindAncestor))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_HeaviestFromAncestor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsHeaviestFromAncestor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).HeaviestFromAncestor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/HeaviestFromAncestor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).HeaviestFromAncestor(ctx, req.(*QueryParamsHeaviestFromAncestor))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IsMostRecentCommonAncestor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsIsMostRecentCommonAncestor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IsMostRecentCommonAncestor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/IsMostRecentCommonAncestor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IsMostRecentCommonAncestor(ctx, req.(*QueryParamsIsMostRecentCommonAncestor))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_GetRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).GetRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/GetRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).GetRequest(ctx, req.(*QueryParamsGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CheckRequests_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsCheckRequests)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CheckRequests(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/CheckRequests",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CheckRequests(ctx, req.(*QueryParamsCheckRequests))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_CheckProof_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsCheckProof)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).CheckProof(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/relayproto.Query/CheckProof",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).CheckProof(ctx, req.(*QueryParamsCheckProof))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "relayproto.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsAncestor",
			Handler:    _Query_IsAncestor_Handler,
		},
		{
			MethodName: "GetRelayGenesis",
			Handler:    _Query_GetRelayGenesis_Handler,
		},
		{
			MethodName: "GetLastReorgLCA",
			Handler:    _Query_GetLastReorgLCA_Handler,
		},
		{
			MethodName: "GetBestDigest",
			Handler:    _Query_GetBestDigest_Handler,
		},
		{
			MethodName: "FindAncestor",
			Handler:    _Query_FindAncestor_Handler,
		},
		{
			MethodName: "HeaviestFromAncestor",
			Handler:    _Query_HeaviestFromAncestor_Handler,
		},
		{
			MethodName: "IsMostRecentCommonAncestor",
			Handler:    _Query_IsMostRecentCommonAncestor_Handler,
		},
		{
			MethodName: "GetRequest",
			Handler:    _Query_GetRequest_Handler,
		},
		{
			MethodName: "CheckRequests",
			Handler:    _Query_CheckRequests_Handler,
		},
		{
			MethodName: "CheckProof",
			Handler:    _Query_CheckProof_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "querier.proto",
}
