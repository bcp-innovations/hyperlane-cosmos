package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*

SPEC: HexAddress

The HexAddress mimics an evm-compatible address for a smart contract.
Due to the nature of cosmos, addresses must be created differently.

Requirements:
- HexAddresses must be unique across all cosmos modules interacting with Hyperlane

Structure
- The HexAddress has 32 bytes and is used for external communication
- For internal usage and storage an uint64 is totally sufficient

HexAddress: <module-specifier (20 byte)> <type (4 byte)> <internal-id (8 byte)>

The struct provides functions to encode and decode the information stored within the address.

To ensure global uniqueness, the HexAddressFactory should be used. It is initialized once per Keeper
and keeps global track of all registered module specifiers.

*/

// Hex Address Factory

var registeredFactoryClasses = map[string]int{}

type HexAddressFactory struct {
	class string
}

func NewHexAddressFactory(class string) (HexAddressFactory, error) {
	// Keeper is called twice, so if the function called more than 2 times
	// one can assume that the developer misconfigured the module.
	// TODO
	//if count, ok := registeredFactoryClasses[class]; ok && count > 1 {
	//	return HexAddressFactory{}, fmt.Errorf("factory class %s already registered", class)
	//}
	//registeredFactoryClasses[class] += 1
	_ = registeredFactoryClasses

	if len(class) > 20 {
		return HexAddressFactory{}, fmt.Errorf("factory class %s too long", class)
	}

	return HexAddressFactory{class: class}, nil
}

func (h HexAddressFactory) IsClassMember(id HexAddress) bool {
	return id.GetClass() == h.class
}

func (h HexAddressFactory) GetClass() string {
	return h.class
}

func (h HexAddressFactory) GenerateId(internalType uint32, internalId uint64) HexAddress {
	address := make([]byte, 20)
	copy(address, h.class)

	internalTypeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(internalTypeBytes, internalType)

	internalIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(internalIdBytes, internalId)

	return HexAddress(slices.Concat(address, internalTypeBytes, internalIdBytes))
}

// Hex Address

type HexAddress [32]byte

func (h HexAddress) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(h[:]))
}

func (h HexAddress) Bytes() []byte {
	return h[:]
}

func (h HexAddress) IsZeroAddress() bool {
	emptyByteVar := make([]byte, 32)
	return bytes.Equal(h[:], emptyByteVar)
}

func (h HexAddress) GetInternalId() uint64 {
	return binary.BigEndian.Uint64(h[24:32])
}

func (h HexAddress) GetClass() string {
	// Trim empty bytes.
	return strings.Trim(string(h[:20]), "\x00")
}

func (h HexAddress) GetType() uint32 {
	return binary.BigEndian.Uint32(h[20:24])
}

func NewZeroAddress() HexAddress {
	return HexAddress{}
}

func DecodeHexAddress(s string) (HexAddress, error) {
	s = strings.TrimPrefix(s, "0x")

	if len(s) != 64 {
		return HexAddress{}, errors.New("invalid hex address length")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return HexAddress{}, err
	}
	return HexAddress(b), nil
}

func DecodeEthHex(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func EncodeEthHex(b []byte) string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(b))
}

func CreateHexAddress(identifier string, id int64) HexAddress {
	idBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idBytes, uint64(id))
	message := append([]byte(identifier), idBytes...)
	return sha256.Sum256(message)
}

func ParseFromCosmosAcc(cosmosAcc string) (HexAddress, error) {
	bech32, err := sdk.AccAddressFromBech32(cosmosAcc)
	if err != nil {
		return [32]byte{}, err
	}

	if len(bech32) > 32 {
		return HexAddress{}, errors.New("invalid length")
	}

	hexAddressBytes := make([]byte, 32)
	copy(hexAddressBytes[32-len(bech32):], bech32)

	return HexAddress(hexAddressBytes), nil
}
