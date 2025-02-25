// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/interchain_security/v1/types.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// MessageIdMultisigISM ...
type MessageIdMultisigISM struct {
	// XXX ...
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// owner ...
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	// validators
	// these are 20 byte long ethereum style addresses
	Validators []string `protobuf:"bytes,3,rep,name=validators,proto3" json:"validators,omitempty"`
	// threshold ...
	Threshold uint32 `protobuf:"varint,4,opt,name=threshold,proto3" json:"threshold,omitempty"`
}

func (m *MessageIdMultisigISM) Reset()         { *m = MessageIdMultisigISM{} }
func (m *MessageIdMultisigISM) String() string { return proto.CompactTextString(m) }
func (*MessageIdMultisigISM) ProtoMessage()    {}
func (*MessageIdMultisigISM) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9ae28ed3623cedf, []int{0}
}
func (m *MessageIdMultisigISM) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MessageIdMultisigISM) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MessageIdMultisigISM.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MessageIdMultisigISM) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageIdMultisigISM.Merge(m, src)
}
func (m *MessageIdMultisigISM) XXX_Size() int {
	return m.Size()
}
func (m *MessageIdMultisigISM) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageIdMultisigISM.DiscardUnknown(m)
}

var xxx_messageInfo_MessageIdMultisigISM proto.InternalMessageInfo

// MerkleRootMultiSigISM ...
type MerkleRootMultiSigISM struct {
	// XXX ...
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// owner ...
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	// validators
	// these are 20 byte long ethereum style addresses
	Validators []string `protobuf:"bytes,3,rep,name=validators,proto3" json:"validators,omitempty"`
	// threshold ...
	Threshold uint32 `protobuf:"varint,4,opt,name=threshold,proto3" json:"threshold,omitempty"`
}

func (m *MerkleRootMultiSigISM) Reset()         { *m = MerkleRootMultiSigISM{} }
func (m *MerkleRootMultiSigISM) String() string { return proto.CompactTextString(m) }
func (*MerkleRootMultiSigISM) ProtoMessage()    {}
func (*MerkleRootMultiSigISM) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9ae28ed3623cedf, []int{1}
}
func (m *MerkleRootMultiSigISM) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MerkleRootMultiSigISM) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MerkleRootMultiSigISM.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MerkleRootMultiSigISM) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MerkleRootMultiSigISM.Merge(m, src)
}
func (m *MerkleRootMultiSigISM) XXX_Size() int {
	return m.Size()
}
func (m *MerkleRootMultiSigISM) XXX_DiscardUnknown() {
	xxx_messageInfo_MerkleRootMultiSigISM.DiscardUnknown(m)
}

var xxx_messageInfo_MerkleRootMultiSigISM proto.InternalMessageInfo

// NoopISM ...
type NoopISM struct {
	// id ...
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// owner ...
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (m *NoopISM) Reset()         { *m = NoopISM{} }
func (m *NoopISM) String() string { return proto.CompactTextString(m) }
func (*NoopISM) ProtoMessage()    {}
func (*NoopISM) Descriptor() ([]byte, []int) {
	return fileDescriptor_b9ae28ed3623cedf, []int{2}
}
func (m *NoopISM) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *NoopISM) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_NoopISM.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *NoopISM) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NoopISM.Merge(m, src)
}
func (m *NoopISM) XXX_Size() int {
	return m.Size()
}
func (m *NoopISM) XXX_DiscardUnknown() {
	xxx_messageInfo_NoopISM.DiscardUnknown(m)
}

var xxx_messageInfo_NoopISM proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MessageIdMultisigISM)(nil), "hyperlane.core.interchain_security.v1.MessageIdMultisigISM")
	proto.RegisterType((*MerkleRootMultiSigISM)(nil), "hyperlane.core.interchain_security.v1.MerkleRootMultiSigISM")
	proto.RegisterType((*NoopISM)(nil), "hyperlane.core.interchain_security.v1.NoopISM")
}

func init() {
	proto.RegisterFile("hyperlane/core/interchain_security/v1/types.proto", fileDescriptor_b9ae28ed3623cedf)
}

var fileDescriptor_b9ae28ed3623cedf = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x93, 0x31, 0x8b, 0xdb, 0x30,
	0x14, 0xc7, 0xad, 0x24, 0x6d, 0x89, 0xa0, 0x1d, 0x4c, 0x0a, 0x6e, 0x28, 0xae, 0x09, 0x14, 0xbc,
	0xd8, 0xc2, 0x74, 0xeb, 0xd6, 0x4c, 0xc9, 0xe0, 0x52, 0xec, 0xad, 0x4b, 0x70, 0x2c, 0x61, 0x8b,
	0x3a, 0x7a, 0x46, 0x92, 0xdd, 0xe6, 0x1b, 0x74, 0xec, 0xde, 0xa5, 0x1f, 0x22, 0x1f, 0xa2, 0xdc,
	0x14, 0x6e, 0xba, 0xf1, 0x48, 0xe0, 0xee, 0x6b, 0x1c, 0x89, 0x72, 0xb9, 0xe3, 0xc8, 0x98, 0xe9,
	0x36, 0xe9, 0xff, 0xde, 0xff, 0xe9, 0xff, 0x03, 0x3d, 0x1c, 0x95, 0xcb, 0x9a, 0xc9, 0x2a, 0x13,
	0x8c, 0xe4, 0x20, 0x19, 0xe1, 0x42, 0x33, 0x99, 0x97, 0x19, 0x17, 0x33, 0xc5, 0xf2, 0x46, 0x72,
	0xbd, 0x24, 0x6d, 0x44, 0xf4, 0xb2, 0x66, 0x2a, 0xac, 0x25, 0x68, 0xb0, 0x3f, 0x1e, 0x2d, 0xe1,
	0xce, 0x12, 0x9e, 0xb0, 0x84, 0x6d, 0x34, 0x7c, 0x97, 0x83, 0x5a, 0x80, 0x9a, 0xed, 0x4d, 0xc4,
	0x5c, 0xcc, 0x84, 0xe1, 0xa0, 0x80, 0x02, 0x8c, 0xbe, 0x3b, 0x19, 0x75, 0x74, 0x83, 0xf0, 0x20,
	0x66, 0x4a, 0x65, 0x05, 0x9b, 0xd2, 0xb8, 0xa9, 0x34, 0x57, 0xbc, 0x98, 0xa6, 0xb1, 0xfd, 0x06,
	0x77, 0x38, 0x75, 0x90, 0x87, 0xfc, 0x5e, 0xd2, 0xe1, 0xd4, 0x0e, 0xf1, 0x0b, 0xf8, 0x29, 0x98,
	0x74, 0x3a, 0x1e, 0xf2, 0xfb, 0x63, 0xe7, 0x72, 0x15, 0x0c, 0x0e, 0xf3, 0xbf, 0x50, 0x2a, 0x99,
	0x52, 0xa9, 0x96, 0x5c, 0x14, 0x89, 0x69, 0xb3, 0x5d, 0x8c, 0xdb, 0xac, 0xe2, 0x34, 0xd3, 0x20,
	0x95, 0xd3, 0xf5, 0xba, 0x7e, 0x3f, 0x79, 0xa4, 0xd8, 0xef, 0x71, 0x5f, 0x97, 0x92, 0xa9, 0x12,
	0x2a, 0xea, 0xf4, 0x3c, 0xe4, 0xbf, 0x4e, 0x1e, 0x84, 0xcf, 0xdf, 0x7e, 0xff, 0xfb, 0x60, 0x5d,
	0xac, 0x82, 0xc9, 0x13, 0xec, 0x36, 0x0a, 0x67, 0xa7, 0xd0, 0x27, 0xf7, 0x6d, 0xd3, 0x63, 0x31,
	0x3d, 0xd4, 0x62, 0xa0, 0x4d, 0xc5, 0x46, 0xb7, 0x08, 0xbf, 0x8d, 0x99, 0xfc, 0x51, 0xb1, 0x04,
	0x40, 0xef, 0x49, 0xd3, 0xe7, 0x49, 0xfa, 0x17, 0xe1, 0x57, 0x5f, 0x01, 0xea, 0x33, 0xb0, 0x9d,
	0x3f, 0xdd, 0xb8, 0xf8, 0xbf, 0x71, 0xd1, 0x7a, 0xe3, 0xa2, 0xeb, 0x8d, 0x8b, 0xfe, 0x6c, 0x5d,
	0x6b, 0xbd, 0x75, 0xad, 0xab, 0xad, 0x6b, 0x7d, 0x8f, 0x0b, 0xae, 0xcb, 0x66, 0x1e, 0xe6, 0xb0,
	0x20, 0xf3, 0xbc, 0x0e, 0xb8, 0x10, 0xd0, 0x66, 0x9a, 0x83, 0x50, 0xe4, 0xf8, 0x7c, 0x60, 0xd2,
	0x92, 0x5f, 0x66, 0x73, 0x4e, 0x85, 0x30, 0x7b, 0x33, 0x7f, 0xb9, 0xff, 0xe0, 0x9f, 0xee, 0x02,
	0x00, 0x00, 0xff, 0xff, 0xf2, 0xbd, 0xf7, 0x7a, 0x6d, 0x03, 0x00, 0x00,
}

func (m *MessageIdMultisigISM) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MessageIdMultisigISM) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MessageIdMultisigISM) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Threshold != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Validators) > 0 {
		for iNdEx := len(m.Validators) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Validators[iNdEx])
			copy(dAtA[i:], m.Validators[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.Validators[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MerkleRootMultiSigISM) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MerkleRootMultiSigISM) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MerkleRootMultiSigISM) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Threshold != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Threshold))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Validators) > 0 {
		for iNdEx := len(m.Validators) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Validators[iNdEx])
			copy(dAtA[i:], m.Validators[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.Validators[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *NoopISM) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *NoopISM) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *NoopISM) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
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
func (m *MessageIdMultisigISM) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.Validators) > 0 {
		for _, s := range m.Validators {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if m.Threshold != 0 {
		n += 1 + sovTypes(uint64(m.Threshold))
	}
	return n
}

func (m *MerkleRootMultiSigISM) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.Validators) > 0 {
		for _, s := range m.Validators {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if m.Threshold != 0 {
		n += 1 + sovTypes(uint64(m.Threshold))
	}
	return n
}

func (m *NoopISM) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MessageIdMultisigISM) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MessageIdMultisigISM: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MessageIdMultisigISM: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validators", wireType)
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
			m.Validators = append(m.Validators, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
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
func (m *MerkleRootMultiSigISM) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MerkleRootMultiSigISM: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MerkleRootMultiSigISM: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validators", wireType)
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
			m.Validators = append(m.Validators, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
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
func (m *NoopISM) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: NoopISM: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: NoopISM: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
