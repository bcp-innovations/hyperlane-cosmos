package types

import (
	"cosmossdk.io/math"
)

var (
	PostDispatchHooksKey             = []byte{SubModuleId, 1}
	PostDispatchHooksSequenceKey     = []byte{SubModuleId, 2}
	InterchainGasPaymasterConfigsKey = []byte{SubModuleId, 3}
	MerkleTreeHooksKey               = []byte{SubModuleId, 4}
)

const (
	SubModuleName       = "post_dispatch"
	SubModuleId   uint8 = 2
)

var TokenExchangeRateScale = math.NewInt(1e10)

const (
	POST_DISPATCH_HOOK_TYPE_UNUSED uint8 = iota
	POST_DISPATCH_HOOK_TYPE_ROUTING
	POST_DISPATCH_HOOK_TYPE_AGGREGATION
	POST_DISPATCH_HOOK_TYPE_MERKLE_TREE
	POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER
	POST_DISPATCH_HOOK_TYPE_FALLBACK_ROUTING
	POST_DISPATCH_HOOK_TYPE_ID_AUTH_ISM
	POST_DISPATCH_HOOK_TYPE_PAUSABLE
	POST_DISPATCH_HOOK_TYPE_PROTOCOL_FEE
	POST_DISPATCH_HOOK_TYPE_LAYER_ZERO_V1
	POST_DISPATCH_HOOK_TYPE_RATE_LIMITED
	POST_DISPATCH_HOOK_TYPE_ARB_L2_TO_L1
	POST_DISPATCH_HOOK_TYPE_OP_L2_TO_L1
)
