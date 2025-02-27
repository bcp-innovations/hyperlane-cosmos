package types

import "cosmossdk.io/errors"

var (
	ErrMailboxDoesNotExist               = errors.Register(SubModuleName, 1, "mailbox does not exist")
	ErrSenderIsNotDesignatedMailbox      = errors.Register(SubModuleName, 2, "sender is not designated mailbox")
	ErrHookDoesNotExistOrIsNotRegistered = errors.Register(SubModuleName, 2, "hook does not exist or isn't registered")
)
