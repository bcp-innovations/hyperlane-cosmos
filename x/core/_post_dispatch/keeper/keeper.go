package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	igps                     collections.Map[uint64, types.InterchainGasPaymaster]
	igpDestinationGasConfigs collections.Map[collections.Pair[uint64, uint32], types.DestinationGasConfig]

	merkleTreeHooks collections.Map[uint64, types.MerkleTreeHook]

	schema    collections.Schema
	idFactory PostDispatchHookFactory

	bankKeeper types.BankKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService, bankKeeper types.BankKeeper) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		igps:                     collections.NewMap(sb, types.PostDispatchHooksKey, "interchain_gas_paymasters", collections.Uint64Key, codec.CollValue[types.InterchainGasPaymaster](cdc)),
		igpDestinationGasConfigs: collections.NewMap(sb, types.InterchainGasPaymasterConfigsKey, "interchain_gas_paymaster_configs", collections.PairKeyCodec(collections.Uint64Key, collections.Uint32Key), codec.CollValue[types.DestinationGasConfig](cdc)),

		merkleTreeHooks: collections.NewMap(sb, types.MerkleTreeHooksKey, "merkle_tree_hooks_key", collections.Uint64Key, codec.CollValue[types.MerkleTreeHook](cdc)),

		bankKeeper: bankKeeper,
	}

	k.idFactory = NewPostDispatchHookFactory(collections.NewSequence(sb, types.PostDispatchHooksSequenceKey, "post_dispatch_hooks_sequence_key"))

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	// Register hooks
	k.idFactory.RegisterHook(types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER, InterchainGasPaymasterHookHandler{k})
	k.idFactory.RegisterHook(types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE, MerkleTreeHookHandler{k})

	return k
}

type PostDispatchHookFactory struct {
	hookMap           map[uint8]types.PostDispatchHookHandler
	HexAddressFactory util.HexAddressFactory
	sequence          collections.Sequence
}

func NewPostDispatchHookFactory(sequence collections.Sequence) PostDispatchHookFactory {
	factory, err := util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER)
	if err != nil {
		panic(err)
	}

	return PostDispatchHookFactory{
		hookMap:           map[uint8]types.PostDispatchHookHandler{},
		HexAddressFactory: factory,
		sequence:          sequence,
	}
}

func (p PostDispatchHookFactory) RegisterHook(hookType uint8, hook types.PostDispatchHookHandler) {
	p.hookMap[hookType] = hook
}

func (p PostDispatchHookFactory) GenerateNewId(ctx context.Context) uint64 {
	next, _ := p.sequence.Next(ctx)
	return next
}

func (p PostDispatchHookFactory) AddressFromId(id uint64) util.HexAddress {
	// TODO what internal Type to use?
	return p.HexAddressFactory.GenerateId(0, id)
}

func (p PostDispatchHookFactory) GetHookHandler(hookType uint8) (types.PostDispatchHookHandler, error) {
	hook, ok := p.hookMap[hookType]
	if !ok {
		return nil, fmt.Errorf("hook type %d not registered", hookType)
	}
	return hook, nil
}

func (k Keeper) PostDispatch(ctx sdk.Context, hookId util.HexAddress, metadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	if !k.idFactory.HexAddressFactory.IsClassMember(hookId) {
		return nil, nil
	}

	// TODO possible overflow
	paymasterHook, err := k.idFactory.GetHookHandler(uint8(hookId.GetType()))
	if err != nil {
		return nil, err
	}

	// TODO figure out internal id
	return paymasterHook.PostDispatch(ctx, hookId.GetInternalId(), metadata, message, maxFee)
}
