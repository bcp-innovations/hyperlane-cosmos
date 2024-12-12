package util

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type HexAddress [32]byte

func (h HexAddress) String() string {
	return fmt.Sprintf("0x%s", hex.EncodeToString(h[:]))
}

func (h HexAddress) Bytes() []byte {
	return h[:]
}

func DecodeHexAddress(s string) (HexAddress, error) {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

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
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return b, nil
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
