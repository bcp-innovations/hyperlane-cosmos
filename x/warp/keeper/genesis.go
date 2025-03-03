package keeper

import (
	"context"

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

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	token, err := k.HypTokens.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	tokens, err := token.Values()
	if err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Params: params,
		Tokens: tokens,
	}, nil
}
