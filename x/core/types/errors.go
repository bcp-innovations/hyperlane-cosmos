package types

import "cosmossdk.io/errors"

// ErrDuplicateAddress error if there is a duplicate address
var ErrDuplicateAddress = errors.Register(ModuleName, 2, "duplicate address")

var ErrNoReceiverISM = errors.New(ModuleName, 3, "no receiver ISM")

// ErrMultipleReceiverIsm should not happen if every Module is wired up correctly
var ErrMultipleReceiverIsm = errors.New(ModuleName, 4, "multiple receiver ISM")

// Invalid mailbox state, a required hook must always be set
var ErrRequiredHookNotSet = errors.New(ModuleName, 5, "required hook not set")

// Invalid mailbox state, a required hook must always be set
var ErrDefaultHookNotSet = errors.New(ModuleName, 6, "default hook not set")
