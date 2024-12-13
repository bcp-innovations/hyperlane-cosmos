// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/warp/v1/events.proto

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

type RemoteTransfer struct {
	DestinationDomain uint32 `protobuf:"varint,1,opt,name=destination_domain,json=destinationDomain,proto3" json:"destination_domain,omitempty"`
	RecipientAddress  string `protobuf:"bytes,2,opt,name=recipient_address,json=recipientAddress,proto3" json:"recipient_address,omitempty"`
}

func (m *RemoteTransfer) Reset()         { *m = RemoteTransfer{} }
func (m *RemoteTransfer) String() string { return proto.CompactTextString(m) }
func (*RemoteTransfer) ProtoMessage()    {}
func (*RemoteTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_0d53f48f6ba8625c, []int{0}
}
func (m *RemoteTransfer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RemoteTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RemoteTransfer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RemoteTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemoteTransfer.Merge(m, src)
}
func (m *RemoteTransfer) XXX_Size() int {
	return m.Size()
}
func (m *RemoteTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_RemoteTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_RemoteTransfer proto.InternalMessageInfo

func (m *RemoteTransfer) GetDestinationDomain() uint32 {
	if m != nil {
		return m.DestinationDomain
	}
	return 0
}

func (m *RemoteTransfer) GetRecipientAddress() string {
	if m != nil {
		return m.RecipientAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*RemoteTransfer)(nil), "hyperlane.warp.v1.RemoteTransfer")
}

func init() { proto.RegisterFile("hyperlane/warp/v1/events.proto", fileDescriptor_0d53f48f6ba8625c) }

var fileDescriptor_0d53f48f6ba8625c = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xcf, 0xb1, 0x4a, 0xc4, 0x40,
	0x10, 0xc6, 0xf1, 0xac, 0x85, 0x60, 0x40, 0x31, 0xa9, 0xae, 0x5a, 0x0e, 0xab, 0x03, 0x49, 0x96,
	0xc3, 0xc6, 0x56, 0xf1, 0x09, 0x82, 0x95, 0xcd, 0xb1, 0x49, 0x46, 0x6f, 0xe1, 0x76, 0x66, 0xd9,
	0x19, 0x57, 0xef, 0x2d, 0x7c, 0x2c, 0xcb, 0x2b, 0x2d, 0x25, 0x79, 0x11, 0x31, 0x07, 0xe1, 0xda,
	0xf9, 0x15, 0xf3, 0xfd, 0x73, 0xbd, 0xdd, 0x07, 0x88, 0x3b, 0x8b, 0x60, 0x3e, 0x6c, 0x0c, 0x26,
	0xad, 0x0d, 0x24, 0x40, 0xe1, 0x3a, 0x44, 0x12, 0x2a, 0x8b, 0xd9, 0xeb, 0x7f, 0xaf, 0xd3, 0xfa,
	0x66, 0x97, 0x5f, 0x35, 0xe0, 0x49, 0xe0, 0x39, 0x5a, 0xe4, 0x57, 0x88, 0x65, 0x95, 0x97, 0x3d,
	0xb0, 0x38, 0xb4, 0xe2, 0x08, 0x37, 0x3d, 0x79, 0xeb, 0x70, 0xa1, 0x96, 0x6a, 0x75, 0xd9, 0x14,
	0x27, 0xf2, 0x34, 0x41, 0x79, 0x9b, 0x17, 0x11, 0x3a, 0x17, 0x1c, 0xa0, 0x6c, 0x6c, 0xdf, 0x47,
	0x60, 0x5e, 0x9c, 0x2d, 0xd5, 0xea, 0xa2, 0xb9, 0x9e, 0xe1, 0xe1, 0x78, 0x7f, 0x6c, 0xbe, 0x07,
	0xad, 0x0e, 0x83, 0x56, 0xbf, 0x83, 0x56, 0x5f, 0xa3, 0xce, 0x0e, 0xa3, 0xce, 0x7e, 0x46, 0x9d,
	0xbd, 0xdc, 0xbf, 0x39, 0xd9, 0xbe, 0xb7, 0x75, 0x47, 0xde, 0xb4, 0x5d, 0xa8, 0x1c, 0x22, 0xa5,
	0xe9, 0x0f, 0x9b, 0x79, 0x75, 0xd5, 0x11, 0x7b, 0x62, 0xf3, 0x79, 0xcc, 0x93, 0x7d, 0x00, 0x6e,
	0xcf, 0xa7, 0xb6, 0xbb, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x98, 0x42, 0x18, 0xfd, 0x00,
	0x00, 0x00,
}

func (m *RemoteTransfer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RemoteTransfer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RemoteTransfer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.RecipientAddress) > 0 {
		i -= len(m.RecipientAddress)
		copy(dAtA[i:], m.RecipientAddress)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.RecipientAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.DestinationDomain != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.DestinationDomain))
		i--
		dAtA[i] = 0x8
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
func (m *RemoteTransfer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DestinationDomain != 0 {
		n += 1 + sovEvents(uint64(m.DestinationDomain))
	}
	l = len(m.RecipientAddress)
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
func (m *RemoteTransfer) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: RemoteTransfer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RemoteTransfer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationDomain", wireType)
			}
			m.DestinationDomain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DestinationDomain |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecipientAddress", wireType)
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
			m.RecipientAddress = string(dAtA[iNdEx:postIndex])
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
