// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/post_dispatch/v1/genesis.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
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

// GenesisState defines the post dispatch submodule's genesis state.
type GenesisState struct {
	Igps            []InterchainGasPaymaster             `protobuf:"bytes,1,rep,name=igps,proto3" json:"igps"`
	IgpGasConfigs   []GenesisDestinationGasConfigWrapper `protobuf:"bytes,2,rep,name=igp_gas_configs,json=igpGasConfigs,proto3" json:"igp_gas_configs"`
	MerkleTreeHooks []MerkleTreeHook                     `protobuf:"bytes,3,rep,name=merkle_tree_hooks,json=merkleTreeHooks,proto3" json:"merkle_tree_hooks"`
	NoopHooks       []NoopHook                           `protobuf:"bytes,4,rep,name=noop_hooks,json=noopHooks,proto3" json:"noop_hooks"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_8864b1a76aa43cd2, []int{0}
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

func (m *GenesisState) GetIgps() []InterchainGasPaymaster {
	if m != nil {
		return m.Igps
	}
	return nil
}

func (m *GenesisState) GetIgpGasConfigs() []GenesisDestinationGasConfigWrapper {
	if m != nil {
		return m.IgpGasConfigs
	}
	return nil
}

func (m *GenesisState) GetMerkleTreeHooks() []MerkleTreeHook {
	if m != nil {
		return m.MerkleTreeHooks
	}
	return nil
}

func (m *GenesisState) GetNoopHooks() []NoopHook {
	if m != nil {
		return m.NoopHooks
	}
	return nil
}

// GenesisDestinationGasConfigWrapper ...
type GenesisDestinationGasConfigWrapper struct {
	// remote_domain ...
	RemoteDomain uint32 `protobuf:"varint,1,opt,name=remote_domain,json=remoteDomain,proto3" json:"remote_domain,omitempty"`
	// gas_oracle ...
	GasOracle *GasOracle `protobuf:"bytes,2,opt,name=gas_oracle,json=gasOracle,proto3" json:"gas_oracle,omitempty"`
	// gas_overhead ...
	GasOverhead cosmossdk_io_math.Int `protobuf:"bytes,3,opt,name=gas_overhead,json=gasOverhead,proto3,customtype=cosmossdk.io/math.Int" json:"gas_overhead"`
	// igp_id is required for the Genesis handling.
	IgpId uint64 `protobuf:"varint,4,opt,name=igp_id,json=igpId,proto3" json:"igp_id,omitempty"`
}

func (m *GenesisDestinationGasConfigWrapper) Reset()         { *m = GenesisDestinationGasConfigWrapper{} }
func (m *GenesisDestinationGasConfigWrapper) String() string { return proto.CompactTextString(m) }
func (*GenesisDestinationGasConfigWrapper) ProtoMessage()    {}
func (*GenesisDestinationGasConfigWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_8864b1a76aa43cd2, []int{1}
}
func (m *GenesisDestinationGasConfigWrapper) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisDestinationGasConfigWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisDestinationGasConfigWrapper.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisDestinationGasConfigWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisDestinationGasConfigWrapper.Merge(m, src)
}
func (m *GenesisDestinationGasConfigWrapper) XXX_Size() int {
	return m.Size()
}
func (m *GenesisDestinationGasConfigWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisDestinationGasConfigWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisDestinationGasConfigWrapper proto.InternalMessageInfo

func (m *GenesisDestinationGasConfigWrapper) GetRemoteDomain() uint32 {
	if m != nil {
		return m.RemoteDomain
	}
	return 0
}

func (m *GenesisDestinationGasConfigWrapper) GetGasOracle() *GasOracle {
	if m != nil {
		return m.GasOracle
	}
	return nil
}

func (m *GenesisDestinationGasConfigWrapper) GetIgpId() uint64 {
	if m != nil {
		return m.IgpId
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "hyperlane.core.post_dispatch.v1.GenesisState")
	proto.RegisterType((*GenesisDestinationGasConfigWrapper)(nil), "hyperlane.core.post_dispatch.v1.GenesisDestinationGasConfigWrapper")
}

func init() {
	proto.RegisterFile("hyperlane/core/post_dispatch/v1/genesis.proto", fileDescriptor_8864b1a76aa43cd2)
}

var fileDescriptor_8864b1a76aa43cd2 = []byte{
	// 484 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xe3, 0x26, 0x54, 0xca, 0x26, 0x51, 0x85, 0x45, 0x25, 0xab, 0x12, 0x4e, 0x14, 0x2e,
	0x01, 0x14, 0x2f, 0x0d, 0x07, 0xae, 0x28, 0xad, 0x54, 0x72, 0xa0, 0x40, 0x40, 0x42, 0xe2, 0x62,
	0x6d, 0xec, 0x61, 0xbd, 0x4a, 0xbc, 0xb3, 0xec, 0x2e, 0x11, 0x79, 0x02, 0xae, 0x3c, 0x56, 0x8f,
	0x3d, 0x22, 0x0e, 0x15, 0x4a, 0xce, 0xbc, 0x03, 0xf2, 0xda, 0x44, 0xe4, 0xe4, 0xdb, 0x78, 0x3c,
	0xf3, 0xfd, 0x33, 0xb3, 0x3f, 0x19, 0x67, 0x1b, 0x05, 0x7a, 0xc5, 0x24, 0xd0, 0x04, 0x35, 0x50,
	0x85, 0xc6, 0xc6, 0xa9, 0x30, 0x8a, 0xd9, 0x24, 0xa3, 0xeb, 0x73, 0xca, 0x41, 0x82, 0x11, 0x26,
	0x52, 0x1a, 0x2d, 0xfa, 0xfd, 0x7d, 0x79, 0x54, 0x94, 0x47, 0x07, 0xe5, 0xd1, 0xfa, 0xfc, 0xec,
	0x69, 0x1d, 0xcf, 0x6e, 0x14, 0x54, 0xb4, 0xb3, 0x07, 0x1c, 0x39, 0xba, 0x90, 0x16, 0x51, 0x99,
	0x1d, 0x7e, 0x6f, 0x92, 0xee, 0x55, 0xa9, 0xfa, 0xde, 0x32, 0x0b, 0xfe, 0x3b, 0xd2, 0x12, 0x5c,
	0x99, 0xc0, 0x1b, 0x34, 0x47, 0x9d, 0xc9, 0x8b, 0xa8, 0x66, 0x86, 0x68, 0x26, 0x2d, 0xe8, 0x24,
	0x63, 0x42, 0x5e, 0x31, 0xf3, 0x96, 0x6d, 0x72, 0x66, 0x2c, 0xe8, 0x69, 0xeb, 0xe6, 0xae, 0xdf,
	0x98, 0x3b, 0x94, 0xff, 0x85, 0x9c, 0x08, 0xae, 0x62, 0xce, 0x4c, 0x9c, 0xa0, 0xfc, 0x2c, 0xb8,
	0x09, 0x8e, 0x1c, 0xfd, 0xa2, 0x96, 0x5e, 0x8d, 0x76, 0x09, 0xc6, 0x0a, 0xc9, 0xac, 0xc0, 0x42,
	0xe5, 0xc2, 0x41, 0x3e, 0x6a, 0xa6, 0xd4, 0x5e, 0xa9, 0x27, 0xb8, 0xda, 0xff, 0x32, 0x3e, 0x23,
	0xf7, 0x73, 0xd0, 0xcb, 0x15, 0xc4, 0x56, 0x03, 0xc4, 0x19, 0xe2, 0xd2, 0x04, 0x4d, 0x27, 0x4a,
	0x6b, 0x45, 0x5f, 0xbb, 0xce, 0x0f, 0x1a, 0xe0, 0x15, 0xe2, 0xb2, 0x12, 0x38, 0xc9, 0x0f, 0xb2,
	0xc6, 0xbf, 0x26, 0x44, 0x22, 0xaa, 0x8a, 0xdd, 0x72, 0xec, 0xc7, 0xb5, 0xec, 0x6b, 0x44, 0xf5,
	0x1f, 0xb5, 0x2d, 0xab, 0x6f, 0x33, 0xfc, 0xe3, 0x91, 0x61, 0xfd, 0xba, 0xfe, 0x23, 0xd2, 0xd3,
	0x90, 0xa3, 0x85, 0x38, 0xc5, 0x9c, 0x09, 0x19, 0x78, 0x03, 0x6f, 0xd4, 0x9b, 0x77, 0xcb, 0xe4,
	0xa5, 0xcb, 0xf9, 0x33, 0x42, 0x8a, 0x6b, 0xa3, 0x66, 0xc9, 0x0a, 0x82, 0xa3, 0x81, 0x37, 0xea,
	0x4c, 0x9e, 0xd4, 0x1f, 0x9b, 0x99, 0x37, 0xae, 0x63, 0xde, 0xe6, 0xff, 0x42, 0xff, 0x25, 0xe9,
	0x3a, 0xd4, 0x1a, 0x74, 0x06, 0x2c, 0x0d, 0x9a, 0x03, 0x6f, 0xd4, 0x9e, 0x3e, 0x2c, 0xa6, 0xff,
	0x75, 0xd7, 0x3f, 0x4d, 0xd0, 0xe4, 0x68, 0x4c, 0xba, 0x8c, 0x04, 0xd2, 0x9c, 0xd9, 0xac, 0x30,
	0xc4, 0xbc, 0x53, 0xf4, 0x57, 0x1d, 0xfe, 0x29, 0x39, 0x2e, 0x9e, 0x5f, 0xa4, 0x41, 0x6b, 0xe0,
	0x8d, 0x5a, 0xf3, 0x7b, 0x82, 0xab, 0x59, 0x3a, 0x4d, 0x6e, 0xb6, 0xa1, 0x77, 0xbb, 0x0d, 0xbd,
	0xdf, 0xdb, 0xd0, 0xfb, 0xb1, 0x0b, 0x1b, 0xb7, 0xbb, 0xb0, 0xf1, 0x73, 0x17, 0x36, 0x3e, 0xcd,
	0xb8, 0xb0, 0xd9, 0xd7, 0x45, 0x94, 0x60, 0x4e, 0x17, 0x89, 0x1a, 0x0b, 0x29, 0x71, 0xed, 0xae,
	0x61, 0xe8, 0x7e, 0x87, 0x71, 0xa9, 0x4c, 0xbf, 0x95, 0xd6, 0x7f, 0x36, 0x89, 0x0f, 0xdd, 0xef,
	0xac, 0xbf, 0x38, 0x76, 0x2e, 0x7f, 0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0xd3, 0xa8, 0x32, 0xde,
	0x7a, 0x03, 0x00, 0x00,
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
	if len(m.NoopHooks) > 0 {
		for iNdEx := len(m.NoopHooks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.NoopHooks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.MerkleTreeHooks) > 0 {
		for iNdEx := len(m.MerkleTreeHooks) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MerkleTreeHooks[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.IgpGasConfigs) > 0 {
		for iNdEx := len(m.IgpGasConfigs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.IgpGasConfigs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Igps) > 0 {
		for iNdEx := len(m.Igps) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Igps[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *GenesisDestinationGasConfigWrapper) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisDestinationGasConfigWrapper) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisDestinationGasConfigWrapper) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.IgpId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.IgpId))
		i--
		dAtA[i] = 0x20
	}
	{
		size := m.GasOverhead.Size()
		i -= size
		if _, err := m.GasOverhead.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.GasOracle != nil {
		{
			size, err := m.GasOracle.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.RemoteDomain != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.RemoteDomain))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Igps) > 0 {
		for _, e := range m.Igps {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.IgpGasConfigs) > 0 {
		for _, e := range m.IgpGasConfigs {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.MerkleTreeHooks) > 0 {
		for _, e := range m.MerkleTreeHooks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.NoopHooks) > 0 {
		for _, e := range m.NoopHooks {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GenesisDestinationGasConfigWrapper) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RemoteDomain != 0 {
		n += 1 + sovGenesis(uint64(m.RemoteDomain))
	}
	if m.GasOracle != nil {
		l = m.GasOracle.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = m.GasOverhead.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.IgpId != 0 {
		n += 1 + sovGenesis(uint64(m.IgpId))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
				return fmt.Errorf("proto: wrong wireType = %d for field Igps", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Igps = append(m.Igps, InterchainGasPaymaster{})
			if err := m.Igps[len(m.Igps)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IgpGasConfigs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IgpGasConfigs = append(m.IgpGasConfigs, GenesisDestinationGasConfigWrapper{})
			if err := m.IgpGasConfigs[len(m.IgpGasConfigs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MerkleTreeHooks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MerkleTreeHooks = append(m.MerkleTreeHooks, MerkleTreeHook{})
			if err := m.MerkleTreeHooks[len(m.MerkleTreeHooks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NoopHooks", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NoopHooks = append(m.NoopHooks, NoopHook{})
			if err := m.NoopHooks[len(m.NoopHooks)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *GenesisDestinationGasConfigWrapper) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisDestinationGasConfigWrapper: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisDestinationGasConfigWrapper: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RemoteDomain", wireType)
			}
			m.RemoteDomain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RemoteDomain |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasOracle", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.GasOracle == nil {
				m.GasOracle = &GasOracle{}
			}
			if err := m.GasOracle.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasOverhead", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.GasOverhead.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IgpId", wireType)
			}
			m.IgpId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IgpId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
