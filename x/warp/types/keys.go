package types

import "cosmossdk.io/collections"

const ModuleName = "warp"

var (
	ParamsKey          = collections.NewPrefix(0)
	HypTokenKey        = collections.NewPrefix(1)
	EnrolledRoutersKey = collections.NewPrefix(2)
)
