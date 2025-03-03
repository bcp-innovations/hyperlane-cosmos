package types

import "cosmossdk.io/errors"

var ErrNoReceiverISM = errors.New(ModuleName, 1, "no receiver ISM")

// Invalid mailbox state, a required hook must always be set
var ErrRequiredHookNotSet = errors.New(ModuleName, 2, "required hook not set")

// Invalid mailbox state, a required hook must always be set
var ErrDefaultHookNotSet = errors.New(ModuleName, 3, "default hook not set")
