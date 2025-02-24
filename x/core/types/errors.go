package types

import "cosmossdk.io/errors"

// ErrDuplicateAddress error if there is a duplicate address
var ErrDuplicateAddress = errors.Register(ModuleName, 2, "duplicate address")

var ErrNoReceiverISM = errors.New(ModuleName, 3, "no receiver ISM")

// ErrMultipleReceiverIsm should not happen if every Module is wired up correctly
var ErrMultipleReceiverIsm = errors.New(ModuleName, 4, "multiple receiver ISM")
