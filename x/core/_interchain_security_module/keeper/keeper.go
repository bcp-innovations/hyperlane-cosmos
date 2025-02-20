package keeper

import (
	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security_module/types"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	isms         collections.Map[uint64, types.HyperlaneInterchainSecurityModule]
	ismsSequence collections.Sequence
	schema       collections.Schema

	hexAddressFactory util.HexAddressFactory
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		isms:              collections.NewMap(sb, types.IsmsKey, "isms", collections.Uint64Key, codec.CollInterfaceValue[types.HyperlaneInterchainSecurityModule](cdc)),
		ismsSequence:      collections.NewSequence(sb, types.IsmsSequenceKey, "isms_sequence"),
		hexAddressFactory: util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k Keeper) Verify(ctx sdk.Context, ismId util.HexAddress, message util.HyperlaneMessage, metadata coreTypes.Metadata) (bool, error) {
	// Global Conventions
	// - Address must be unique
	// - Hook must check if id exists (and correct recipient)
	// module_name / class / type / custom

	if !k.hexAddressFactory.IsMember(ismId) {
		return false, nil
	}

	ism, err := k.isms.Get(ctx, 0)
	if err != nil {
		return false, err
	}

	return ism.Verify(ctx, metadata, message)
}
