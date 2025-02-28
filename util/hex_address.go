package util

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

*/

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

// Custom Proto Type Implementation below
//
// For custom type serialization we prefer readability to storage space
// In the entire CosmosSDK ecosystem, there is always the string representation used for addresses.
// We therefore store the 66 (0x + 32 bytes hex encoded = 66 bytes) hex representation of the address.

const HEX_ADDRESS_LENGTH = 66

func (t HexAddress) Marshal() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *HexAddress) MarshalTo(data []byte) (n int, err error) {
	n = copy(data, t.String())
	if n != HEX_ADDRESS_LENGTH {
		return n, fmt.Errorf("invalid hex address length: %d", n)
	}
	return n, nil
}

func (t *HexAddress) Unmarshal(data []byte) error {
	if len(data) != HEX_ADDRESS_LENGTH {
		return errors.New("invalid hex address length")
	}
	addr, err := DecodeHexAddress(string(data))
	if err != nil {
		return err
	}
	copy(t[:], addr.Bytes())
	return nil
}

func (t *HexAddress) Size() int {
	return HEX_ADDRESS_LENGTH
}

func (t HexAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *HexAddress) UnmarshalJSON(data []byte) error {
	address, err := DecodeHexAddress(string(data))
	if err != nil {
		return err
	}
	copy(t[:], address.Bytes())
	return nil
}

func (t HexAddress) Compare(other HexAddress) int {
	return bytes.Compare(t.Bytes(), other.Bytes())
}

func (t HexAddress) Equal(other HexAddress) bool {
	return bytes.Equal(t.Bytes(), other.Bytes())
}
