// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: hyperlane/core/interchain_security/v1/genesis.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
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

// GenesisState defines the 01_interchain_security submodule's genesis state.
type GenesisState struct {
	// accounts are the accounts present at genesis.
	Isms                      []*types.Any                             `protobuf:"bytes,1,rep,name=isms,proto3" json:"isms,omitempty"`
	ValidatorStorageLocations []GenesisValidatorStorageLocationWrapper `protobuf:"bytes,2,rep,name=validator_storage_locations,json=validatorStorageLocations,proto3" json:"validator_storage_locations"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_908eb45c3c27ef24, []int{0}
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

func (m *GenesisState) GetIsms() []*types.Any {
	if m != nil {
		return m.Isms
	}
	return nil
}

func (m *GenesisState) GetValidatorStorageLocations() []GenesisValidatorStorageLocationWrapper {
	if m != nil {
		return m.ValidatorStorageLocations
	}
	return nil
}

// GenesisValidatorStorageLocationWrapper stores the information for
// validator, mailbox and storage-location which validators have announced
type GenesisValidatorStorageLocationWrapper struct {
	MailboxId        uint64 `protobuf:"varint,1,opt,name=mailbox_id,json=mailboxId,proto3" json:"mailbox_id,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Index            uint64 `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	StorageLocation  string `protobuf:"bytes,4,opt,name=storage_location,json=storageLocation,proto3" json:"storage_location,omitempty"`
}

func (m *GenesisValidatorStorageLocationWrapper) Reset() {
	*m = GenesisValidatorStorageLocationWrapper{}
}
func (m *GenesisValidatorStorageLocationWrapper) String() string { return proto.CompactTextString(m) }
func (*GenesisValidatorStorageLocationWrapper) ProtoMessage()    {}
func (*GenesisValidatorStorageLocationWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_908eb45c3c27ef24, []int{1}
}
func (m *GenesisValidatorStorageLocationWrapper) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisValidatorStorageLocationWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisValidatorStorageLocationWrapper.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisValidatorStorageLocationWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisValidatorStorageLocationWrapper.Merge(m, src)
}
func (m *GenesisValidatorStorageLocationWrapper) XXX_Size() int {
	return m.Size()
}
func (m *GenesisValidatorStorageLocationWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisValidatorStorageLocationWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisValidatorStorageLocationWrapper proto.InternalMessageInfo

func (m *GenesisValidatorStorageLocationWrapper) GetMailboxId() uint64 {
	if m != nil {
		return m.MailboxId
	}
	return 0
}

func (m *GenesisValidatorStorageLocationWrapper) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *GenesisValidatorStorageLocationWrapper) GetIndex() uint64 {
	if m != nil {
		return m.Index
	}
	return 0
}

func (m *GenesisValidatorStorageLocationWrapper) GetStorageLocation() string {
	if m != nil {
		return m.StorageLocation
	}
	return ""
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "hyperlane.core.interchain_security.v1.GenesisState")
	proto.RegisterType((*GenesisValidatorStorageLocationWrapper)(nil), "hyperlane.core.interchain_security.v1.GenesisValidatorStorageLocationWrapper")
}

func init() {
	proto.RegisterFile("hyperlane/core/interchain_security/v1/genesis.proto", fileDescriptor_908eb45c3c27ef24)
}

var fileDescriptor_908eb45c3c27ef24 = []byte{
	// 391 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xd1, 0xca, 0xd3, 0x30,
	0x14, 0xc7, 0x1b, 0xbf, 0x2a, 0x2c, 0x0a, 0xce, 0xb2, 0x8b, 0x6e, 0x62, 0x1d, 0x03, 0xa5, 0x22,
	0x4b, 0x9c, 0x7b, 0x82, 0xed, 0x46, 0x04, 0x45, 0xe8, 0x40, 0xc1, 0x9b, 0x92, 0xb6, 0xb1, 0x0b,
	0xb4, 0x39, 0x25, 0xc9, 0xca, 0xfa, 0x16, 0x82, 0xcf, 0xe2, 0x3b, 0xec, 0x4a, 0x76, 0xe9, 0x95,
	0xc8, 0xf6, 0x22, 0xb2, 0x66, 0x9b, 0x28, 0xfb, 0x60, 0x77, 0x27, 0x27, 0xe7, 0xff, 0xcb, 0xf9,
	0x9f, 0x1c, 0x3c, 0x5d, 0x36, 0x15, 0x57, 0x05, 0x93, 0x9c, 0xa6, 0xa0, 0x38, 0x15, 0xd2, 0x70,
	0x95, 0x2e, 0x99, 0x90, 0xb1, 0xe6, 0xe9, 0x4a, 0x09, 0xd3, 0xd0, 0x7a, 0x42, 0x73, 0x2e, 0xb9,
	0x16, 0x9a, 0x54, 0x0a, 0x0c, 0x78, 0xcf, 0xce, 0x22, 0x72, 0x10, 0x91, 0x0b, 0x22, 0x52, 0x4f,
	0x06, 0xfd, 0x1c, 0x20, 0x2f, 0x38, 0x6d, 0x45, 0xc9, 0xea, 0x0b, 0x65, 0xb2, 0xb1, 0x84, 0x41,
	0x2f, 0x87, 0x1c, 0xda, 0x90, 0x1e, 0x22, 0x9b, 0x1d, 0xfd, 0x40, 0xf8, 0xc1, 0x1b, 0xfb, 0xd2,
	0xc2, 0x30, 0xc3, 0xbd, 0x10, 0xbb, 0x42, 0x97, 0xda, 0x47, 0xc3, 0x9b, 0xf0, 0xfe, 0xeb, 0x1e,
	0xb1, 0x40, 0x72, 0x02, 0x92, 0x99, 0x6c, 0xa2, 0xb6, 0xc2, 0xfb, 0x86, 0xf0, 0xe3, 0x9a, 0x15,
	0x22, 0x63, 0x06, 0x54, 0xac, 0x0d, 0x28, 0x96, 0xf3, 0xb8, 0x80, 0x94, 0x19, 0x01, 0x52, 0xfb,
	0x77, 0x5a, 0xc2, 0x7b, 0x72, 0x55, 0xe7, 0xe4, 0xd8, 0xc4, 0xc7, 0x13, 0x70, 0x61, 0x79, 0xef,
	0x8e, 0xb8, 0x4f, 0x8a, 0x55, 0x15, 0x57, 0x73, 0x77, 0xf3, 0xeb, 0xa9, 0x13, 0xf5, 0xeb, 0x5b,
	0xca, 0xf4, 0xe8, 0x3b, 0xc2, 0xcf, 0xaf, 0x63, 0x79, 0x4f, 0x30, 0x2e, 0x99, 0x28, 0x12, 0x58,
	0xc7, 0x22, 0xf3, 0xd1, 0x10, 0x85, 0x6e, 0xd4, 0x39, 0x66, 0xde, 0x66, 0xde, 0x4b, 0xfc, 0xe8,
	0xaf, 0x3d, 0x96, 0x65, 0x8a, 0xeb, 0x83, 0x29, 0x14, 0x76, 0xa2, 0xee, 0xf9, 0x62, 0x66, 0xf3,
	0x5e, 0x0f, 0xdf, 0x15, 0x32, 0xe3, 0x6b, 0xff, 0xa6, 0xc5, 0xd8, 0x83, 0xf7, 0x02, 0x77, 0xff,
	0x9f, 0x8b, 0xef, 0xb6, 0x84, 0x87, 0xfa, 0xdf, 0x9e, 0xe6, 0x62, 0xb3, 0x0b, 0xd0, 0x76, 0x17,
	0xa0, 0xdf, 0xbb, 0x00, 0x7d, 0xdd, 0x07, 0xce, 0x76, 0x1f, 0x38, 0x3f, 0xf7, 0x81, 0xf3, 0xf9,
	0x43, 0x2e, 0xcc, 0x72, 0x95, 0x90, 0x14, 0x4a, 0x9a, 0xa4, 0xd5, 0x58, 0x48, 0x09, 0xb5, 0x75,
	0x4b, 0xcf, 0xb3, 0x1d, 0xa7, 0xa0, 0x4b, 0xd0, 0x74, 0x6d, 0x77, 0xea, 0xd5, 0x24, 0xbe, 0xb4,
	0x56, 0xa6, 0xa9, 0xb8, 0x4e, 0xee, 0xb5, 0x9f, 0x39, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0xcd,
	0xd1, 0x1a, 0xb1, 0x89, 0x02, 0x00, 0x00,
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
	if len(m.ValidatorStorageLocations) > 0 {
		for iNdEx := len(m.ValidatorStorageLocations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ValidatorStorageLocations[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Isms) > 0 {
		for iNdEx := len(m.Isms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Isms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *GenesisValidatorStorageLocationWrapper) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisValidatorStorageLocationWrapper) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisValidatorStorageLocationWrapper) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.StorageLocation) > 0 {
		i -= len(m.StorageLocation)
		copy(dAtA[i:], m.StorageLocation)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.StorageLocation)))
		i--
		dAtA[i] = 0x22
	}
	if m.Index != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x18
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.MailboxId != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.MailboxId))
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
	if len(m.Isms) > 0 {
		for _, e := range m.Isms {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.ValidatorStorageLocations) > 0 {
		for _, e := range m.ValidatorStorageLocations {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *GenesisValidatorStorageLocationWrapper) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MailboxId != 0 {
		n += 1 + sovGenesis(uint64(m.MailboxId))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Index != 0 {
		n += 1 + sovGenesis(uint64(m.Index))
	}
	l = len(m.StorageLocation)
	if l > 0 {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Isms", wireType)
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
			m.Isms = append(m.Isms, &types.Any{})
			if err := m.Isms[len(m.Isms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorStorageLocations", wireType)
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
			m.ValidatorStorageLocations = append(m.ValidatorStorageLocations, GenesisValidatorStorageLocationWrapper{})
			if err := m.ValidatorStorageLocations[len(m.ValidatorStorageLocations)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *GenesisValidatorStorageLocationWrapper) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: GenesisValidatorStorageLocationWrapper: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisValidatorStorageLocationWrapper: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MailboxId", wireType)
			}
			m.MailboxId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MailboxId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
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
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StorageLocation", wireType)
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
			m.StorageLocation = string(dAtA[iNdEx:postIndex])
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
