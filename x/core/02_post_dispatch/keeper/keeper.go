package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	Igps                     collections.Map[uint64, types.InterchainGasPaymaster]
	IgpDestinationGasConfigs collections.Map[collections.Pair[uint64, uint32], types.DestinationGasConfig]

	merkleTreeHooks collections.Map[uint64, types.MerkleTreeHook]

	noopHooks collections.Map[uint64, types.NoopHook]

	aggregationHooks collections.Map[uint64, types.AggregationHook]

	pausableHooks    collections.Map[uint64, types.PausableHook]
	rateLimitedHooks collections.Map[uint64, types.RateLimitedHook]
	tokenRateLimits  collections.Map[collections.Pair[uint64, []byte], types.TokenRateLimit]

	schema collections.Schema

	coreKeeper types.CoreKeeper
	bankKeeper types.BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService, bankKeeper types.BankKeeper) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		Igps:                     collections.NewMap(sb, types.InterchainGasPaymasterKey, "interchain_gas_paymasters", collections.Uint64Key, codec.CollValue[types.InterchainGasPaymaster](cdc)),
		IgpDestinationGasConfigs: collections.NewMap(sb, types.InterchainGasPaymasterConfigsKey, "interchain_gas_paymaster_configs", collections.PairKeyCodec(collections.Uint64Key, collections.Uint32Key), codec.CollValue[types.DestinationGasConfig](cdc)),

		merkleTreeHooks: collections.NewMap(sb, types.MerkleTreeHooksKey, "merkle_tree_hooks_key", collections.Uint64Key, codec.CollValue[types.MerkleTreeHook](cdc)),
		noopHooks:       collections.NewMap(sb, types.NoopHooksKey, "noop_hooks_key", collections.Uint64Key, codec.CollValue[types.NoopHook](cdc)),

		aggregationHooks: collections.NewMap(sb, types.AggregationHooksKey, "aggregation_hooks_key", collections.Uint64Key, codec.CollValue[types.AggregationHook](cdc)),

		pausableHooks:    collections.NewMap(sb, types.PausableHooksKey, "pausable_hooks_key", collections.Uint64Key, codec.CollValue[types.PausableHook](cdc)),
		rateLimitedHooks: collections.NewMap(sb, types.RateLimitedHooksKey, "rate_limited_hooks_key", collections.Uint64Key, codec.CollValue[types.RateLimitedHook](cdc)),
		tokenRateLimits:  collections.NewMap(sb, types.TokenRateLimitsKey, "token_rate_limits_key", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey), codec.CollValue[types.TokenRateLimit](cdc)),

		bankKeeper: bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k *Keeper) SetCoreKeeper(coreKeeper types.CoreKeeper) {
	if k.coreKeeper != nil {
		panic("core keeper already set")
	}

	k.coreKeeper = coreKeeper

	router := coreKeeper.PostDispatchRouter()
	// add default post dispatch hooks
	router.RegisterModule(types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE, MerkleTreeHookHandler{*k})
	router.RegisterModule(types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER, InterchainGasPaymasterHookHandler{*k})
	router.RegisterModule(types.POST_DISPATCH_HOOK_TYPE_UNUSED, NoopHookHandler{*k})
	router.RegisterModule(types.POST_DISPATCH_HOOK_TYPE_AGGREGATION, AggregationHookHandler{*k})
	router.RegisterModule(types.POST_DISPATCH_HOOK_TYPE_PAUSABLE, PausableHookHandler{*k})
	router.RegisterModule(types.POST_DISPATCH_HOOK_TYPE_RATE_LIMITED, RateLimitedHookHandler{*k})
}
