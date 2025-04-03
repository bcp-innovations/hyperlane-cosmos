package keeper

import (
	"fmt"

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

	isms := make([]types.HyperlaneInterchainSecurityModule, len(data.Isms))
	for i, rawIsm := range data.Isms {
		var item types.HyperlaneInterchainSecurityModule
		// TODO find generic solution with proto unpack any
		switch rawIsm.TypeUrl {
		case "/hyperlane.core.interchain_security.v1.NoopISM":
			item = &types.NoopISM{}
		case "/hyperlane.core.interchain_security.v1.MessageIdMultisigISM":
			item = &types.MessageIdMultisigISM{}
		case "/hyperlane.core.interchain_security.v1.MerkleRootMultisigISM":
			item = &types.MerkleRootMultisigISM{}
		case "/hyperlane.core.interchain_security.v1.RoutingISM":
			item = &types.RoutingISM{}
		default:
			panic(fmt.Sprintf("unsupported type %s", rawIsm.TypeUrl))
		}
		if err := proto.Unmarshal(rawIsm.Value, item); err != nil {
			panic(err)
		}
		isms[i] = item
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

		if err = k.storageLocations.Set(ctx, collections.Join3(storageLocation.MailboxId, validatorBytes, storageLocation.Index), storageLocation.StorageLocation); err != nil {
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

	iterStorageLocations, err := k.storageLocations.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	storageLocations, err := iterStorageLocations.KeyValues()
	if err != nil {
		panic(err)
	}

	wrappedLocations := make([]types.GenesisValidatorStorageLocationWrapper, len(storageLocations))
	for i := range storageLocations {
		location := types.GenesisValidatorStorageLocationWrapper{
			MailboxId:        storageLocations[i].Key.K1(),
			ValidatorAddress: util.EncodeEthHex(storageLocations[i].Key.K2()),
			Index:            storageLocations[i].Key.K3(),
			StorageLocation:  storageLocations[i].Value,
		}
		wrappedLocations[i] = location
	}

	return &types.GenesisState{
		Isms:                      ismsAny,
		ValidatorStorageLocations: wrappedLocations,
	}
}
