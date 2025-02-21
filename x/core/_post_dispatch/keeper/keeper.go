package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	interchainGasPaymasters         collections.Map[uint64, types.InterchainGasPaymaster]
	interchainGasPaymastersSequence collections.Sequence
	schema                          collections.Schema

	hexAddressFactory util.HexAddressFactory
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	factory, err := util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER)
	if err != nil {
		panic(err)
	}

	k := Keeper{
		interchainGasPaymasters:         collections.NewMap(sb, types.PostDispatchHooksKey, "interchain_gas_paymasters", collections.Uint64Key, codec.CollValue[types.InterchainGasPaymaster](cdc)),
		interchainGasPaymastersSequence: collections.NewSequence(sb, types.PostDispatchHooksSequenceKey, "interchain_gas_paymasters_sequence"),
		hexAddressFactory:               factory,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k Keeper) SwitchHook(ctx sdk.Context, hookId util.HexAddress) (types.PostDispatchHook, error) {
	switch hookId.GetType() {
	case uint32(types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER):
		hook, err := k.interchainGasPaymasters.Get(ctx, hookId.GetInternalId())
		if err != nil {
			return nil, err
		}
		return InterchainGasPaymasterHook{hook, k}, nil
		// TODO add other cases
	}

	return nil, errors.New("invalid hook id")
}

func (k Keeper) PostDispatch(ctx sdk.Context, hookId util.HexAddress, metadata any, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	if !k.hexAddressFactory.IsClassMember(hookId) {
		return nil, nil
	}

	hook, err := k.SwitchHook(ctx, hookId)
	if err != nil {
		return nil, err
	}

	return hook.PostDispatch(ctx, metadata, message, maxFee)
}
