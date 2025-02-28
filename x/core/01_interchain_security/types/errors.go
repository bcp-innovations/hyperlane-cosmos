package types

import "cosmossdk.io/errors"

var (
	ErrUnexpectedError              = errors.Register(SubModuleName, 1, "unexpected error")
	ErrInvalidMultisigConfiguration = errors.Register(SubModuleName, 2, "invalid multisig configuration")
)
