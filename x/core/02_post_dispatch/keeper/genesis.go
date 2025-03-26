package keeper

import (
	"cosmossdk.io/collections"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data *types.GenesisState) {
	if data == nil || data.Igps == nil || data.IgpGasConfigs == nil ||
		data.MerkleTreeHooks == nil || data.NoopHooks == nil {
		panic("cannot init genesis state, some data not available")
	}

	for _, igp := range data.Igps {
		if err := k.Igps.Set(ctx, igp.Id.GetInternalId(), igp); err != nil {
			panic(err)
		}
	}

	for _, destinationGasConfig := range data.IgpGasConfigs {
		cfg := types.DestinationGasConfig{
			RemoteDomain: destinationGasConfig.RemoteDomain,
			GasOracle:    destinationGasConfig.GasOracle,
			GasOverhead:  destinationGasConfig.GasOverhead,
		}
		key := collections.Join(destinationGasConfig.IgpId, destinationGasConfig.RemoteDomain)
		if err := k.IgpDestinationGasConfigs.Set(ctx, key, cfg); err != nil {
			panic(err)
		}
	}

	for _, merkleTreeHook := range data.MerkleTreeHooks {
		if err := k.merkleTreeHooks.Set(ctx, merkleTreeHook.Id.GetInternalId(), merkleTreeHook); err != nil {
			panic(err)
		}
	}

	for _, noopHook := range data.NoopHooks {
		if err := k.noopHooks.Set(ctx, noopHook.Id.GetInternalId(), noopHook); err != nil {
			panic(err)
		}
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	iterIgp, err := k.Igps.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	igps, err := iterIgp.Values()
	if err != nil {
		panic(err)
	}

	iterConfigs, err := k.IgpDestinationGasConfigs.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	destinationGasConfigs, err := iterConfigs.KeyValues()
	if err != nil {
		panic(err)
	}

	gasConfigs := make([]types.GenesisDestinationGasConfigWrapper, len(destinationGasConfigs))
	for i := range destinationGasConfigs {
		cfg := types.GenesisDestinationGasConfigWrapper{
			RemoteDomain: destinationGasConfigs[i].Value.RemoteDomain,
			GasOracle:    destinationGasConfigs[i].Value.GasOracle,
			GasOverhead:  destinationGasConfigs[i].Value.GasOverhead,
			IgpId:        destinationGasConfigs[i].Key.K1(),
		}
		gasConfigs[i] = cfg
	}

	iterMerkleTreeHooks, err := k.merkleTreeHooks.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	merkleTreeHooks, err := iterMerkleTreeHooks.Values()
	if err != nil {
		panic(err)
	}

	iterNoopHooks, err := k.noopHooks.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	noopHooks, err := iterNoopHooks.Values()
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Igps:            igps,
		IgpGasConfigs:   gasConfigs,
		MerkleTreeHooks: merkleTreeHooks,
		NoopHooks:       noopHooks,
	}
}
