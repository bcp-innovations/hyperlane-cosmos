package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var (
	InterchainGasPaymasterKey        = []byte{SubModuleId, 1}
	InterchainGasPaymasterConfigsKey = []byte{SubModuleId, 3}
	MerkleTreeHooksKey               = []byte{SubModuleId, 4}
	NoopHooksKey                     = []byte{SubModuleId, 5}
	AggregationHooksKey              = []byte{SubModuleId, 6}
	RateLimitedHooksKey              = []byte{SubModuleId, 7}
	TokenRateLimitsKey               = []byte{SubModuleId, 8}
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

const RateLimitDurationSeconds uint64 = 86400

var RateLimitDuration = math.NewInt(int64(RateLimitDurationSeconds))

// TokenRateLimitKey creates and returns a new token rate limit key.
func TokenRateLimitKey(hookId, tokenId util.HexAddress) collections.Pair[uint64, []byte] {
	return collections.Join(hookId.GetInternalId(), tokenId.Bytes())
}
