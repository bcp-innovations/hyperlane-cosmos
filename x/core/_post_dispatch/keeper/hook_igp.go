package keeper

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InterchainGasPaymasterHook struct {
	igp types.InterchainGasPaymaster
	k   Keeper
}

var _ types.PostDispatchHook = InterchainGasPaymasterHook{}

func (i InterchainGasPaymasterHook) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER
}

func (i InterchainGasPaymasterHook) SupportsMetadata(metadata any) bool {
	//TODO implement me
	panic("implement me")
}

func (i InterchainGasPaymasterHook) PostDispatch(ctx sdk.Context, metadata any, message util.HyperlaneMessage) error {
	//TODO implement me
	panic("implement me")
}

func (i InterchainGasPaymasterHook) QuotePostDispatch(ctx sdk.Context, metadata any, message util.HyperlaneMessage) (sdk.Coins, error) {
	//TODO implement me
	panic("implement me")
}
