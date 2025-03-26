package keeper

import (
	"context"

	"cosmossdk.io/collections"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	if err := k.Params.Set(ctx, data.Params); err != nil {
		return err
	}
	for _, token := range data.Tokens {
		if err := k.HypTokens.Set(ctx, token.Id.GetInternalId(), token); err != nil {
			return err
		}
	}

	for _, r := range data.RemoteRouters {
		if err := k.EnrolledRouters.Set(ctx, collections.Join(r.TokenId, r.RemoteRouter.ReceiverDomain), r.RemoteRouter); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	tokenIterator, err := k.HypTokens.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	tokens, err := tokenIterator.Values()
	if err != nil {
		return nil, err
	}

	genesisRouters := make([]types.GenesisRemoteRouterWrapper, 0)
	err = k.EnrolledRouters.Walk(ctx, nil, func(key collections.Pair[uint64, uint32], value types.RemoteRouter) (stop bool, err error) {
		genesisRouters = append(genesisRouters, types.GenesisRemoteRouterWrapper{
			TokenId:      key.K1(),
			RemoteRouter: value,
		})
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Tokens:        tokens,
		Params:        params,
		RemoteRouters: genesisRouters,
	}, nil
}
