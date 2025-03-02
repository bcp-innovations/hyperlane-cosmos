package util

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StandardHookMetadata is application specific, thereby different of the EVM implementation.
type StandardHookMetadata struct {
	Address sdk.AccAddress

	GasLimit math.Int

	CustomHookMetadata []byte
}
