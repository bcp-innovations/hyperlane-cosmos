package types

import "cosmossdk.io/errors"

var (
	ErrUnexpectedError              = errors.Register(SubModuleName, 1, "unexpected error")
	ErrInvalidMultisigConfiguration = errors.Register(SubModuleName, 2, "invalid multisig configuration")
	ErrInvalidAnnounce              = errors.Register(SubModuleName, 3, "invalid announce")
	ErrMailboxDoesNotExist          = errors.Register(SubModuleName, 4, "mailbox does not exist")
	ErrInvalidSignature             = errors.Register(SubModuleName, 5, "invalid signature")
	ErrInvalidISMType               = errors.Register(SubModuleName, 6, "invalid ism type")
	ErrUnkownIsmId                  = errors.Register(SubModuleName, 7, "unknown ism id")
	ErrNoRouteFound                 = errors.Register(SubModuleName, 8, "no route found")
	ErrUnauthorized                 = errors.Register(SubModuleName, 9, "unauthorized")
	ErrInvalidOwner                 = errors.Register(SubModuleName, 10, "invalid owner")
	ErrDuplicatedDomains            = errors.Register(SubModuleName, 11, "route for domain already exists")
	ErrInvalidThreshold             = errors.Register(SubModuleName, 12, "invalid threshold")
	ErrInvalidISM                   = errors.Register(SubModuleName, 13, "invalid ism")
	ErrInsufficientVerifications    = errors.Register(SubModuleName, 14, "insufficient verifications")
)
