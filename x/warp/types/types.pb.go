// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/warp/v1/types.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	github_com_bcp_innovations_hyperlane_cosmos_util "github.com/bcp-innovations/hyperlane-cosmos/util"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// HypTokenType ...
type HypTokenType int32

const (
	// HYP_TOKEN_TYPE_UNSPECIFIED ...
	HYP_TOKEN_TYPE_UNSPECIFIED HypTokenType = 0
	// HYP_TOKEN_TYPE_COLLATERAL ...
	HYP_TOKEN_TYPE_COLLATERAL HypTokenType = 1
	// HYP_TOKEN_TYPE_SYNTHETIC ...
	HYP_TOKEN_TYPE_SYNTHETIC HypTokenType = 2
)

var HypTokenType_name = map[int32]string{
	0: "HYP_TOKEN_TYPE_UNSPECIFIED",
	1: "HYP_TOKEN_TYPE_COLLATERAL",
	2: "HYP_TOKEN_TYPE_SYNTHETIC",
}

var HypTokenType_value = map[string]int32{
	"HYP_TOKEN_TYPE_UNSPECIFIED": 0,
	"HYP_TOKEN_TYPE_COLLATERAL":  1,
	"HYP_TOKEN_TYPE_SYNTHETIC":   2,
}

func (x HypTokenType) String() string {
	return proto.EnumName(HypTokenType_name, int32(x))
}

func (HypTokenType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_7372986c61417e18, []int{0}
}

// Params defines the parameters of the module.
type Params struct {
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_7372986c61417e18, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params Params     `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Tokens []HypToken `protobuf:"bytes,2,rep,name=tokens,proto3" json:"tokens"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_7372986c61417e18, []int{1}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetTokens() []HypToken {
	if m != nil {
		return m.Tokens
	}
	return nil
}

// HypToken ...
type HypToken struct {
	Id                github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress `protobuf:"bytes,1,opt,name=id,proto3,customtype=github.com/bcp-innovations/hyperlane-cosmos/util.HexAddress" json:"id"`
	Owner             string                                                      `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	TokenType         HypTokenType                                                `protobuf:"varint,3,opt,name=token_type,json=tokenType,proto3,enum=hyperlane.warp.v1.HypTokenType" json:"token_type,omitempty"`
	OriginMailbox     []byte                                                      `protobuf:"bytes,4,opt,name=origin_mailbox,json=originMailbox,proto3" json:"origin_mailbox,omitempty"`
	OriginDenom       string                                                      `protobuf:"bytes,5,opt,name=origin_denom,json=originDenom,proto3" json:"origin_denom,omitempty"`
	CollateralBalance cosmossdk_io_math.Int                                       `protobuf:"bytes,6,opt,name=collateral_balance,json=collateralBalance,proto3,customtype=cosmossdk.io/math.Int" json:"collateral_balance"`
	IsmId             string                                                      `protobuf:"bytes,7,opt,name=ism_id,json=ismId,proto3" json:"ism_id,omitempty"`
}

func (m *HypToken) Reset()         { *m = HypToken{} }
func (m *HypToken) String() string { return proto.CompactTextString(m) }
func (*HypToken) ProtoMessage()    {}
func (*HypToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_7372986c61417e18, []int{2}
}
func (m *HypToken) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *HypToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_HypToken.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *HypToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HypToken.Merge(m, src)
}
func (m *HypToken) XXX_Size() int {
	return m.Size()
}
func (m *HypToken) XXX_DiscardUnknown() {
	xxx_messageInfo_HypToken.DiscardUnknown(m)
}

var xxx_messageInfo_HypToken proto.InternalMessageInfo

func (m *HypToken) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *HypToken) GetTokenType() HypTokenType {
	if m != nil {
		return m.TokenType
	}
	return HYP_TOKEN_TYPE_UNSPECIFIED
}

func (m *HypToken) GetOriginMailbox() []byte {
	if m != nil {
		return m.OriginMailbox
	}
	return nil
}

func (m *HypToken) GetOriginDenom() string {
	if m != nil {
		return m.OriginDenom
	}
	return ""
}

func (m *HypToken) GetIsmId() string {
	if m != nil {
		return m.IsmId
	}
	return ""
}

// RemoteRouter ...
type RemoteRouter struct {
	ReceiverDomain   uint32                `protobuf:"varint,1,opt,name=receiver_domain,json=receiverDomain,proto3" json:"receiver_domain,omitempty"`
	ReceiverContract string                `protobuf:"bytes,2,opt,name=receiver_contract,json=receiverContract,proto3" json:"receiver_contract,omitempty"`
	Gas              cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=gas,proto3,customtype=cosmossdk.io/math.Int" json:"gas"`
}

func (m *RemoteRouter) Reset()         { *m = RemoteRouter{} }
func (m *RemoteRouter) String() string { return proto.CompactTextString(m) }
func (*RemoteRouter) ProtoMessage()    {}
func (*RemoteRouter) Descriptor() ([]byte, []int) {
	return fileDescriptor_7372986c61417e18, []int{3}
}
func (m *RemoteRouter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RemoteRouter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RemoteRouter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RemoteRouter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteRouter.Merge(m, src)
}
func (m *RemoteRouter) XXX_Size() int {
	return m.Size()
}
func (m *RemoteRouter) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteRouter.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteRouter proto.InternalMessageInfo

func (m *RemoteRouter) GetReceiverDomain() uint32 {
	if m != nil {
		return m.ReceiverDomain
	}
	return 0
}

func (m *RemoteRouter) GetReceiverContract() string {
	if m != nil {
		return m.ReceiverContract
	}
	return ""
}

func init() {
	proto.RegisterEnum("hyperlane.warp.v1.HypTokenType", HypTokenType_name, HypTokenType_value)
	proto.RegisterType((*Params)(nil), "hyperlane.warp.v1.Params")
	proto.RegisterType((*GenesisState)(nil), "hyperlane.warp.v1.GenesisState")
	proto.RegisterType((*HypToken)(nil), "hyperlane.warp.v1.HypToken")
	proto.RegisterType((*RemoteRouter)(nil), "hyperlane.warp.v1.RemoteRouter")
}

func init() { proto.RegisterFile("hyperlane/warp/v1/types.proto", fileDescriptor_7372986c61417e18) }

var fileDescriptor_7372986c61417e18 = []byte{
	// 642 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0xcf, 0x4e, 0x13, 0x41,
	0x1c, 0xee, 0x16, 0xa8, 0x32, 0x14, 0x6c, 0x27, 0x90, 0x2c, 0x55, 0x96, 0xda, 0x68, 0x6c, 0x30,
	0xec, 0x06, 0xbc, 0x18, 0x35, 0x24, 0xb4, 0x54, 0xdb, 0x58, 0xa1, 0xd9, 0xae, 0x07, 0xbc, 0x6c,
	0xa6, 0xbb, 0x93, 0x32, 0x61, 0x77, 0x66, 0xdd, 0x19, 0x0a, 0x7d, 0x03, 0x63, 0x3c, 0x78, 0xf1,
	0x09, 0xbc, 0x78, 0xe4, 0xe0, 0x43, 0x70, 0x24, 0x9e, 0x8c, 0x07, 0x62, 0xe0, 0xc0, 0x5b, 0x18,
	0xb3, 0x33, 0x5b, 0x54, 0x34, 0x1a, 0x2f, 0x9b, 0xd9, 0xef, 0xcf, 0xce, 0x37, 0xdf, 0xec, 0x0f,
	0x2c, 0xec, 0x0c, 0x23, 0x1c, 0x07, 0x88, 0x62, 0x6b, 0x1f, 0xc5, 0x91, 0x35, 0x58, 0xb1, 0xc4,
	0x30, 0xc2, 0xdc, 0x8c, 0x62, 0x26, 0x18, 0x2c, 0x5e, 0xd0, 0x66, 0x42, 0x9b, 0x83, 0x95, 0xd2,
	0xbc, 0xc7, 0x78, 0xc8, 0xb8, 0x2b, 0x05, 0x96, 0x7a, 0x51, 0xea, 0xd2, 0x6c, 0x9f, 0xf5, 0x99,
	0xc2, 0x93, 0x55, 0x8a, 0x16, 0x51, 0x48, 0x28, 0xb3, 0xe4, 0x53, 0x41, 0x95, 0x5b, 0x20, 0xd7,
	0x41, 0x31, 0x0a, 0xf9, 0x83, 0xd2, 0xeb, 0xf3, 0xc3, 0xa5, 0xb9, 0x4b, 0x21, 0x14, 0x57, 0x79,
	0xa3, 0x81, 0xfc, 0x13, 0x4c, 0x31, 0x27, 0xbc, 0x2b, 0x90, 0xc0, 0xf0, 0x11, 0xc8, 0x45, 0x92,
	0xd2, 0xb5, 0xb2, 0x56, 0x9d, 0x5a, 0x9d, 0x37, 0x7f, 0x8b, 0x67, 0x2a, 0x6f, 0x6d, 0xf2, 0xe8,
	0x64, 0x31, 0xf3, 0xe1, 0xfc, 0x70, 0x49, 0xb3, 0x53, 0x0f, 0x5c, 0x03, 0x39, 0xc1, 0x76, 0x31,
	0xe5, 0x7a, 0xb6, 0x3c, 0x56, 0x9d, 0x5a, 0xbd, 0xfe, 0x07, 0x77, 0x73, 0x18, 0x39, 0x89, 0xe6,
	0x17, 0xbf, 0x72, 0x55, 0xbe, 0x65, 0xc1, 0xd5, 0x11, 0x0f, 0xbb, 0x20, 0x4b, 0x7c, 0x19, 0x63,
	0xb2, 0x56, 0x4f, 0xb4, 0x5f, 0x4e, 0x16, 0x1f, 0xf6, 0x89, 0xd8, 0xd9, 0xeb, 0x99, 0x1e, 0x0b,
	0xad, 0x9e, 0x17, 0x2d, 0x13, 0x4a, 0xd9, 0x00, 0x09, 0xc2, 0x28, 0xb7, 0x2e, 0xb6, 0x5a, 0x56,
	0x8d, 0x59, 0x7b, 0x82, 0x04, 0x66, 0x13, 0x1f, 0xac, 0xfb, 0x7e, 0x8c, 0x39, 0xb7, 0xb3, 0xc4,
	0x87, 0x26, 0x98, 0x60, 0xfb, 0x14, 0xc7, 0x7a, 0x56, 0x7e, 0x57, 0xff, 0xf4, 0x71, 0x79, 0x36,
	0x2d, 0x38, 0x95, 0x75, 0x45, 0x4c, 0x68, 0xdf, 0x56, 0x32, 0xb8, 0x06, 0x80, 0xcc, 0xe6, 0x26,
	0x57, 0xa6, 0x8f, 0x95, 0xb5, 0xea, 0xcc, 0xea, 0xe2, 0x5f, 0x4e, 0xe5, 0x0c, 0x23, 0x6c, 0x4f,
	0x8a, 0xd1, 0x12, 0xde, 0x06, 0x33, 0x2c, 0x26, 0x7d, 0x42, 0xdd, 0x10, 0x91, 0xa0, 0xc7, 0x0e,
	0xf4, 0xf1, 0xb2, 0x56, 0xcd, 0xdb, 0xd3, 0x0a, 0x7d, 0xa6, 0x40, 0x78, 0x13, 0xe4, 0x53, 0x99,
	0x8f, 0x29, 0x0b, 0xf5, 0x89, 0x24, 0x9d, 0x3d, 0xa5, 0xb0, 0x8d, 0x04, 0x82, 0x6d, 0x00, 0x3d,
	0x16, 0x04, 0x48, 0xe0, 0x18, 0x05, 0x6e, 0x0f, 0x05, 0x88, 0x7a, 0x58, 0xcf, 0xc9, 0x63, 0x2c,
	0xa4, 0xf5, 0xcc, 0xa9, 0xa3, 0x70, 0x7f, 0xd7, 0x24, 0xcc, 0x0a, 0x91, 0xd8, 0x31, 0x5b, 0x54,
	0xd8, 0xc5, 0x1f, 0xc6, 0x9a, 0xf2, 0xc1, 0x39, 0x90, 0x23, 0x3c, 0x74, 0x89, 0xaf, 0x5f, 0x91,
	0x5b, 0x4d, 0x10, 0x1e, 0xb6, 0xfc, 0xca, 0x3b, 0x0d, 0xe4, 0x6d, 0x1c, 0x32, 0x81, 0x6d, 0xb6,
	0x27, 0x70, 0x0c, 0xef, 0x80, 0x6b, 0x31, 0xf6, 0x30, 0x19, 0xe0, 0xd8, 0xf5, 0x59, 0x88, 0x08,
	0x95, 0x37, 0x32, 0x6d, 0xcf, 0x8c, 0xe0, 0x0d, 0x89, 0xc2, 0xbb, 0xa0, 0x78, 0x21, 0xf4, 0x18,
	0x15, 0x31, 0xf2, 0x84, 0x2a, 0xd9, 0x2e, 0x8c, 0x88, 0x7a, 0x8a, 0x43, 0x0b, 0x8c, 0xf5, 0x11,
	0x97, 0x75, 0xfe, 0x33, 0x7c, 0xa2, 0x5c, 0x7a, 0x09, 0xf2, 0x3f, 0x37, 0x0c, 0x0d, 0x50, 0x6a,
	0x6e, 0x77, 0x5c, 0x67, 0xeb, 0x69, 0x63, 0xd3, 0x75, 0xb6, 0x3b, 0x0d, 0xf7, 0xf9, 0x66, 0xb7,
	0xd3, 0xa8, 0xb7, 0x1e, 0xb7, 0x1a, 0x1b, 0x85, 0x0c, 0x5c, 0x00, 0xf3, 0x97, 0xf8, 0xfa, 0x56,
	0xbb, 0xbd, 0xee, 0x34, 0xec, 0xf5, 0x76, 0x41, 0x83, 0x37, 0x80, 0x7e, 0x89, 0xee, 0x6e, 0x6f,
	0x3a, 0xcd, 0x86, 0xd3, 0xaa, 0x17, 0xb2, 0xa5, 0xf1, 0x57, 0xef, 0x8d, 0x4c, 0xcd, 0x3e, 0x3a,
	0x35, 0xb4, 0xe3, 0x53, 0x43, 0xfb, 0x7a, 0x6a, 0x68, 0x6f, 0xcf, 0x8c, 0xcc, 0xf1, 0x99, 0x91,
	0xf9, 0x7c, 0x66, 0x64, 0x5e, 0xdc, 0xff, 0x9f, 0x9f, 0xf0, 0x40, 0xcd, 0x9b, 0x9c, 0xf8, 0x5e,
	0x4e, 0xce, 0xe6, 0xbd, 0xef, 0x01, 0x00, 0x00, 0xff, 0xff, 0xed, 0xc9, 0xc9, 0xcc, 0x13, 0x04,
	0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Tokens) > 0 {
		for iNdEx := len(m.Tokens) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Tokens[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *HypToken) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *HypToken) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *HypToken) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.IsmId) > 0 {
		i -= len(m.IsmId)
		copy(dAtA[i:], m.IsmId)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.IsmId)))
		i--
		dAtA[i] = 0x3a
	}
	{
		size := m.CollateralBalance.Size()
		i -= size
		if _, err := m.CollateralBalance.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.OriginDenom) > 0 {
		i -= len(m.OriginDenom)
		copy(dAtA[i:], m.OriginDenom)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.OriginDenom)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.OriginMailbox) > 0 {
		i -= len(m.OriginMailbox)
		copy(dAtA[i:], m.OriginMailbox)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.OriginMailbox)))
		i--
		dAtA[i] = 0x22
	}
	if m.TokenType != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.TokenType))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.Id.Size()
		i -= size
		if _, err := m.Id.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *RemoteRouter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RemoteRouter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RemoteRouter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Gas.Size()
		i -= size
		if _, err := m.Gas.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.ReceiverContract) > 0 {
		i -= len(m.ReceiverContract)
		copy(dAtA[i:], m.ReceiverContract)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.ReceiverContract)))
		i--
		dAtA[i] = 0x12
	}
	if m.ReceiverDomain != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.ReceiverDomain))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovTypes(uint64(l))
	if len(m.Tokens) > 0 {
		for _, e := range m.Tokens {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func (m *HypToken) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Id.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.TokenType != 0 {
		n += 1 + sovTypes(uint64(m.TokenType))
	}
	l = len(m.OriginMailbox)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = len(m.OriginDenom)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.CollateralBalance.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = len(m.IsmId)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *RemoteRouter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ReceiverDomain != 0 {
		n += 1 + sovTypes(uint64(m.ReceiverDomain))
	}
	l = len(m.ReceiverContract)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.Gas.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tokens", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Tokens = append(m.Tokens, HypToken{})
			if err := m.Tokens[len(m.Tokens)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *HypToken) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: HypToken: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: HypToken: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Id.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenType", wireType)
			}
			m.TokenType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TokenType |= HypTokenType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginMailbox", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OriginMailbox = append(m.OriginMailbox[:0], dAtA[iNdEx:postIndex]...)
			if m.OriginMailbox == nil {
				m.OriginMailbox = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OriginDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralBalance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.CollateralBalance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsmId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IsmId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *RemoteRouter) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: RemoteRouter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RemoteRouter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiverDomain", wireType)
			}
			m.ReceiverDomain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ReceiverDomain |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceiverContract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceiverContract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gas", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Gas.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
