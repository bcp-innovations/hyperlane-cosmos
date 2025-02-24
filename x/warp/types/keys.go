package types

import "cosmossdk.io/collections"

const ModuleName = "warp"

var (
	ParamsKey          = collections.NewPrefix(0)
	HypTokenKey        = collections.NewPrefix(1)
	HypTokensCountKey  = collections.NewPrefix(2)
	EnrolledRoutersKey = collections.NewPrefix(3)
)

const HEX_ADDRESS_CLASS_IDENTIFIER = "hyperlanewarp"
