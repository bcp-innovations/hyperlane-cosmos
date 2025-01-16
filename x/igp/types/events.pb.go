// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/igp/v1/events.proto

package types

import (
	fmt "fmt"
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

type GasPayment struct {
	MessageId   string `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	Destination uint32 `protobuf:"varint,2,opt,name=destination,proto3" json:"destination,omitempty"`
	GasAmount   uint64 `protobuf:"varint,3,opt,name=gas_amount,json=gasAmount,proto3" json:"gas_amount,omitempty"`
	Payment     uint64 `protobuf:"varint,4,opt,name=payment,proto3" json:"payment,omitempty"`
	IgpId       string `protobuf:"bytes,5,opt,name=igp_id,json=igpId,proto3" json:"igp_id,omitempty"`
}

func (m *GasPayment) Reset()         { *m = GasPayment{} }
func (m *GasPayment) String() string { return proto.CompactTextString(m) }
func (*GasPayment) ProtoMessage()    {}
func (*GasPayment) Descriptor() ([]byte, []int) {
	return fileDescriptor_5c3e8b8d99d4a845, []int{0}
}
func (m *GasPayment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GasPayment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GasPayment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GasPayment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GasPayment.Merge(m, src)
}
func (m *GasPayment) XXX_Size() int {
	return m.Size()
}
func (m *GasPayment) XXX_DiscardUnknown() {
	xxx_messageInfo_GasPayment.DiscardUnknown(m)
}

var xxx_messageInfo_GasPayment proto.InternalMessageInfo

func (m *GasPayment) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *GasPayment) GetDestination() uint32 {
	if m != nil {
		return m.Destination
	}
	return 0
}

func (m *GasPayment) GetGasAmount() uint64 {
	if m != nil {
		return m.GasAmount
	}
	return 0
}

func (m *GasPayment) GetPayment() uint64 {
	if m != nil {
		return m.Payment
	}
	return 0
}

func (m *GasPayment) GetIgpId() string {
	if m != nil {
		return m.IgpId
	}
	return ""
}

func init() {
	proto.RegisterType((*GasPayment)(nil), "hyperlane.igp.v1.GasPayment")
}

func init() { proto.RegisterFile("hyperlane/igp/v1/events.proto", fileDescriptor_5c3e8b8d99d4a845) }

var fileDescriptor_5c3e8b8d99d4a845 = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x86, 0x63, 0x68, 0x8b, 0x62, 0x84, 0x84, 0x22, 0x21, 0x79, 0xa9, 0x15, 0x31, 0x65, 0x69,
	0xac, 0x8a, 0x81, 0x19, 0x16, 0xd4, 0x0d, 0x32, 0xb2, 0x54, 0x4e, 0x62, 0xb9, 0x96, 0xb0, 0x7d,
	0xea, 0xb9, 0x11, 0x79, 0x0b, 0x5e, 0x80, 0xf7, 0x61, 0xec, 0xc8, 0x88, 0x92, 0x17, 0x41, 0x35,
	0x50, 0x31, 0xde, 0xf7, 0x4b, 0xf7, 0xdd, 0xfd, 0x74, 0xbe, 0xe9, 0x41, 0x6d, 0x5f, 0xa4, 0x53,
	0xc2, 0x68, 0x10, 0xdd, 0x52, 0xa8, 0x4e, 0xb9, 0x80, 0x25, 0x6c, 0x7d, 0xf0, 0xd9, 0xe5, 0x31,
	0x2e, 0x8d, 0x86, 0xb2, 0x5b, 0x5e, 0xbf, 0x13, 0x4a, 0x1f, 0x24, 0x3e, 0xca, 0xde, 0x2a, 0x17,
	0xb2, 0x39, 0xa5, 0x56, 0x21, 0x4a, 0xad, 0xd6, 0xa6, 0x65, 0x24, 0x27, 0x45, 0x5a, 0xa5, 0xbf,
	0x64, 0xd5, 0x66, 0x39, 0x3d, 0x6f, 0x15, 0x06, 0xe3, 0x64, 0x30, 0xde, 0xb1, 0x93, 0x9c, 0x14,
	0x17, 0xd5, 0x7f, 0x74, 0x58, 0xa0, 0x25, 0xae, 0xa5, 0xf5, 0x3b, 0x17, 0xd8, 0x69, 0x4e, 0x8a,
	0x49, 0x95, 0x6a, 0x89, 0x77, 0x11, 0x64, 0x8c, 0x9e, 0xc1, 0x8f, 0x8a, 0x4d, 0x62, 0xf6, 0x37,
	0x66, 0x57, 0x74, 0x66, 0x34, 0x1c, 0xac, 0xd3, 0x68, 0x9d, 0x1a, 0x0d, 0xab, 0xf6, 0xfe, 0xe9,
	0x63, 0xe0, 0x64, 0x3f, 0x70, 0xf2, 0x35, 0x70, 0xf2, 0x36, 0xf2, 0x64, 0x3f, 0xf2, 0xe4, 0x73,
	0xe4, 0xc9, 0xf3, 0xad, 0x36, 0x61, 0xb3, 0xab, 0xcb, 0xc6, 0x5b, 0x51, 0x37, 0xb0, 0x30, 0xce,
	0xf9, 0x2e, 0x1e, 0x81, 0xe2, 0xf8, 0xe6, 0xa2, 0xf1, 0x68, 0x3d, 0x8a, 0xd7, 0x58, 0x47, 0xe8,
	0x41, 0x61, 0x3d, 0x8b, 0x5d, 0xdc, 0x7c, 0x07, 0x00, 0x00, 0xff, 0xff, 0x32, 0x5e, 0x5a, 0xbf,
	0x2c, 0x01, 0x00, 0x00,
}

func (m *GasPayment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GasPayment) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GasPayment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.IgpId) > 0 {
		i -= len(m.IgpId)
		copy(dAtA[i:], m.IgpId)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.IgpId)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Payment != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Payment))
		i--
		dAtA[i] = 0x20
	}
	if m.GasAmount != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.GasAmount))
		i--
		dAtA[i] = 0x18
	}
	if m.Destination != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Destination))
		i--
		dAtA[i] = 0x10
	}
	if len(m.MessageId) > 0 {
		i -= len(m.MessageId)
		copy(dAtA[i:], m.MessageId)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.MessageId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GasPayment) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.MessageId)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Destination != 0 {
		n += 1 + sovEvents(uint64(m.Destination))
	}
	if m.GasAmount != 0 {
		n += 1 + sovEvents(uint64(m.GasAmount))
	}
	if m.Payment != 0 {
		n += 1 + sovEvents(uint64(m.Payment))
	}
	l = len(m.IgpId)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GasPayment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: GasPayment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GasPayment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Destination", wireType)
			}
			m.Destination = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Destination |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasAmount", wireType)
			}
			m.GasAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Payment", wireType)
			}
			m.Payment = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Payment |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IgpId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IgpId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
