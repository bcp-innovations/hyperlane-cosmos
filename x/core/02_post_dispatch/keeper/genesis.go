package keeper

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data *types.GenesisState) {
	if data == nil {
		return
	}

	// TODO init state
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	// TODO implement

	return &types.GenesisState{}
}
