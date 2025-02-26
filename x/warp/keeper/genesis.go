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

	// TODO:
	// for _, token := range data.Tokens {
	// 	if err := k.HypTokens.Set(ctx, token.Id, token); err != nil {
	// 		return err
	// 	}
	// }

	// if err := k.HypTokensCount.Set(ctx, uint64(len(data.Tokens))); err != nil {
	// 	return err
	// }

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var tokens []types.HypToken
	// TODO:
	// if err := k.HypTokens.Walk(ctx, nil, func(key uint64, value types.HypToken) (stop bool, err error) {
	// 	tokens = append(tokens, value)

	// 	return false, nil
	// }); err != nil {
	// 	return nil, err
	// }

	return &types.GenesisState{
		Params: params,
		Tokens: tokens,
	}, nil
}
