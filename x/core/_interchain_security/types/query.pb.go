// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/interchain_security/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryIsmsRequest ...
type QueryIsmsRequest struct {
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryIsmsRequest) Reset()         { *m = QueryIsmsRequest{} }
func (m *QueryIsmsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryIsmsRequest) ProtoMessage()    {}
func (*QueryIsmsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d5cff9b810eaec0b, []int{0}
}
func (m *QueryIsmsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIsmsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIsmsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIsmsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIsmsRequest.Merge(m, src)
}
func (m *QueryIsmsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryIsmsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIsmsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIsmsRequest proto.InternalMessageInfo

func (m *QueryIsmsRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryIsmsResponse ...
type QueryIsmsResponse struct {
	Isms []*types.Any `protobuf:"bytes,1,rep,name=isms,proto3" json:"isms,omitempty"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryIsmsResponse) Reset()         { *m = QueryIsmsResponse{} }
func (m *QueryIsmsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryIsmsResponse) ProtoMessage()    {}
func (*QueryIsmsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d5cff9b810eaec0b, []int{1}
}
func (m *QueryIsmsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIsmsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIsmsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIsmsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIsmsResponse.Merge(m, src)
}
func (m *QueryIsmsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryIsmsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIsmsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIsmsResponse proto.InternalMessageInfo

func (m *QueryIsmsResponse) GetIsms() []*types.Any {
	if m != nil {
		return m.Isms
	}
	return nil
}

func (m *QueryIsmsResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryIsmRequest ...
type QueryIsmRequest struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *QueryIsmRequest) Reset()         { *m = QueryIsmRequest{} }
func (m *QueryIsmRequest) String() string { return proto.CompactTextString(m) }
func (*QueryIsmRequest) ProtoMessage()    {}
func (*QueryIsmRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d5cff9b810eaec0b, []int{2}
}
func (m *QueryIsmRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIsmRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIsmRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIsmRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIsmRequest.Merge(m, src)
}
func (m *QueryIsmRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryIsmRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIsmRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIsmRequest proto.InternalMessageInfo

func (m *QueryIsmRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// QueryIsmResponse ...
type QueryIsmResponse struct {
	Ism types.Any `protobuf:"bytes,1,opt,name=ism,proto3" json:"ism"`
}

func (m *QueryIsmResponse) Reset()         { *m = QueryIsmResponse{} }
func (m *QueryIsmResponse) String() string { return proto.CompactTextString(m) }
func (*QueryIsmResponse) ProtoMessage()    {}
func (*QueryIsmResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d5cff9b810eaec0b, []int{3}
}
func (m *QueryIsmResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIsmResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIsmResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIsmResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIsmResponse.Merge(m, src)
}
func (m *QueryIsmResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryIsmResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIsmResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIsmResponse proto.InternalMessageInfo

func (m *QueryIsmResponse) GetIsm() types.Any {
	if m != nil {
		return m.Ism
	}
	return types.Any{}
}

func init() {
	proto.RegisterType((*QueryIsmsRequest)(nil), "hyperlane.core.interchain_security.v1.QueryIsmsRequest")
	proto.RegisterType((*QueryIsmsResponse)(nil), "hyperlane.core.interchain_security.v1.QueryIsmsResponse")
	proto.RegisterType((*QueryIsmRequest)(nil), "hyperlane.core.interchain_security.v1.QueryIsmRequest")
	proto.RegisterType((*QueryIsmResponse)(nil), "hyperlane.core.interchain_security.v1.QueryIsmResponse")
}

func init() {
	proto.RegisterFile("hyperlane/core/interchain_security/v1/query.proto", fileDescriptor_d5cff9b810eaec0b)
}

var fileDescriptor_d5cff9b810eaec0b = []byte{
	// 532 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4f, 0x6b, 0x13, 0x41,
	0x1c, 0xcd, 0x6c, 0xab, 0xd0, 0x29, 0xa8, 0x1d, 0x02, 0xb6, 0x41, 0xb6, 0x35, 0xa0, 0x96, 0x42,
	0x66, 0xd8, 0x08, 0xd6, 0xab, 0x05, 0xff, 0xe4, 0x50, 0xd4, 0xed, 0xad, 0x97, 0x30, 0xbb, 0x19,
	0x37, 0x03, 0xd9, 0x99, 0xed, 0xce, 0xee, 0xe2, 0x22, 0xbd, 0xf8, 0x09, 0x04, 0x11, 0xbf, 0x82,
	0x47, 0x0f, 0x5e, 0xfc, 0x06, 0xc5, 0x53, 0xc1, 0x8b, 0x20, 0x88, 0x24, 0x82, 0x5f, 0x43, 0x76,
	0x66, 0x36, 0x69, 0x6b, 0x40, 0x73, 0x09, 0x33, 0xf3, 0xfb, 0xbd, 0x97, 0xf7, 0xde, 0xcc, 0x6f,
	0xa1, 0x37, 0x2c, 0x13, 0x96, 0x8e, 0xa8, 0x60, 0x24, 0x94, 0x29, 0x23, 0x5c, 0x64, 0x2c, 0x0d,
	0x87, 0x94, 0x8b, 0xbe, 0x62, 0x61, 0x9e, 0xf2, 0xac, 0x24, 0x85, 0x47, 0x8e, 0x72, 0x96, 0x96,
	0x38, 0x49, 0x65, 0x26, 0xd1, 0xad, 0x29, 0x04, 0x57, 0x10, 0x3c, 0x07, 0x82, 0x0b, 0xaf, 0xb5,
	0x13, 0x4a, 0x15, 0x4b, 0x45, 0x02, 0xaa, 0x98, 0xc1, 0x93, 0xc2, 0x0b, 0x58, 0x46, 0x3d, 0x92,
	0xd0, 0x88, 0x0b, 0x9a, 0x71, 0x29, 0x0c, 0x65, 0xeb, 0x46, 0x24, 0x65, 0x34, 0x62, 0x84, 0x26,
	0x9c, 0x50, 0x21, 0x64, 0xa6, 0x8b, 0xca, 0x56, 0xd7, 0x68, 0xcc, 0x85, 0x24, 0xfa, 0xd7, 0x1e,
	0x35, 0x23, 0x19, 0x49, 0xbd, 0x24, 0xd5, 0xca, 0x9e, 0x6e, 0x58, 0x1a, 0xbd, 0x0b, 0xf2, 0x17,
	0x84, 0x8a, 0xb2, 0x2e, 0x19, 0x35, 0x7d, 0x83, 0x31, 0x1b, 0x53, 0x6a, 0x1f, 0xc2, 0x6b, 0xcf,
	0x2b, 0x79, 0x3d, 0x15, 0x2b, 0x9f, 0x1d, 0xe5, 0x4c, 0x65, 0xe8, 0x11, 0x84, 0x33, 0x91, 0xeb,
	0x60, 0x0b, 0x6c, 0xaf, 0x76, 0x6f, 0x63, 0x0b, 0xab, 0x1c, 0x61, 0x93, 0x88, 0x75, 0x84, 0x9f,
	0xd1, 0x88, 0x59, 0xac, 0x7f, 0x06, 0xd9, 0xfe, 0x0e, 0xe0, 0xda, 0x19, 0x72, 0x95, 0x48, 0xa1,
	0x18, 0x3a, 0x86, 0xcb, 0x5c, 0xc5, 0x6a, 0x1d, 0x6c, 0x2d, 0x6d, 0xaf, 0x76, 0x9b, 0xd8, 0xc8,
	0xc6, 0xb5, 0x6c, 0xfc, 0x40, 0x94, 0x7b, 0x07, 0x5f, 0x3e, 0x75, 0x9e, 0x5e, 0x48, 0xba, 0xf0,
	0x70, 0x7f, 0x4e, 0xda, 0xfd, 0x58, 0x0e, 0xf2, 0x11, 0xc3, 0x4f, 0xea, 0xee, 0xde, 0xb4, 0xe7,
	0xc0, 0xb6, 0xec, 0xeb, 0x0e, 0x5f, 0xff, 0x2d, 0x7a, 0x7c, 0xce, 0x9c, 0xa3, 0xcd, 0xdd, 0xf9,
	0xa7, 0x39, 0xa3, 0xfd, 0x9c, 0xbb, 0x9b, 0xf0, 0x6a, 0x6d, 0xae, 0x0e, 0xee, 0x0a, 0x74, 0xf8,
	0x40, 0x07, 0xb6, 0xe2, 0x3b, 0x7c, 0xd0, 0x7e, 0x38, 0x0b, 0x77, 0x6a, 0xdf, 0x83, 0x4b, 0x5c,
	0xc5, 0x36, 0xd5, 0xf9, 0xee, 0x57, 0x4e, 0x7e, 0x6c, 0x36, 0x3e, 0xfc, 0xfe, 0xb8, 0x03, 0xfc,
	0xaa, 0xb7, 0xfb, 0xd9, 0x81, 0x97, 0x34, 0x0f, 0x7a, 0x07, 0xe0, 0x72, 0x15, 0x26, 0xda, 0xc5,
	0xff, 0xf5, 0x0e, 0xf1, 0xc5, 0xbb, 0x6d, 0xdd, 0x5f, 0x1c, 0x68, 0x84, 0xb7, 0x5b, 0xaf, 0xbf,
	0xfe, 0x7a, 0xeb, 0x34, 0x11, 0x22, 0xb3, 0xa9, 0x29, 0x3c, 0xa2, 0x43, 0x7d, 0x0f, 0xe0, 0x52,
	0x4f, 0xc5, 0xe8, 0xde, 0x82, 0xec, 0xb5, 0xaa, 0xdd, 0x85, 0x71, 0x56, 0xd4, 0xa6, 0x16, 0xb5,
	0x81, 0xae, 0xff, 0x2d, 0x8a, 0xbc, 0xe2, 0x83, 0xe3, 0xbd, 0xe8, 0x64, 0xec, 0x82, 0xd3, 0xb1,
	0x0b, 0x7e, 0x8e, 0x5d, 0xf0, 0x66, 0xe2, 0x36, 0x4e, 0x27, 0x6e, 0xe3, 0xdb, 0xc4, 0x6d, 0x1c,
	0xee, 0x47, 0x3c, 0x1b, 0xe6, 0x01, 0x0e, 0x65, 0x4c, 0x82, 0x30, 0xe9, 0x70, 0x21, 0x64, 0x61,
	0x46, 0x6f, 0x46, 0xd6, 0xb1, 0x73, 0xfc, 0xd2, 0x7c, 0x20, 0xe6, 0x3d, 0x40, 0x92, 0x95, 0x09,
	0x53, 0xc1, 0x65, 0x7d, 0x85, 0x77, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x64, 0xc5, 0x6c, 0xc5,
	0x54, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Isms ...
	Isms(ctx context.Context, in *QueryIsmsRequest, opts ...grpc.CallOption) (*QueryIsmsResponse, error)
	// Ism ...
	Ism(ctx context.Context, in *QueryIsmRequest, opts ...grpc.CallOption) (*QueryIsmResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Isms(ctx context.Context, in *QueryIsmsRequest, opts ...grpc.CallOption) (*QueryIsmsResponse, error) {
	out := new(QueryIsmsResponse)
	err := c.cc.Invoke(ctx, "/hyperlane.core.interchain_security.v1.Query/Isms", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Ism(ctx context.Context, in *QueryIsmRequest, opts ...grpc.CallOption) (*QueryIsmResponse, error) {
	out := new(QueryIsmResponse)
	err := c.cc.Invoke(ctx, "/hyperlane.core.interchain_security.v1.Query/Ism", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Isms ...
	Isms(context.Context, *QueryIsmsRequest) (*QueryIsmsResponse, error)
	// Ism ...
	Ism(context.Context, *QueryIsmRequest) (*QueryIsmResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Isms(ctx context.Context, req *QueryIsmsRequest) (*QueryIsmsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Isms not implemented")
}
func (*UnimplementedQueryServer) Ism(ctx context.Context, req *QueryIsmRequest) (*QueryIsmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ism not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Isms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIsmsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Isms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hyperlane.core.interchain_security.v1.Query/Isms",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Isms(ctx, req.(*QueryIsmsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Ism_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIsmRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Ism(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hyperlane.core.interchain_security.v1.Query/Ism",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Ism(ctx, req.(*QueryIsmRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_serviceDesc = _Query_serviceDesc
var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "hyperlane.core.interchain_security.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Isms",
			Handler:    _Query_Isms_Handler,
		},
		{
			MethodName: "Ism",
			Handler:    _Query_Ism_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hyperlane/core/interchain_security/v1/query.proto",
}

func (m *QueryIsmsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIsmsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIsmsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryIsmsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIsmsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIsmsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Isms) > 0 {
		for iNdEx := len(m.Isms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Isms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryIsmRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIsmRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIsmRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryIsmResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIsmResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIsmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Ism.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryIsmsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryIsmsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Isms) > 0 {
		for _, e := range m.Isms {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryIsmRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryIsmResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Ism.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryIsmsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryIsmsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIsmsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryIsmsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryIsmsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIsmsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Isms", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Isms = append(m.Isms, &types.Any{})
			if err := m.Isms[len(m.Isms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryIsmRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryIsmRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIsmRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryIsmResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryIsmResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIsmResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ism", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Ism.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
