// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.4
// source: plumber.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
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

type PlumberRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	Addr string `protobuf:"bytes,2,opt,name=addr,proto3" json:"addr,omitempty"`
}

func (x *PlumberRequest) Reset() {
	*x = PlumberRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plumber_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlumberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlumberRequest) ProtoMessage() {}

func (x *PlumberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_plumber_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlumberRequest.ProtoReflect.Descriptor instead.
func (*PlumberRequest) Descriptor() ([]byte, []int) {
	return file_plumber_proto_rawDescGZIP(), []int{0}
}

func (x *PlumberRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PlumberRequest) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

type PlumberResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *PlumberResponse) Reset() {
	*x = PlumberResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_plumber_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlumberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlumberResponse) ProtoMessage() {}

func (x *PlumberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_plumber_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlumberResponse.ProtoReflect.Descriptor instead.
func (*PlumberResponse) Descriptor() ([]byte, []int) {
	return file_plumber_proto_rawDescGZIP(), []int{1}
}

func (x *PlumberResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_plumber_proto protoreflect.FileDescriptor

var file_plumber_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x38, 0x0a, 0x0e, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x64, 0x64, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x61, 0x64, 0x64, 0x72,
	0x22, 0x25, 0x0a, 0x0f, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xce, 0x01, 0x0a, 0x07, 0x50, 0x6c, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0c, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4c, 0x69,
	0x6e, 0x6b, 0x73, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x75, 0x6d,
	0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x40, 0x0a, 0x0f, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x45, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x42, 0x0a, 0x11, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x44,
	0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x6c, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_plumber_proto_rawDescOnce sync.Once
	file_plumber_proto_rawDescData = file_plumber_proto_rawDesc
)

func file_plumber_proto_rawDescGZIP() []byte {
	file_plumber_proto_rawDescOnce.Do(func() {
		file_plumber_proto_rawDescData = protoimpl.X.CompressGZIP(file_plumber_proto_rawDescData)
	})
	return file_plumber_proto_rawDescData
}

var file_plumber_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_plumber_proto_goTypes = []interface{}{
	(*PlumberRequest)(nil),  // 0: proto.PlumberRequest
	(*PlumberResponse)(nil), // 1: proto.PlumberResponse
}
var file_plumber_proto_depIdxs = []int32{
	0, // 0: proto.Plumber.PlumberLinks:input_type -> proto.PlumberRequest
	0, // 1: proto.Plumber.PlumberExchange:input_type -> proto.PlumberRequest
	0, // 2: proto.Plumber.PlumberDisconnect:input_type -> proto.PlumberRequest
	1, // 3: proto.Plumber.PlumberLinks:output_type -> proto.PlumberResponse
	1, // 4: proto.Plumber.PlumberExchange:output_type -> proto.PlumberResponse
	1, // 5: proto.Plumber.PlumberDisconnect:output_type -> proto.PlumberResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_plumber_proto_init() }
func file_plumber_proto_init() {
	if File_plumber_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_plumber_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlumberRequest); i {
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
		file_plumber_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlumberResponse); i {
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
			RawDescriptor: file_plumber_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_plumber_proto_goTypes,
		DependencyIndexes: file_plumber_proto_depIdxs,
		MessageInfos:      file_plumber_proto_msgTypes,
	}.Build()
	File_plumber_proto = out.File
	file_plumber_proto_rawDesc = nil
	file_plumber_proto_goTypes = nil
	file_plumber_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PlumberClient is the client API for Plumber service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PlumberClient interface {
	PlumberLinks(ctx context.Context, in *PlumberRequest, opts ...grpc.CallOption) (*PlumberResponse, error)
	PlumberExchange(ctx context.Context, in *PlumberRequest, opts ...grpc.CallOption) (*PlumberResponse, error)
	PlumberDisconnect(ctx context.Context, in *PlumberRequest, opts ...grpc.CallOption) (*PlumberResponse, error)
}

type plumberClient struct {
	cc grpc.ClientConnInterface
}

func NewPlumberClient(cc grpc.ClientConnInterface) PlumberClient {
	return &plumberClient{cc}
}

func (c *plumberClient) PlumberLinks(ctx context.Context, in *PlumberRequest, opts ...grpc.CallOption) (*PlumberResponse, error) {
	out := new(PlumberResponse)
	err := c.cc.Invoke(ctx, "/proto.Plumber/PlumberLinks", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plumberClient) PlumberExchange(ctx context.Context, in *PlumberRequest, opts ...grpc.CallOption) (*PlumberResponse, error) {
	out := new(PlumberResponse)
	err := c.cc.Invoke(ctx, "/proto.Plumber/PlumberExchange", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plumberClient) PlumberDisconnect(ctx context.Context, in *PlumberRequest, opts ...grpc.CallOption) (*PlumberResponse, error) {
	out := new(PlumberResponse)
	err := c.cc.Invoke(ctx, "/proto.Plumber/PlumberDisconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlumberServer is the server API for Plumber service.
type PlumberServer interface {
	PlumberLinks(context.Context, *PlumberRequest) (*PlumberResponse, error)
	PlumberExchange(context.Context, *PlumberRequest) (*PlumberResponse, error)
	PlumberDisconnect(context.Context, *PlumberRequest) (*PlumberResponse, error)
}

// UnimplementedPlumberServer can be embedded to have forward compatible implementations.
type UnimplementedPlumberServer struct {
}

func (*UnimplementedPlumberServer) PlumberLinks(context.Context, *PlumberRequest) (*PlumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlumberLinks not implemented")
}
func (*UnimplementedPlumberServer) PlumberExchange(context.Context, *PlumberRequest) (*PlumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlumberExchange not implemented")
}
func (*UnimplementedPlumberServer) PlumberDisconnect(context.Context, *PlumberRequest) (*PlumberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlumberDisconnect not implemented")
}

func RegisterPlumberServer(s *grpc.Server, srv PlumberServer) {
	s.RegisterService(&_Plumber_serviceDesc, srv)
}

func _Plumber_PlumberLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlumberServer).PlumberLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plumber/PlumberLinks",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlumberServer).PlumberLinks(ctx, req.(*PlumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plumber_PlumberExchange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlumberServer).PlumberExchange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plumber/PlumberExchange",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlumberServer).PlumberExchange(ctx, req.(*PlumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plumber_PlumberDisconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlumberServer).PlumberDisconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Plumber/PlumberDisconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlumberServer).PlumberDisconnect(ctx, req.(*PlumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Plumber_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Plumber",
	HandlerType: (*PlumberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PlumberLinks",
			Handler:    _Plumber_PlumberLinks_Handler,
		},
		{
			MethodName: "PlumberExchange",
			Handler:    _Plumber_PlumberExchange_Handler,
		},
		{
			MethodName: "PlumberDisconnect",
			Handler:    _Plumber_PlumberDisconnect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plumber.proto",
}
