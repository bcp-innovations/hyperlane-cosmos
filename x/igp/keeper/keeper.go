package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
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
	Igp         collections.Map[[]byte, types.Igp]
	IgpSequence collections.Sequence
	Params      collections.Item[types.Params]
	Schema      collections.Schema

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

		Igp:         collections.NewMap(sb, types.IgpKey, "igp", collections.BytesKey, codec.CollValue[types.Igp](cdc)),
		IgpSequence: collections.NewSequence(sb, types.IgpSequenceKey, "igp_sequence"),
		Params:      collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),

		bankKeeper: bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}
