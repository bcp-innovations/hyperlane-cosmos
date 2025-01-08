package types

import "cosmossdk.io/errors"

var (
	ErrNotEnoughCollateral = errors.Register(ModuleName, 2, "not enough collateral")
)
