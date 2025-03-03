package util

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StandardHookMetadata is application specific, thereby different to the EVM implementation.
type StandardHookMetadata struct {
	// Address is used to determine the message sender.
	// E.g., IGP uses the address to determine the fee payer.
	Address sdk.AccAddress

	// GasLimit of the destination transaction.
	GasLimit math.Int

	// CustomHookMetadata can be used to pass custom data to a PostDispatch hook.
	// The hook is responsible for validating and decoding the custom metadata.
	CustomHookMetadata []byte
}
