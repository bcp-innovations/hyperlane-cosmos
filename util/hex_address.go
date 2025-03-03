package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
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

Encoding:
- 0:20 bytes are the module identifier
- 20:24 bytes is the module type. This is used to identify the correct module implementation for this type in the router. See `util./router.go`
- 24:32 bytes is the internal id. Can be used for internal collection storage.
*/
const (
	HEX_ADDRESS_LENGTH = 32
	module_length      = 20
	type_length        = 4
	id_length          = 8

	module_offset              = 0
	type_offset                = module_offset + module_length // 20
	id_offset                  = type_offset + type_length     // 24
	ENCODED_HEX_ADDRESS_LENGTH = 66                            // raw 32 bytes of the address stored as a hex string (0x + 32 bytes hex encoded = 66 bytes)
)

type HexAddress [HEX_ADDRESS_LENGTH]byte

func (h HexAddress) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(h[:]))
}

func (h HexAddress) Bytes() []byte {
	return h[:]
}

func (h HexAddress) IsZeroAddress() bool {
	emptyByteVar := make([]byte, HEX_ADDRESS_LENGTH)
	return bytes.Equal(h[:], emptyByteVar)
}

func (h HexAddress) GetInternalId() uint64 {
	return binary.BigEndian.Uint64(h[id_offset : id_offset+id_length])
}

func (h HexAddress) GetType() uint32 {
	return binary.BigEndian.Uint32(h[type_offset : type_offset+type_length])
}

func NewZeroAddress() HexAddress {
	return HexAddress{}
}

func DecodeHexAddress(s string) (HexAddress, error) {
	s = strings.TrimPrefix(s, "0x")

	// hex encodes two characters per byte
	if len(s) != HEX_ADDRESS_LENGTH*2 {
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

func CreateMockHexAddress(identifier string, id int64) HexAddress {
	idBytes := make([]byte, id_length)
	binary.BigEndian.PutUint64(idBytes, uint64(id))
	message := append([]byte(identifier), idBytes...)
	return sha256.Sum256(message)
}

func ParseFromCosmosAcc(cosmosAcc string) (HexAddress, error) {
	bech32, err := sdk.AccAddressFromBech32(cosmosAcc)
	if err != nil {
		return [HEX_ADDRESS_LENGTH]byte{}, err
	}

	if len(bech32) > HEX_ADDRESS_LENGTH {
		return HexAddress{}, errors.New("invalid length")
	}

	hexAddressBytes := make([]byte, HEX_ADDRESS_LENGTH)
	copy(hexAddressBytes[HEX_ADDRESS_LENGTH-len(bech32):], bech32)

	return HexAddress(hexAddressBytes), nil
}

func GenerateHexAddress(moduleSpecifier [module_length]byte, internalType uint32, internalId uint64) HexAddress {
	internalTypeBytes := make([]byte, type_length)
	binary.BigEndian.PutUint32(internalTypeBytes, internalType)

	internalIdBytes := make([]byte, id_length)
	binary.BigEndian.PutUint64(internalIdBytes, internalId)

	return HexAddress(slices.Concat(moduleSpecifier[:], internalTypeBytes, internalIdBytes))
}

// Custom Proto Type Implementation below
//
// For custom type serialization we prefer readability to storage space
// In the entire CosmosSDK ecosystem, there is always the string representation used for addresses.
// We therefore store the 66 (0x + 32 bytes hex encoded = 66 bytes) hex representation of the address.

func (h HexAddress) Marshal() ([]byte, error) {
	return []byte(h.String()), nil
}

func (h *HexAddress) MarshalTo(data []byte) (n int, err error) {
	n = copy(data, h.String())
	if n != ENCODED_HEX_ADDRESS_LENGTH {
		return n, fmt.Errorf("invalid hex address length: %d", n)
	}
	return n, nil
}

func (h *HexAddress) Unmarshal(data []byte) error {
	if len(data) != ENCODED_HEX_ADDRESS_LENGTH {
		return errors.New("invalid hex address length")
	}
	addr, err := DecodeHexAddress(string(data))
	if err != nil {
		return err
	}
	copy(h[:], addr.Bytes())
	return nil
}

func (h *HexAddress) Size() int {
	return ENCODED_HEX_ADDRESS_LENGTH
}

func (h HexAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *HexAddress) UnmarshalJSON(data []byte) error {
	address, err := DecodeHexAddress(string(data))
	if err != nil {
		return err
	}
	copy(h[:], address.Bytes())
	return nil
}

func (h HexAddress) Compare(other HexAddress) int {
	return bytes.Compare(h.Bytes(), other.Bytes())
}

func (h HexAddress) Equal(other HexAddress) bool {
	return bytes.Equal(h.Bytes(), other.Bytes())
}
