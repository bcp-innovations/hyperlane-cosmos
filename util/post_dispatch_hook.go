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

	metadata.Address = raw[STANDARD_HOOK_METADATA_REFUND_ADDRESS_OFFSET:STANDARD_HOOK_METADATA_MIN_METADATA_LENGTH]

	return metadata, nil
}

func (metadata StandardHookMetadata) Bytes() []byte {
	variant := make([]byte, 2)
	binary.BigEndian.PutUint16(variant, metadata.Variant)

	value := metadata.Value.BigInt().Bytes()
	// TODO return error
	valuePadding := make([]byte, 32-len(value))

	gasLimit := metadata.GasLimit.BigInt().Bytes()
	gasPadding := make([]byte, 32-len(gasLimit))

	// TODO assert length = 20
	metadata.Address.Bytes()

	return slices.Concat(
		variant,
		valuePadding,
		value,
		gasPadding,
		gasLimit,
		metadata.Address.Bytes(),
	)
}

func CreateEmptyStandardHookMetadata() (StandardHookMetadata, error) {
	mockBytes := make([]byte, STANDARD_HOOK_METADATA_MIN_METADATA_LENGTH)

	binary.BigEndian.PutUint16(mockBytes[STANDARD_HOOK_METADATA_VARIANT_OFFSET:STANDARD_HOOK_METADATA_MSG_VALUE_OFFSET], 0)

	value := big.NewInt(0)
	copy(mockBytes[STANDARD_HOOK_METADATA_MSG_VALUE_OFFSET:STANDARD_HOOK_METADATA_GAS_LIMIT_OFFSET], value.FillBytes(make([]byte, 32)))

	gasLimit := big.NewInt(0)
	copy(mockBytes[STANDARD_HOOK_METADATA_GAS_LIMIT_OFFSET:STANDARD_HOOK_METADATA_REFUND_ADDRESS_OFFSET], gasLimit.FillBytes(make([]byte, 32)))

	copy(mockBytes[STANDARD_HOOK_METADATA_REFUND_ADDRESS_OFFSET:STANDARD_HOOK_METADATA_MIN_METADATA_LENGTH], sdk.AccAddress("mocked_addr________").Bytes())

	return ParseStandardHookMetadata(mockBytes)
}
