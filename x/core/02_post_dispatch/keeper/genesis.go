package keeper

import (
	"fmt"

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

	for _, aggregationHook := range data.AggregationHooks {
		if err := k.validateAggregationHooks(ctx, aggregationHook.Hooks); err != nil {
			panic(err)
		}
		if err := k.aggregationHooks.Set(ctx, aggregationHook.Id.GetInternalId(), aggregationHook); err != nil {
			panic(err)
		}
	}

	for _, hook := range data.RateLimitedHooks {
		if exists, err := k.coreKeeper.MailboxIdExists(ctx, hook.MailboxId); !exists || err != nil {
			panic(types.ErrMailboxDoesNotExist)
		}
		if err := k.rateLimitedHooks.Set(ctx, hook.Id.GetInternalId(), hook); err != nil {
			panic(err)
		}
	}

	for _, tokenRateLimit := range data.TokenRateLimits {
		if _, err := k.rateLimitedHooks.Get(ctx, tokenRateLimit.HookId.GetInternalId()); err != nil {
			panic(err)
		}
		if tokenRateLimit.MaxCapacity.LT(types.RateLimitDuration) {
			panic(fmt.Sprintf("max capacity must be at least %s", types.RateLimitDuration.String()))
		}
		if err := k.tokenRateLimits.Set(ctx, types.TokenRateLimitKey(tokenRateLimit.HookId, tokenRateLimit.TokenId), tokenRateLimit); err != nil {
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

	iterAggregationHooks, err := k.aggregationHooks.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	aggregationHooks, err := iterAggregationHooks.Values()
	if err != nil {
		panic(err)
	}

	iterRateLimitedHooks, err := k.rateLimitedHooks.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	rateLimitedHooks, err := iterRateLimitedHooks.Values()
	if err != nil {
		panic(err)
	}

	iterTokenRateLimits, err := k.tokenRateLimits.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	tokenRateLimits, err := iterTokenRateLimits.Values()
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Igps:             igps,
		IgpGasConfigs:    gasConfigs,
		MerkleTreeHooks:  merkleTreeHooks,
		NoopHooks:        noopHooks,
		AggregationHooks: aggregationHooks,
		RateLimitedHooks: rateLimitedHooks,
		TokenRateLimits:  tokenRateLimits,
	}
}
