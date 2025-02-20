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

	k := Keeper{
		interchainGasPaymasters:         collections.NewMap(sb, types.PostDispatchHooksKey, "interchain_gas_paymasters", collections.Uint64Key, codec.CollValue[types.InterchainGasPaymaster](cdc)),
		interchainGasPaymastersSequence: collections.NewSequence(sb, types.PostDispatchHooksSequenceKey, "interchain_gas_paymasters_sequence"),
		hexAddressFactory:               util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k Keeper) SwitchHook(ctx sdk.Context, hookId util.HexAddress) (types.PostDispatchHook, error) {
	switch hookId.GetClass() {
	case "123":
		hook, err := k.interchainGasPaymasters.Get(ctx, hookId.GetInternalId())
		if err != nil {
			return nil, err
		}
		return InterchainGasPaymasterHook{hook, k}, nil
		// TODO add other cases
	}

	return nil, errors.New("invalid hook id")
}

func (k Keeper) PostDispatch(ctx sdk.Context, hookId util.HexAddress, message util.HyperlaneMessage, metadata any) error {
	if !k.hexAddressFactory.IsClassMember(hookId) {
		return nil
	}

	hook, err := k.SwitchHook(ctx, hookId)
	if err != nil {
		return err
	}

	return hook.PostDispatch(ctx, metadata, message)
}

func (k Keeper) QuoteDispatch(ctx sdk.Context, hookId util.HexAddress, message util.HyperlaneMessage, metadata any) (sdk.Coins, error) {
	if !k.hexAddressFactory.IsClassMember(hookId) {
		return nil, nil
	}

	hook, err := k.SwitchHook(ctx, hookId)
	if err != nil {
		return nil, err
	}

	return hook.QuotePostDispatch(ctx, metadata, message)
}
