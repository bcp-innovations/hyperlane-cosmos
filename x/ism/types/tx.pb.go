// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/ism/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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

type MsgCreateMultisigIsm struct {
	// sender is the message sender.
	Creator  string       `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	MultiSig *MultiSigIsm `protobuf:"bytes,2,opt,name=multi_sig,json=multiSig,proto3" json:"multi_sig,omitempty"`
}

func (m *MsgCreateMultisigIsm) Reset()         { *m = MsgCreateMultisigIsm{} }
func (m *MsgCreateMultisigIsm) String() string { return proto.CompactTextString(m) }
func (*MsgCreateMultisigIsm) ProtoMessage()    {}
func (*MsgCreateMultisigIsm) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b013040feeda308, []int{0}
}
func (m *MsgCreateMultisigIsm) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateMultisigIsm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateMultisigIsm.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateMultisigIsm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateMultisigIsm.Merge(m, src)
}
func (m *MsgCreateMultisigIsm) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateMultisigIsm) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateMultisigIsm.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateMultisigIsm proto.InternalMessageInfo

func (m *MsgCreateMultisigIsm) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgCreateMultisigIsm) GetMultiSig() *MultiSigIsm {
	if m != nil {
		return m.MultiSig
	}
	return nil
}

type MsgCreateMultisigIsmResponse struct {
}

func (m *MsgCreateMultisigIsmResponse) Reset()         { *m = MsgCreateMultisigIsmResponse{} }
func (m *MsgCreateMultisigIsmResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreateMultisigIsmResponse) ProtoMessage()    {}
func (*MsgCreateMultisigIsmResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b013040feeda308, []int{1}
}
func (m *MsgCreateMultisigIsmResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateMultisigIsmResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateMultisigIsmResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateMultisigIsmResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateMultisigIsmResponse.Merge(m, src)
}
func (m *MsgCreateMultisigIsmResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateMultisigIsmResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateMultisigIsmResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateMultisigIsmResponse proto.InternalMessageInfo

type MsgCreateNoopIsm struct {
	// sender is the message sender.
	Creator string `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
}

func (m *MsgCreateNoopIsm) Reset()         { *m = MsgCreateNoopIsm{} }
func (m *MsgCreateNoopIsm) String() string { return proto.CompactTextString(m) }
func (*MsgCreateNoopIsm) ProtoMessage()    {}
func (*MsgCreateNoopIsm) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b013040feeda308, []int{2}
}
func (m *MsgCreateNoopIsm) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateNoopIsm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateNoopIsm.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateNoopIsm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateNoopIsm.Merge(m, src)
}
func (m *MsgCreateNoopIsm) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateNoopIsm) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateNoopIsm.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateNoopIsm proto.InternalMessageInfo

func (m *MsgCreateNoopIsm) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

type MsgCreateNoopIsmResponse struct {
}

func (m *MsgCreateNoopIsmResponse) Reset()         { *m = MsgCreateNoopIsmResponse{} }
func (m *MsgCreateNoopIsmResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreateNoopIsmResponse) ProtoMessage()    {}
func (*MsgCreateNoopIsmResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b013040feeda308, []int{3}
}
func (m *MsgCreateNoopIsmResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateNoopIsmResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateNoopIsmResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateNoopIsmResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateNoopIsmResponse.Merge(m, src)
}
func (m *MsgCreateNoopIsmResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateNoopIsmResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateNoopIsmResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateNoopIsmResponse proto.InternalMessageInfo

// MsgUpdateParams is the Msg/UpdateParams request type.
type MsgUpdateParams struct {
	// authority is the address that controls the module
	// NOTE: Defaults to the governance module unless overwritten.
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	// params defines the module parameters to update.
	// NOTE: All parameters must be supplied.
	Params Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}

func (m *MsgUpdateParams) Reset()         { *m = MsgUpdateParams{} }
func (m *MsgUpdateParams) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParams) ProtoMessage()    {}
func (*MsgUpdateParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b013040feeda308, []int{4}
}
func (m *MsgUpdateParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParams.Merge(m, src)
}
func (m *MsgUpdateParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParams proto.InternalMessageInfo

func (m *MsgUpdateParams) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *MsgUpdateParams) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

// MsgUpdateParamsResponse defines the response structure for executing a
// MsgUpdateParams message.
type MsgUpdateParamsResponse struct {
}

func (m *MsgUpdateParamsResponse) Reset()         { *m = MsgUpdateParamsResponse{} }
func (m *MsgUpdateParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateParamsResponse) ProtoMessage()    {}
func (*MsgUpdateParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_7b013040feeda308, []int{5}
}
func (m *MsgUpdateParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateParamsResponse.Merge(m, src)
}
func (m *MsgUpdateParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgCreateMultisigIsm)(nil), "hyperlane.ism.v1.MsgCreateMultisigIsm")
	proto.RegisterType((*MsgCreateMultisigIsmResponse)(nil), "hyperlane.ism.v1.MsgCreateMultisigIsmResponse")
	proto.RegisterType((*MsgCreateNoopIsm)(nil), "hyperlane.ism.v1.MsgCreateNoopIsm")
	proto.RegisterType((*MsgCreateNoopIsmResponse)(nil), "hyperlane.ism.v1.MsgCreateNoopIsmResponse")
	proto.RegisterType((*MsgUpdateParams)(nil), "hyperlane.ism.v1.MsgUpdateParams")
	proto.RegisterType((*MsgUpdateParamsResponse)(nil), "hyperlane.ism.v1.MsgUpdateParamsResponse")
}

func init() { proto.RegisterFile("hyperlane/ism/v1/tx.proto", fileDescriptor_7b013040feeda308) }

var fileDescriptor_7b013040feeda308 = []byte{
	// 495 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xbf, 0x6f, 0xd3, 0x40,
	0x18, 0x8d, 0x8b, 0x28, 0xe4, 0x00, 0xd1, 0x5a, 0x95, 0xea, 0x58, 0xa9, 0x29, 0x1e, 0x50, 0x89,
	0x54, 0x1f, 0x6d, 0x07, 0xa4, 0xb0, 0x95, 0x89, 0x21, 0x08, 0x82, 0x58, 0x10, 0xa8, 0xba, 0x84,
	0xe3, 0x72, 0x22, 0xe7, 0xcf, 0xf2, 0x77, 0x89, 0x9a, 0x0d, 0x31, 0x32, 0xf1, 0x17, 0x20, 0xb1,
	0x31, 0x66, 0xe4, 0x4f, 0xe8, 0xd8, 0x91, 0x09, 0xa1, 0x64, 0xc8, 0xbf, 0x81, 0x7c, 0x39, 0xbb,
	0xcd, 0x0f, 0x1a, 0x96, 0x28, 0xf7, 0xbd, 0xf7, 0xbd, 0xf7, 0x7c, 0x4f, 0x47, 0x2a, 0x9d, 0x41,
	0xc2, 0xd3, 0x2e, 0x8b, 0x39, 0x95, 0xa8, 0x68, 0xff, 0x80, 0xea, 0xd3, 0x28, 0x49, 0x41, 0x83,
	0xbb, 0x51, 0x40, 0x91, 0x44, 0x15, 0xf5, 0x0f, 0xfc, 0x2d, 0x01, 0x02, 0x0c, 0x48, 0xb3, 0x7f,
	0x53, 0x9e, 0xbf, 0xdd, 0x06, 0x54, 0x80, 0x54, 0xa1, 0xc8, 0xf6, 0x15, 0x0a, 0x0b, 0x6c, 0x32,
	0x25, 0x63, 0xa0, 0xe6, 0xd7, 0x8e, 0xaa, 0x8b, 0x76, 0x83, 0x84, 0xa3, 0x45, 0x2b, 0x02, 0x40,
	0x74, 0x39, 0x35, 0xa7, 0x56, 0xef, 0x03, 0x65, 0xf1, 0x60, 0x0a, 0x85, 0xdf, 0x1d, 0xb2, 0xd5,
	0x40, 0xf1, 0x34, 0xe5, 0x4c, 0xf3, 0x46, 0xaf, 0xab, 0x25, 0x4a, 0xf1, 0x0c, 0x95, 0xeb, 0x91,
	0x1b, 0xed, 0x6c, 0x08, 0xa9, 0xe7, 0xec, 0x3a, 0x7b, 0xe5, 0x66, 0x7e, 0x74, 0xeb, 0xa4, 0xac,
	0x32, 0xe2, 0x09, 0x4a, 0xe1, 0xad, 0xed, 0x3a, 0x7b, 0xb7, 0x0e, 0x77, 0xa2, 0xf9, 0x6f, 0x8a,
	0x8c, 0xd6, 0x2b, 0xa3, 0xd5, 0xbc, 0xa9, 0xec, 0xa1, 0x7e, 0xf4, 0x79, 0x32, 0xac, 0xe5, 0x4a,
	0x5f, 0x26, 0xc3, 0x5a, 0x38, 0x1b, 0x7c, 0x36, 0xca, 0x74, 0x3d, 0x0c, 0x48, 0x75, 0x59, 0xc4,
	0x26, 0xc7, 0x04, 0x62, 0xe4, 0xe1, 0x3b, 0xb2, 0x51, 0xe0, 0xcf, 0x01, 0x92, 0x2b, 0xe3, 0xd7,
	0xe9, 0x7c, 0x84, 0xe0, 0x1f, 0x11, 0xac, 0x54, 0xe8, 0x13, 0x6f, 0x7e, 0x56, 0x58, 0x7f, 0x73,
	0xc8, 0xdd, 0x06, 0x8a, 0xd7, 0xc9, 0x7b, 0xa6, 0xf9, 0x0b, 0x96, 0x32, 0x85, 0x6e, 0x95, 0x94,
	0x59, 0x4f, 0x77, 0x20, 0x95, 0x7a, 0x60, 0xcd, 0x2f, 0x06, 0xee, 0x13, 0xb2, 0x9e, 0x18, 0x9e,
	0xbd, 0x3a, 0x6f, 0xf1, 0xea, 0xa6, 0x3a, 0xc7, 0xe5, 0xb3, 0xdf, 0xf7, 0x4a, 0x3f, 0x26, 0xc3,
	0x9a, 0xd3, 0xb4, 0x2b, 0xf5, 0x47, 0x59, 0xf6, 0x0b, 0xb1, 0x2c, 0xfd, 0xce, 0x42, 0xfa, 0xcb,
	0x61, 0xc2, 0x0a, 0xd9, 0x9e, 0x1b, 0xe5, 0xd9, 0x0f, 0x7f, 0xae, 0x91, 0x6b, 0x0d, 0x14, 0xee,
	0x47, 0xb2, 0xb9, 0x58, 0xff, 0x83, 0x25, 0x8d, 0x2e, 0xe9, 0xc0, 0x8f, 0xfe, 0x8f, 0x97, 0x9b,
	0xba, 0x27, 0xe4, 0xce, 0x6c, 0x51, 0xe1, 0x15, 0x02, 0x96, 0xe3, 0xd7, 0x56, 0x73, 0x0a, 0x83,
	0xb7, 0xe4, 0xf6, 0x4c, 0x1b, 0xf7, 0x97, 0xee, 0x5e, 0xa6, 0xf8, 0x0f, 0x57, 0x52, 0x72, 0x75,
	0xff, 0xfa, 0xa7, 0xac, 0x8f, 0xe3, 0x97, 0x67, 0xa3, 0xc0, 0x39, 0x1f, 0x05, 0xce, 0x9f, 0x51,
	0xe0, 0x7c, 0x1d, 0x07, 0xa5, 0xf3, 0x71, 0x50, 0xfa, 0x35, 0x0e, 0x4a, 0x6f, 0x1e, 0x0b, 0xa9,
	0x3b, 0xbd, 0x56, 0xd4, 0x06, 0x45, 0x5b, 0xed, 0x64, 0x5f, 0xc6, 0x31, 0xf4, 0x99, 0x96, 0x10,
	0x23, 0x2d, 0x5c, 0xf6, 0xed, 0xcb, 0x3e, 0x35, 0x95, 0x99, 0x97, 0xda, 0x5a, 0x37, 0xef, 0xf1,
	0xe8, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0b, 0x72, 0x5c, 0x52, 0x39, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	CreateMultisigIsm(ctx context.Context, in *MsgCreateMultisigIsm, opts ...grpc.CallOption) (*MsgCreateMultisigIsmResponse, error)
	CreateNoopIsm(ctx context.Context, in *MsgCreateNoopIsm, opts ...grpc.CallOption) (*MsgCreateNoopIsmResponse, error)
	// UpdateParams updates the module parameters.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateMultisigIsm(ctx context.Context, in *MsgCreateMultisigIsm, opts ...grpc.CallOption) (*MsgCreateMultisigIsmResponse, error) {
	out := new(MsgCreateMultisigIsmResponse)
	err := c.cc.Invoke(ctx, "/hyperlane.ism.v1.Msg/CreateMultisigIsm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateNoopIsm(ctx context.Context, in *MsgCreateNoopIsm, opts ...grpc.CallOption) (*MsgCreateNoopIsmResponse, error) {
	out := new(MsgCreateNoopIsmResponse)
	err := c.cc.Invoke(ctx, "/hyperlane.ism.v1.Msg/CreateNoopIsm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/hyperlane.ism.v1.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	CreateMultisigIsm(context.Context, *MsgCreateMultisigIsm) (*MsgCreateMultisigIsmResponse, error)
	CreateNoopIsm(context.Context, *MsgCreateNoopIsm) (*MsgCreateNoopIsmResponse, error)
	// UpdateParams updates the module parameters.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) CreateMultisigIsm(ctx context.Context, req *MsgCreateMultisigIsm) (*MsgCreateMultisigIsmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMultisigIsm not implemented")
}
func (*UnimplementedMsgServer) CreateNoopIsm(ctx context.Context, req *MsgCreateNoopIsm) (*MsgCreateNoopIsmResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateNoopIsm not implemented")
}
func (*UnimplementedMsgServer) UpdateParams(ctx context.Context, req *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateMultisigIsm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateMultisigIsm)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateMultisigIsm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hyperlane.ism.v1.Msg/CreateMultisigIsm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateMultisigIsm(ctx, req.(*MsgCreateMultisigIsm))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateNoopIsm_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateNoopIsm)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateNoopIsm(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hyperlane.ism.v1.Msg/CreateNoopIsm",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateNoopIsm(ctx, req.(*MsgCreateNoopIsm))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hyperlane.ism.v1.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "hyperlane.ism.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMultisigIsm",
			Handler:    _Msg_CreateMultisigIsm_Handler,
		},
		{
			MethodName: "CreateNoopIsm",
			Handler:    _Msg_CreateNoopIsm_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "hyperlane/ism/v1/tx.proto",
}

func (m *MsgCreateMultisigIsm) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateMultisigIsm) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateMultisigIsm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MultiSig != nil {
		{
			size, err := m.MultiSig.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateMultisigIsmResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateMultisigIsmResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateMultisigIsmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgCreateNoopIsm) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateNoopIsm) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateNoopIsm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateNoopIsmResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateNoopIsmResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateNoopIsmResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgCreateMultisigIsm) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.MultiSig != nil {
		l = m.MultiSig.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgCreateMultisigIsmResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgCreateNoopIsm) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgCreateNoopIsmResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgUpdateParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Params.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgUpdateParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgCreateMultisigIsm) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgCreateMultisigIsm: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateMultisigIsm: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiSig", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.MultiSig == nil {
				m.MultiSig = &MultiSigIsm{}
			}
			if err := m.MultiSig.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgCreateMultisigIsmResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgCreateMultisigIsmResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateMultisigIsmResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgCreateNoopIsm) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgCreateNoopIsm: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateNoopIsm: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgCreateNoopIsmResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgCreateNoopIsmResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateNoopIsmResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgUpdateParams) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdateParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgUpdateParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgUpdateParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
