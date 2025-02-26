package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PostDispatchHookHandler interface {
	HookType() uint8
	SupportsMetadata(metadata []byte) bool
	PostDispatch(ctx sdk.Context, hookId uint64, metadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
}

var (
	// TODO prefix with submodule
	PostDispatchHooksKey             = collections.NewPrefix(1)
	PostDispatchHooksSequenceKey     = collections.NewPrefix(2)
	InterchainGasPaymasterConfigsKey = collections.NewPrefix(3)
	MerkleTreeHooksKey               = collections.NewPrefix(4)
)

const (
	SubModuleName = "post_dispatch"

	HEX_ADDRESS_CLASS_IDENTIFIER = "corepostdispatch"
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
