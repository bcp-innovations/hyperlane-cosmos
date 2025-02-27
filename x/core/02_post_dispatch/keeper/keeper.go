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

	schema collections.Schema

	coreKeeper types.CoreKeeper
	bankKeeper types.BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService, bankKeeper types.BankKeeper) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		Igps:                     collections.NewMap(sb, types.PostDispatchHooksKey, "interchain_gas_paymasters", collections.Uint64Key, codec.CollValue[types.InterchainGasPaymaster](cdc)),
		IgpDestinationGasConfigs: collections.NewMap(sb, types.InterchainGasPaymasterConfigsKey, "interchain_gas_paymaster_configs", collections.PairKeyCodec(collections.Uint64Key, collections.Uint32Key), codec.CollValue[types.DestinationGasConfig](cdc)),

		merkleTreeHooks: collections.NewMap(sb, types.MerkleTreeHooksKey, "merkle_tree_hooks_key", collections.Uint64Key, codec.CollValue[types.MerkleTreeHook](cdc)),

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
}
