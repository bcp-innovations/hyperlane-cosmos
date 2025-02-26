package util

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"slices"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	STANDARD_HOOK_METADATA_VARIANT_OFFSET        = 0
	STANDARD_HOOK_METADATA_MSG_VALUE_OFFSET      = 2
	STANDARD_HOOK_METADATA_GAS_LIMIT_OFFSET      = 34
	STANDARD_HOOK_METADATA_REFUND_ADDRESS_OFFSET = 66
	STANDARD_HOOK_METADATA_MIN_METADATA_LENGTH   = 86

	STANDARD_HOOK_METADATA_VARIANT = 1
)

type StandardHookMetadata struct {
	Variant  uint16
	Value    math.Int
	GasLimit math.Int
	Address  sdk.AccAddress
}

func ParseStandardHookMetadata(raw []byte) (StandardHookMetadata, error) {
	metadata := StandardHookMetadata{}

	if len(raw) < STANDARD_HOOK_METADATA_MIN_METADATA_LENGTH {
		return metadata, fmt.Errorf("invalid standard hook metadata")
	}

	metadata.Variant = binary.BigEndian.Uint16(raw[STANDARD_HOOK_METADATA_VARIANT_OFFSET:STANDARD_HOOK_METADATA_MSG_VALUE_OFFSET])
	valueBig := new(big.Int)
	valueBig.SetBytes(raw[STANDARD_HOOK_METADATA_MSG_VALUE_OFFSET:STANDARD_HOOK_METADATA_GAS_LIMIT_OFFSET])
	metadata.Value = math.NewIntFromBigInt(valueBig)

	gasLimitBig := new(big.Int)
	gasLimitBig.SetBytes(raw[STANDARD_HOOK_METADATA_GAS_LIMIT_OFFSET:STANDARD_HOOK_METADATA_REFUND_ADDRESS_OFFSET])
	metadata.GasLimit = math.NewIntFromBigInt(gasLimitBig)

	// TODO only 20 bytes cosmos addresses supported
	metadata.Address = raw[STANDARD_HOOK_METADATA_REFUND_ADDRESS_OFFSET+12 : STANDARD_HOOK_METADATA_MIN_METADATA_LENGTH]

	return metadata, nil
}

func (metadata StandardHookMetadata) Bytes() []byte {
	variant := make([]byte, 4)
	binary.BigEndian.PutUint16(variant, metadata.Variant)

	value := metadata.Value.BigInt().Bytes()
	gasLimit := metadata.GasLimit.BigInt().Bytes()

	// TODO return error

	return slices.Concat(
		variant,
		value,
		gasLimit,
	)
}
