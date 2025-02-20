package util

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// <class> / <type> / <custom> <uint64-id>
// 24
// class: ISM,PostDistpatchHook,Receiver
// type: multiSigISM, noopISM -- IGP, Posthooks,  --- Collateral, Synthetic

// _post_dispatch_hooks: 0,1,2,3,4,5
// 3rdparty_post_dispatch_hooks: 0,1,2,3,4,5

// mappingId <>
// core -> Verify(ismId(0,1,2,3,4,5))
// core -> Verify(ismHexAddress(type,class))

// verify-modiule
/*
func Verify(ismAddress) {
if ismAddress.getClass != "me" {
return nil

Register() {

}

/query/multisig

/query/3rdparty-sig

/hyperlane/isms/id


*/

var registeredFactoryClasses = map[string]struct{}{}

type HexAddressFactory struct {
	class string
}

func NewHexAddressFactory(class string) HexAddressFactory {
	return HexAddressFactory{class: class}
}

func (h HexAddressFactory) IsMember(id HexAddress) bool {
	return id.GetClass() == h.class
}

func (h HexAddressFactory) GenerateId(internalType uint32, internalId uint64) HexAddress {
	return HexAddress{}
}

func (h HexAddressFactory) GetClass() string {
	return h.class
}

type HexAddress [32]byte

func (h HexAddress) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(h[:]))
}

func (h HexAddress) Bytes() []byte {
	return h[:]
}

func (h HexAddress) GetInternalId() uint64 {
	// TODO stub
	return 0
}

func (h HexAddress) GetClass() string {
	// TODO stub
	return "coreism"
}

func (h HexAddress) GetType() uint32 {
	// TODO stub
	return 0
}

func NewHexAddress(class, _type string, id uint64) HexAddress {
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
