package types

import "cosmossdk.io/collections"

const ModuleName = "igp"

const Denom = "tkyve"

var (
	ParamsKey      = collections.NewPrefix(0)
	IgpKey         = collections.NewPrefix(1)
	IgpSequenceKey = collections.NewPrefix(2)
)

const (
	MaxDestinationGasConfigs        = 50
	TokenExchangeRateScale   uint64 = 1e10
)
