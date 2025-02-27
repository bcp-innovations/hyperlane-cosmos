// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/v1/genesis.proto

package types

import (
	fmt "fmt"
	types "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	types1 "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
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

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	// params defines all the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// ism_genesis
	IsmGenesis *types.GenesisState `protobuf:"bytes,2,opt,name=ism_genesis,json=ismGenesis,proto3" json:"ism_genesis,omitempty"`
	// post_dispatch_genesis
	PostDispatchGenesis *types1.GenesisState `protobuf:"bytes,3,opt,name=post_dispatch_genesis,json=postDispatchGenesis,proto3" json:"post_dispatch_genesis,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_9329350a78ea2d1f, []int{0}
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

func (m *GenesisState) GetIsmGenesis() *types.GenesisState {
	if m != nil {
		return m.IsmGenesis
	}
	return nil
}

func (m *GenesisState) GetPostDispatchGenesis() *types1.GenesisState {
	if m != nil {
		return m.PostDispatchGenesis
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "hyperlane.core.v1.GenesisState")
}

func init() { proto.RegisterFile("hyperlane/core/v1/genesis.proto", fileDescriptor_9329350a78ea2d1f) }

var fileDescriptor_9329350a78ea2d1f = []byte{
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xb1, 0x4a, 0x3b, 0x41,
	0x10, 0xc6, 0xef, 0xf2, 0x87, 0xc0, 0xff, 0x62, 0x93, 0x53, 0x21, 0x1e, 0x78, 0x11, 0x2b, 0x11,
	0x6e, 0x97, 0x98, 0xc6, 0xc2, 0x2a, 0x08, 0xb6, 0x12, 0xad, 0x6c, 0xc2, 0x66, 0x5d, 0xee, 0x16,
	0xdc, 0x9d, 0x65, 0x77, 0x13, 0xcc, 0x03, 0xd8, 0xfb, 0x18, 0x96, 0x3e, 0x46, 0xca, 0x94, 0x56,
	0x22, 0x49, 0xe1, 0x6b, 0xc8, 0x6d, 0x96, 0xc3, 0xdc, 0xd9, 0x1c, 0xc7, 0x7c, 0xdf, 0xfc, 0x66,
	0xbe, 0xd9, 0xa8, 0x5f, 0x2c, 0x14, 0xd3, 0x4f, 0x44, 0x32, 0x4c, 0x41, 0x33, 0x3c, 0x1f, 0xe0,
	0x9c, 0x49, 0x66, 0xb8, 0x41, 0x4a, 0x83, 0x85, 0xb8, 0x5b, 0x19, 0x50, 0x69, 0x40, 0xf3, 0x41,
	0x72, 0x90, 0x43, 0x0e, 0x4e, 0xc5, 0xe5, 0xdf, 0xd6, 0x98, 0x74, 0x89, 0xe0, 0x12, 0xb0, 0xfb,
	0xfa, 0xd2, 0x71, 0x13, 0x6e, 0x17, 0x8a, 0x79, 0x74, 0x32, 0xac, 0xc9, 0x5c, 0x5a, 0xa6, 0x69,
	0x41, 0xb8, 0x9c, 0x18, 0x46, 0x67, 0x9a, 0xdb, 0x45, 0x63, 0x9f, 0x24, 0xab, 0x35, 0x29, 0x30,
	0x76, 0xf2, 0xc8, 0x8d, 0x22, 0x96, 0x16, 0x0d, 0xfb, 0xe9, 0x4b, 0x2b, 0xda, 0xbb, 0xd9, 0x56,
	0xee, 0x2c, 0xb1, 0x2c, 0xbe, 0x8a, 0xda, 0x8a, 0x68, 0x22, 0x4c, 0x2f, 0x3c, 0x09, 0xcf, 0x3a,
	0x17, 0x47, 0xa8, 0x11, 0x10, 0xdd, 0x3a, 0xc3, 0xe8, 0xff, 0xf2, 0xb3, 0x1f, 0xbc, 0x7d, 0xbf,
	0x9f, 0x87, 0x63, 0xdf, 0x13, 0xdf, 0x47, 0x1d, 0x6e, 0xc4, 0xc4, 0xcf, 0xe8, 0xb5, 0x1c, 0x62,
	0x58, 0x47, 0xfc, 0x11, 0xa4, 0xc4, 0xfe, 0xde, 0x63, 0x1c, 0x71, 0x23, 0x7c, 0x21, 0x26, 0xd1,
	0xe1, 0x4e, 0x8c, 0x8a, 0xff, 0xcf, 0xf1, 0xb3, 0x3a, 0x7f, 0xc7, 0xdc, 0x20, 0xef, 0x97, 0xf2,
	0xb5, 0x57, 0xbd, 0x32, 0x1a, 0x2f, 0xd7, 0x69, 0xb8, 0x5a, 0xa7, 0xe1, 0xd7, 0x3a, 0x0d, 0x5f,
	0x37, 0x69, 0xb0, 0xda, 0xa4, 0xc1, 0xc7, 0x26, 0x0d, 0x1e, 0x2e, 0x73, 0x6e, 0x8b, 0xd9, 0x14,
	0x51, 0x10, 0x78, 0x4a, 0x55, 0xc6, 0xa5, 0x84, 0x39, 0xb1, 0x1c, 0xa4, 0xc1, 0xd5, 0xdc, 0x8c,
	0x82, 0x11, 0x60, 0xf0, 0xf3, 0xf6, 0xe8, 0xee, 0x15, 0xa7, 0x6d, 0x77, 0xe2, 0xe1, 0x4f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x24, 0x56, 0x72, 0x07, 0x44, 0x02, 0x00, 0x00,
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
	if m.PostDispatchGenesis != nil {
		{
			size, err := m.PostDispatchGenesis.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.IsmGenesis != nil {
		{
			size, err := m.IsmGenesis.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenesis(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.IsmGenesis != nil {
		l = m.IsmGenesis.Size()
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.PostDispatchGenesis != nil {
		l = m.PostDispatchGenesis.Size()
		n += 1 + l + sovGenesis(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
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
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsmGenesis", wireType)
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
			if m.IsmGenesis == nil {
				m.IsmGenesis = &types.GenesisState{}
			}
			if err := m.IsmGenesis.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostDispatchGenesis", wireType)
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
			if m.PostDispatchGenesis == nil {
				m.PostDispatchGenesis = &types1.GenesisState{}
			}
			if err := m.PostDispatchGenesis.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
