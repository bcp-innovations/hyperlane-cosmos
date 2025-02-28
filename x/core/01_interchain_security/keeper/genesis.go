package keeper

import (
	"cosmossdk.io/collections"
	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

func InitGenesis(ctx sdk.Context, k Keeper, data *types.GenesisState) {
	if data == nil || data.Isms == nil {
		return
	}

	isms, err := util.UnpackAnys[types.HyperlaneInterchainSecurityModule](data.Isms)
	if err != nil {
		panic(err)
	}

	for _, ism := range isms {
		id, err := ism.GetId()
		if err != nil {
			panic(err)
		}
		if err := k.isms.Set(ctx, id.GetInternalId(), ism); err != nil {
			panic(err)
		}
	}

	for _, storageLocation := range data.ValidatorStorageLocations {
		validatorBytes, err := util.DecodeEthHex(storageLocation.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		if err = k.storageLocations.Set(ctx, collections.Join3(storageLocation.MailboxId.Bytes(), validatorBytes, storageLocation.Index), storageLocation.StorageLocation); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	iter, err := k.isms.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	isms, err := iter.Values()
	if err != nil {
		panic(err)
	}

	msgs := make([]proto.Message, len(isms))
	for i, ism := range isms {
		msgs[i] = ism
	}
	ismsAny, err := util.PackAnys(msgs)
	if err != nil {
		panic(err)
	}

	// TODO add export for storage locations

	return &types.GenesisState{
		Isms: ismsAny,
	}
}
