package types

import "cosmossdk.io/errors"

var (
	ErrMailboxDoesNotExist               = errors.Register(SubModuleName, 1, "mailbox does not exist")
	ErrSenderIsNotDesignatedMailbox      = errors.Register(SubModuleName, 2, "sender is not designated mailbox")
	ErrHookDoesNotExistOrIsNotRegistered = errors.Register(SubModuleName, 3, "hook does not exist or isn't registered")
	ErrUnauthorized                      = errors.Register(SubModuleName, 4, "unauthorized")
	ErrInvalidOwner                      = errors.Register(SubModuleName, 5, "invalid owner")
	ErrInvalidAggregationHook            = errors.Register(SubModuleName, 6, "invalid aggregation hook")
	ErrInvalidRateLimitedHook            = errors.Register(SubModuleName, 7, "invalid rate limited hook")
	ErrRateLimitNotConfigured            = errors.Register(SubModuleName, 8, "rate limit not configured")
	ErrRateLimitExceeded                 = errors.Register(SubModuleName, 9, "rate limit exceeded")
	ErrInvalidDispatchedMessage          = errors.Register(SubModuleName, 11, "invalid dispatched message")
	ErrRateLimitNotSet                   = errors.Register(SubModuleName, 12, "rate limit not set")
	ErrInvalidPausableHook               = errors.Register(SubModuleName, 13, "invalid pausable hook")
	ErrPausableHookPaused                = errors.Register(SubModuleName, 14, "pausable hook paused")
)
