package types

import "cosmossdk.io/collections"

const ModuleName = "igp"

const (
	UNUSED uint32 = iota
	ROUTING
	AGGREGATION
	LEGACY_MULTISIG
	MERKLE_ROOT_MULTISIG
	MESSAGE_ID_MULTISIG
	NULL // used with relayer carrying no metadata
	CCIP_READ
	ARB_L2_TO_L1
	WEIGHTED_MERKLE_ROOT_MULTISIG
	WEIGHTED_MESSAGE_ID_MULTISIG
	OP_L2_TO_L1
)

var (
	ParamsKey = collections.NewPrefix(0)
)
