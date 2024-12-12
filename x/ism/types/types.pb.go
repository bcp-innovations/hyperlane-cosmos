// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/ism/v1/types.proto

package types

import (
	fmt "fmt"
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

// Params defines the parameters of the module.
type Params struct {
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_440bb1e52436f0d3, []int{0}
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
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_440bb1e52436f0d3, []int{1}
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

type Ism struct {
	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	IsmType uint32 `protobuf:"varint,2,opt,name=ism_type,json=ismType,proto3" json:"ism_type,omitempty"`
	Creator string `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	// Types that are valid to be assigned to Ism:
	//	*Ism_MultiSig
	//	*Ism_Noop
	Ism isIsm_Ism `protobuf_oneof:"ism"`
}

func (m *Ism) Reset()         { *m = Ism{} }
func (m *Ism) String() string { return proto.CompactTextString(m) }
func (*Ism) ProtoMessage()    {}
func (*Ism) Descriptor() ([]byte, []int) {
	return fileDescriptor_440bb1e52436f0d3, []int{2}
}
func (m *Ism) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Ism) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Ism.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Ism) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ism.Merge(m, src)
}
func (m *Ism) XXX_Size() int {
	return m.Size()
}
func (m *Ism) XXX_DiscardUnknown() {
	xxx_messageInfo_Ism.DiscardUnknown(m)
}

var xxx_messageInfo_Ism proto.InternalMessageInfo

type isIsm_Ism interface {
	isIsm_Ism()
	MarshalTo([]byte) (int, error)
	Size() int
}

type Ism_MultiSig struct {
	MultiSig *MultiSigIsm `protobuf:"bytes,4,opt,name=multi_sig,json=multiSig,proto3,oneof" json:"multi_sig,omitempty"`
}
type Ism_Noop struct {
	Noop *NoopIsm `protobuf:"bytes,5,opt,name=noop,proto3,oneof" json:"noop,omitempty"`
}

func (*Ism_MultiSig) isIsm_Ism() {}
func (*Ism_Noop) isIsm_Ism()     {}

func (m *Ism) GetIsm() isIsm_Ism {
	if m != nil {
		return m.Ism
	}
	return nil
}

func (m *Ism) GetMultiSig() *MultiSigIsm {
	if x, ok := m.GetIsm().(*Ism_MultiSig); ok {
		return x.MultiSig
	}
	return nil
}

func (m *Ism) GetNoop() *NoopIsm {
	if x, ok := m.GetIsm().(*Ism_Noop); ok {
		return x.Noop
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Ism) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Ism_MultiSig)(nil),
		(*Ism_Noop)(nil),
	}
}

type MultiSigIsm struct {
	ValidatorPubKeys []string `protobuf:"bytes,1,rep,name=validator_pub_keys,json=validatorPubKeys,proto3" json:"validator_pub_keys,omitempty"`
	Threshold        uint32   `protobuf:"varint,2,opt,name=threshold,proto3" json:"threshold,omitempty"`
}

func (m *MultiSigIsm) Reset()         { *m = MultiSigIsm{} }
func (m *MultiSigIsm) String() string { return proto.CompactTextString(m) }
func (*MultiSigIsm) ProtoMessage()    {}
func (*MultiSigIsm) Descriptor() ([]byte, []int) {
	return fileDescriptor_440bb1e52436f0d3, []int{3}
}
func (m *MultiSigIsm) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MultiSigIsm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MultiSigIsm.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MultiSigIsm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MultiSigIsm.Merge(m, src)
}
func (m *MultiSigIsm) XXX_Size() int {
	return m.Size()
}
func (m *MultiSigIsm) XXX_DiscardUnknown() {
	xxx_messageInfo_MultiSigIsm.DiscardUnknown(m)
}

var xxx_messageInfo_MultiSigIsm proto.InternalMessageInfo

func (m *MultiSigIsm) GetValidatorPubKeys() []string {
	if m != nil {
		return m.ValidatorPubKeys
	}
	return nil
}

func (m *MultiSigIsm) GetThreshold() uint32 {
	if m != nil {
		return m.Threshold
	}
	return 0
}

type NoopIsm struct {
}

func (m *NoopIsm) Reset()         { *m = NoopIsm{} }
func (m *NoopIsm) String() string { return proto.CompactTextString(m) }
func (*NoopIsm) ProtoMessage()    {}
func (*NoopIsm) Descriptor() ([]byte, []int) {
	return fileDescriptor_440bb1e52436f0d3, []int{4}
}
func (m *NoopIsm) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NoopIsm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NoopIsm.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NoopIsm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoopIsm.Merge(m, src)
}
func (m *NoopIsm) XXX_Size() int {
	return m.Size()
}
func (m *NoopIsm) XXX_DiscardUnknown() {
	xxx_messageInfo_NoopIsm.DiscardUnknown(m)
}

var xxx_messageInfo_NoopIsm proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "hyperlane.ism.v1.Params")
	proto.RegisterType((*GenesisState)(nil), "hyperlane.ism.v1.GenesisState")
	proto.RegisterType((*Ism)(nil), "hyperlane.ism.v1.Ism")
	proto.RegisterType((*MultiSigIsm)(nil), "hyperlane.ism.v1.MultiSigIsm")
	proto.RegisterType((*NoopIsm)(nil), "hyperlane.ism.v1.NoopIsm")
}

func init() { proto.RegisterFile("hyperlane/ism/v1/types.proto", fileDescriptor_440bb1e52436f0d3) }

var fileDescriptor_440bb1e52436f0d3 = []byte{
	// 458 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0x3f, 0x6f, 0xd3, 0x40,
	0x14, 0xb7, 0x93, 0x34, 0xa9, 0xaf, 0x80, 0xca, 0x29, 0x83, 0x53, 0x15, 0xb7, 0x32, 0x4b, 0x54,
	0x51, 0x9f, 0xda, 0x6e, 0x85, 0x85, 0x48, 0x08, 0xaa, 0xa8, 0x55, 0xe5, 0x20, 0xa4, 0xb2, 0x58,
	0x4e, 0x7c, 0x72, 0x4e, 0xcd, 0xf9, 0xac, 0x7b, 0x97, 0x80, 0xbf, 0x01, 0x62, 0xe2, 0x23, 0x30,
	0x32, 0x76, 0xe0, 0x43, 0x74, 0xac, 0x98, 0x98, 0x50, 0x95, 0x0c, 0xfd, 0x1a, 0xc8, 0x77, 0x26,
	0x40, 0xb3, 0x58, 0x7e, 0xbf, 0x3f, 0xcf, 0xbf, 0xf7, 0xfc, 0xd0, 0xf6, 0xb8, 0xc8, 0xa9, 0x9c,
	0xc4, 0x19, 0x25, 0x0c, 0x38, 0x99, 0x1d, 0x10, 0x55, 0xe4, 0x14, 0x82, 0x5c, 0x0a, 0x25, 0xf0,
	0xe6, 0x92, 0x0d, 0x18, 0xf0, 0x60, 0x76, 0xb0, 0xd5, 0x19, 0x09, 0xe0, 0x02, 0x22, 0xcd, 0x13,
	0x53, 0x18, 0xf1, 0x56, 0x3b, 0x15, 0xa9, 0x30, 0x78, 0xf9, 0x56, 0xa1, 0x8f, 0x63, 0xce, 0x32,
	0x41, 0xf4, 0xd3, 0x40, 0xfe, 0x53, 0xd4, 0x3c, 0x8f, 0x65, 0xcc, 0xe1, 0xb8, 0xf3, 0xf9, 0xee,
	0x6a, 0xaf, 0xfd, 0x7f, 0x04, 0x43, 0xf9, 0x7d, 0xf4, 0xe0, 0x35, 0xcd, 0x28, 0x30, 0x18, 0xa8,
	0x58, 0x51, 0xfc, 0x1c, 0x35, 0x73, 0xcd, 0xb8, 0xf6, 0xae, 0xdd, 0xdd, 0x38, 0x74, 0x83, 0xfb,
	0xd9, 0x02, 0xe3, 0xec, 0x39, 0xd7, 0xbf, 0x76, 0xac, 0x6f, 0x77, 0x57, 0x7b, 0x76, 0x58, 0x59,
	0xfc, 0x5b, 0x1b, 0xd5, 0x4f, 0x80, 0xe3, 0x47, 0xa8, 0xc6, 0x12, 0xdd, 0xc0, 0x09, 0x6b, 0x2c,
	0xc1, 0x1d, 0xb4, 0xce, 0x80, 0x47, 0xe5, 0xc8, 0x6e, 0x6d, 0xd7, 0xee, 0x3e, 0x0c, 0x5b, 0x0c,
	0xf8, 0xdb, 0x22, 0xa7, 0xf8, 0x10, 0xb5, 0x46, 0x92, 0xc6, 0x4a, 0x48, 0xb7, 0x5e, 0xea, 0x7b,
	0xee, 0x8f, 0xef, 0xfb, 0xed, 0x6a, 0xe0, 0x97, 0x49, 0x22, 0x29, 0xc0, 0x40, 0x49, 0x96, 0xa5,
	0xe1, 0x1f, 0x21, 0x7e, 0x81, 0x1c, 0x3e, 0x9d, 0x28, 0x16, 0x01, 0x4b, 0xdd, 0x86, 0x8e, 0xf9,
	0x64, 0x35, 0xe6, 0x69, 0x29, 0x19, 0xb0, 0xf4, 0x04, 0xf8, 0x1b, 0x2b, 0x5c, 0xe7, 0x55, 0x89,
	0x09, 0x6a, 0x64, 0x42, 0xe4, 0xee, 0x9a, 0x36, 0x76, 0x56, 0x8d, 0x67, 0x42, 0xe4, 0xc6, 0xa4,
	0x85, 0xc7, 0x8d, 0x4f, 0x5f, 0x77, 0xac, 0xde, 0x1a, 0xaa, 0x33, 0xe0, 0xfe, 0x05, 0xda, 0xf8,
	0xa7, 0x31, 0x7e, 0x86, 0xf0, 0x2c, 0x9e, 0xb0, 0xa4, 0xcc, 0x15, 0xe5, 0xd3, 0x61, 0x74, 0x49,
	0x8b, 0x72, 0x75, 0xf5, 0xae, 0x13, 0x6e, 0x2e, 0x99, 0xf3, 0xe9, 0xb0, 0x4f, 0x0b, 0xc0, 0xdb,
	0xc8, 0x51, 0x63, 0x49, 0x61, 0x2c, 0x26, 0x49, 0xb5, 0x88, 0xbf, 0x80, 0xef, 0xa0, 0x56, 0xf5,
	0xe9, 0xde, 0xe9, 0xf5, 0xdc, 0xb3, 0x6f, 0xe6, 0x9e, 0x7d, 0x3b, 0xf7, 0xec, 0x2f, 0x0b, 0xcf,
	0xba, 0x59, 0x78, 0xd6, 0xcf, 0x85, 0x67, 0xbd, 0x3f, 0x4a, 0x99, 0x1a, 0x4f, 0x87, 0xc1, 0x48,
	0x70, 0xd2, 0xbf, 0x78, 0xf7, 0xea, 0x8c, 0xaa, 0x0f, 0x42, 0x5e, 0x92, 0xe5, 0x14, 0xfb, 0x66,
	0x77, 0xe4, 0xa3, 0xfe, 0xcb, 0xfa, 0xca, 0x86, 0x4d, 0x7d, 0x10, 0x47, 0xbf, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x36, 0xe0, 0xa5, 0xb5, 0x86, 0x02, 0x00, 0x00,
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

func (m *Ism) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Ism) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Ism) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Ism != nil {
		{
			size := m.Ism.Size()
			i -= size
			if _, err := m.Ism.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x1a
	}
	if m.IsmType != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.IsmType))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Ism_MultiSig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Ism_MultiSig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.MultiSig != nil {
		{
			size, err := m.MultiSig.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	return len(dAtA) - i, nil
}
func (m *Ism_Noop) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Ism_Noop) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Noop != nil {
		{
			size, err := m.Noop.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	return len(dAtA) - i, nil
}
func (m *MultiSigIsm) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MultiSigIsm) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MultiSigIsm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Threshold != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ValidatorPubKeys) > 0 {
		for iNdEx := len(m.ValidatorPubKeys) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ValidatorPubKeys[iNdEx])
			copy(dAtA[i:], m.ValidatorPubKeys[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.ValidatorPubKeys[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *NoopIsm) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NoopIsm) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NoopIsm) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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
	return n
}

func (m *Ism) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.IsmType != 0 {
		n += 1 + sovTypes(uint64(m.IsmType))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Ism != nil {
		n += m.Ism.Size()
	}
	return n
}

func (m *Ism_MultiSig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MultiSig != nil {
		l = m.MultiSig.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *Ism_Noop) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Noop != nil {
		l = m.Noop.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *MultiSigIsm) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ValidatorPubKeys) > 0 {
		for _, s := range m.ValidatorPubKeys {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if m.Threshold != 0 {
		n += 1 + sovTypes(uint64(m.Threshold))
	}
	return n
}

func (m *NoopIsm) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
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
func (m *Ism) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Ism: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Ism: illegal tag %d (wire type %d)", fieldNum, wire)
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
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsmType", wireType)
			}
			m.IsmType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IsmType |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
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
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MultiSig", wireType)
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
			v := &MultiSigIsm{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Ism = &Ism_MultiSig{v}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Noop", wireType)
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
			v := &NoopIsm{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Ism = &Ism_Noop{v}
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
func (m *MultiSigIsm) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MultiSigIsm: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MultiSigIsm: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorPubKeys", wireType)
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
			m.ValidatorPubKeys = append(m.ValidatorPubKeys, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Threshold", wireType)
			}
			m.Threshold = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Threshold |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *NoopIsm) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NoopIsm: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NoopIsm: illegal tag %d (wire type %d)", fieldNum, wire)
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
