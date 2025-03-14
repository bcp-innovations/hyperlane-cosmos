package types

import "cosmossdk.io/errors"

var (
	ErrNotEnoughCollateral = errors.Register(ModuleName, 1, "not enough collateral")
	ErrTokenNotFound       = errors.Register(ModuleName, 2, "token not found")
)
