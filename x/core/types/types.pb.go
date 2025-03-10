// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/v1/types.proto

package types

import (
	fmt "fmt"
	github_com_bcp_innovations_hyperlane_cosmos_util "github.com/bcp-innovations/hyperlane-cosmos/util"
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

// Mailbox ...
type Mailbox struct {
	Id github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress `protobuf:"bytes,1,opt,name=id,proto3,customtype=github.com/bcp-innovations/hyperlane-cosmos/util.HexAddress" json:"id"`
	// owner ...
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	// message_sent ...
	MessageSent uint32 `protobuf:"varint,3,opt,name=message_sent,json=messageSent,proto3" json:"message_sent,omitempty"`
	// message_received ...
	MessageReceived uint32 `protobuf:"varint,4,opt,name=message_received,json=messageReceived,proto3" json:"message_received,omitempty"`
	// default_ism ...
	DefaultIsm github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress `protobuf:"bytes,5,opt,name=default_ism,json=defaultIsm,proto3,customtype=github.com/bcp-innovations/hyperlane-cosmos/util.HexAddress" json:"default_ism"`
	// default_hook
	DefaultHook *github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress `protobuf:"bytes,8,opt,name=default_hook,json=defaultHook,proto3,customtype=github.com/bcp-innovations/hyperlane-cosmos/util.HexAddress" json:"default_hook,omitempty"`
	// required_hook
	RequiredHook *github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress `protobuf:"bytes,9,opt,name=required_hook,json=requiredHook,proto3,customtype=github.com/bcp-innovations/hyperlane-cosmos/util.HexAddress" json:"required_hook,omitempty"`
	// domain
	LocalDomain uint32 `protobuf:"varint,10,opt,name=local_domain,json=localDomain,proto3" json:"local_domain,omitempty"`
}

func (m *Mailbox) Reset()         { *m = Mailbox{} }
func (m *Mailbox) String() string { return proto.CompactTextString(m) }
func (*Mailbox) ProtoMessage()    {}
func (*Mailbox) Descriptor() ([]byte, []int) {
	return fileDescriptor_d14de0fc8fa7fd67, []int{0}
}
func (m *Mailbox) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Mailbox) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Mailbox.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Mailbox) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Mailbox.Merge(m, src)
}
func (m *Mailbox) XXX_Size() int {
	return m.Size()
}
func (m *Mailbox) XXX_DiscardUnknown() {
	xxx_messageInfo_Mailbox.DiscardUnknown(m)
}

var xxx_messageInfo_Mailbox proto.InternalMessageInfo

func (m *Mailbox) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *Mailbox) GetMessageSent() uint32 {
	if m != nil {
		return m.MessageSent
	}
	return 0
}

func (m *Mailbox) GetMessageReceived() uint32 {
	if m != nil {
		return m.MessageReceived
	}
	return 0
}

func (m *Mailbox) GetLocalDomain() uint32 {
	if m != nil {
		return m.LocalDomain
	}
	return 0
}

func init() {
	proto.RegisterType((*Mailbox)(nil), "hyperlane.core.v1.Mailbox")
}

func init() { proto.RegisterFile("hyperlane/core/v1/types.proto", fileDescriptor_d14de0fc8fa7fd67) }

var fileDescriptor_d14de0fc8fa7fd67 = []byte{
	// 389 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0x31, 0x6f, 0xd4, 0x30,
	0x14, 0xc7, 0x93, 0xa3, 0x05, 0xea, 0xa6, 0x02, 0xac, 0x0e, 0xa6, 0x12, 0x69, 0xcb, 0x54, 0x86,
	0x4b, 0x54, 0xb1, 0x20, 0x31, 0x51, 0x18, 0xca, 0xc0, 0x92, 0xdb, 0x58, 0xa2, 0x24, 0x7e, 0x4d,
	0xac, 0x26, 0x7e, 0xc1, 0x76, 0xc2, 0xf5, 0x5b, 0xf0, 0x61, 0x90, 0xf8, 0x0a, 0x1d, 0x4f, 0x4c,
	0x88, 0xe1, 0x84, 0xee, 0xbe, 0x08, 0x3a, 0x3b, 0xb9, 0x1d, 0x71, 0x9b, 0xfd, 0x7b, 0x7f, 0xbd,
	0xdf, 0xb3, 0xfc, 0xc8, 0x8b, 0xea, 0xae, 0x05, 0x55, 0x67, 0x12, 0xe2, 0x02, 0x15, 0xc4, 0xfd,
	0x65, 0x6c, 0xee, 0x5a, 0xd0, 0x51, 0xab, 0xd0, 0x20, 0x7d, 0xb6, 0x2d, 0x47, 0x9b, 0x72, 0xd4,
	0x5f, 0x9e, 0x3c, 0x2f, 0x50, 0x37, 0xa8, 0x53, 0x1b, 0x88, 0xdd, 0xc5, 0xa5, 0x4f, 0x8e, 0x4b,
	0x2c, 0xd1, 0xf1, 0xcd, 0xc9, 0xd1, 0x97, 0x3f, 0xf6, 0xc8, 0xa3, 0x4f, 0x99, 0xa8, 0x73, 0x9c,
	0xd3, 0x19, 0x99, 0x08, 0xce, 0xfc, 0x33, 0xff, 0xe2, 0xe0, 0xea, 0xfd, 0xfd, 0xf2, 0xd4, 0xfb,
	0xbd, 0x3c, 0x7d, 0x5b, 0x0a, 0x53, 0x75, 0x79, 0x54, 0x60, 0x13, 0xe7, 0x45, 0x3b, 0x15, 0x52,
	0x62, 0x9f, 0x19, 0x81, 0x52, 0xc7, 0x5b, 0xfd, 0xd4, 0x89, 0xe2, 0xce, 0x88, 0x3a, 0xba, 0x86,
	0xf9, 0x3b, 0xce, 0x15, 0x68, 0x9d, 0x4c, 0x04, 0xa7, 0x11, 0xd9, 0xc7, 0xaf, 0x12, 0x14, 0x9b,
	0xd8, 0xbe, 0xec, 0xe7, 0xf7, 0xe9, 0xf1, 0x30, 0xd7, 0x10, 0x9b, 0x19, 0x25, 0x64, 0x99, 0xb8,
	0x18, 0x3d, 0x27, 0x41, 0x03, 0x5a, 0x67, 0x25, 0xa4, 0x1a, 0xa4, 0x61, 0x0f, 0xce, 0xfc, 0x8b,
	0xa3, 0xe4, 0x70, 0x60, 0x33, 0x90, 0x86, 0xbe, 0x22, 0x4f, 0xc7, 0x88, 0x82, 0x02, 0x44, 0x0f,
	0x9c, 0xed, 0xd9, 0xd8, 0x93, 0x81, 0x27, 0x03, 0xa6, 0x9c, 0x1c, 0x72, 0xb8, 0xc9, 0xba, 0xda,
	0xa4, 0x42, 0x37, 0x6c, 0x7f, 0x77, 0x6f, 0x23, 0x43, 0xdf, 0x8f, 0xba, 0xa1, 0x37, 0x24, 0x18,
	0x2d, 0x15, 0xe2, 0x2d, 0x7b, 0xbc, 0xd5, 0xf8, 0xff, 0xab, 0x19, 0xc7, 0xbf, 0x46, 0xbc, 0xa5,
	0x15, 0x39, 0x52, 0xf0, 0xa5, 0x13, 0x0a, 0xb8, 0x13, 0x1d, 0xec, 0x4e, 0x14, 0x8c, 0x9d, 0xad,
	0xe9, 0x9c, 0x04, 0x35, 0x16, 0x59, 0x9d, 0x72, 0x6c, 0x32, 0x21, 0x19, 0x71, 0xbf, 0x60, 0xd9,
	0x07, 0x8b, 0xae, 0x92, 0xfb, 0x55, 0xe8, 0x2f, 0x56, 0xa1, 0xff, 0x67, 0x15, 0xfa, 0xdf, 0xd6,
	0xa1, 0xb7, 0x58, 0x87, 0xde, 0xaf, 0x75, 0xe8, 0x7d, 0x7e, 0xf3, 0x2f, 0x73, 0xcc, 0xdd, 0x6a,
	0xdb, 0xbd, 0xce, 0x1f, 0xda, 0xa5, 0x7c, 0xfd, 0x37, 0x00, 0x00, 0xff, 0xff, 0x01, 0xd9, 0xc6,
	0xac, 0xf9, 0x02, 0x00, 0x00,
}

func (m *Mailbox) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Mailbox) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Mailbox) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LocalDomain != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.LocalDomain))
		i--
		dAtA[i] = 0x50
	}
	if m.RequiredHook != nil {
		{
			size := m.RequiredHook.Size()
			i -= size
			if _, err := m.RequiredHook.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if m.DefaultHook != nil {
		{
			size := m.DefaultHook.Size()
			i -= size
			if _, err := m.DefaultHook.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	{
		size := m.DefaultIsm.Size()
		i -= size
		if _, err := m.DefaultIsm.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.MessageReceived != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.MessageReceived))
		i--
		dAtA[i] = 0x20
	}
	if m.MessageSent != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.MessageSent))
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
func (m *Mailbox) Size() (n int) {
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
	if m.MessageSent != 0 {
		n += 1 + sovTypes(uint64(m.MessageSent))
	}
	if m.MessageReceived != 0 {
		n += 1 + sovTypes(uint64(m.MessageReceived))
	}
	l = m.DefaultIsm.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.DefaultHook != nil {
		l = m.DefaultHook.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.RequiredHook != nil {
		l = m.RequiredHook.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.LocalDomain != 0 {
		n += 1 + sovTypes(uint64(m.LocalDomain))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Mailbox) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Mailbox: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Mailbox: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field MessageSent", wireType)
			}
			m.MessageSent = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MessageSent |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageReceived", wireType)
			}
			m.MessageReceived = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MessageReceived |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultIsm", wireType)
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
			if err := m.DefaultIsm.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DefaultHook", wireType)
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
			var v github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress
			m.DefaultHook = &v
			if err := m.DefaultHook.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequiredHook", wireType)
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
			var v github_com_bcp_innovations_hyperlane_cosmos_util.HexAddress
			m.RequiredHook = &v
			if err := m.RequiredHook.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LocalDomain", wireType)
			}
			m.LocalDomain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LocalDomain |= uint32(b&0x7F) << shift
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
