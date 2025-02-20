package keeper

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct{}

func (k Keeper) PostDispatch(ctx sdk.Context, hookId uint64, message types.HyperlaneMessage, metadata types.Metadata) error {
	// TODO Call stuff and do stuff

	return nil
}

func (k Keeper) QuoteDispatch(ctx sdk.Context, hookId uint64, message types.HyperlaneMessage, metadata types.Metadata) (sdk.Coins, error) {
	// TODO Call stuff and do stuff

	return nil, nil
}
