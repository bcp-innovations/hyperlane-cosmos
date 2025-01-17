package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management
	Igp                        collections.Map[[]byte, types.Igp]
	IgpDestinationGasConfigMap collections.Map[collections.Pair[[]byte, uint32], types.DestinationGasConfig]
	IgpSequence                collections.Sequence
	Params                     collections.Item[types.Params]
	Schema                     collections.Schema

	bankKeeper bankkeeper.Keeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService storetypes.KVStoreService,
	authority string,
	bankKeeper bankkeeper.Keeper,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		Igp:                        collections.NewMap(sb, types.IgpKey, "igp", collections.BytesKey, codec.CollValue[types.Igp](cdc)),
		IgpDestinationGasConfigMap: collections.NewMap(sb, types.IgpDestinationGasConfigMapKey, "igp_destination_gas_config_map", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), codec.CollValue[types.DestinationGasConfig](cdc)),
		IgpSequence:                collections.NewSequence(sb, types.IgpSequenceKey, "igp_sequence"),
		Params:                     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),

		bankKeeper: bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) IgpIdExists(ctx context.Context, igpId util.HexAddress) (bool, error) {
	igp, err := k.Igp.Has(ctx, igpId.Bytes())
	if err != nil {
		return false, err
	}
	return igp, nil
}
