package types

import "cosmossdk.io/errors"

var (
	ErrUnexpectedError              = errors.Register(SubModuleName, 1, "unexpected error")
	ErrInvalidMultisigConfiguration = errors.Register(SubModuleName, 2, "invalid multisig configuration")
	ErrInvalidAnnounce              = errors.Register(SubModuleName, 3, "invalid announce")
	ErrMailboxDoesNotExist          = errors.Register(SubModuleName, 4, "mailbox does not exist")
	ErrInvalidSignature             = errors.Register(SubModuleName, 5, "invalid signature")
)
