package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AggregationHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = AggregationHookHandler{}

func (a AggregationHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := a.k.aggregationHooks.Has(ctx, hookId.GetInternalId())
	if err != nil {
		return false, err
	}
	return has, nil
}

func (a AggregationHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_AGGREGATION
}

func (a AggregationHookHandler) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	aggregationHook, err := a.k.aggregationHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return sdk.NewCoins(), err
	}

	chargedCoins := sdk.NewCoins()
	remainingCoins := maxFee
	for _, childHookId := range aggregationHook.Hooks {
		if childHookId.GetType() == uint32(types.POST_DISPATCH_HOOK_TYPE_AGGREGATION) {
			return sdk.NewCoins(), errors.Wrapf(types.ErrInvalidAggregationHook, "nested aggregation hook is not allowed: %s", childHookId.String())
		}

		handler, err := a.k.coreKeeper.PostDispatchRouter().GetModule(childHookId)
		if err != nil {
			return sdk.NewCoins(), err
		}

		childChargedCoins, err := (*handler).PostDispatch(ctx, mailboxId, childHookId, metadata, message, remainingCoins)
		if err != nil {
			return sdk.NewCoins(), err
		}

		remainingCoins, err = subtractCoins(remainingCoins, childChargedCoins)
		if err != nil {
			return sdk.NewCoins(), err
		}
		chargedCoins = chargedCoins.Add(childChargedCoins...)
	}

	return chargedCoins, nil
}

func (a AggregationHookHandler) QuoteDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage) (sdk.Coins, error) {
	aggregationHook, err := a.k.aggregationHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return sdk.NewCoins(), err
	}

	quotedCoins := sdk.NewCoins()
	for _, childHookId := range aggregationHook.Hooks {
		if childHookId.GetType() == uint32(types.POST_DISPATCH_HOOK_TYPE_AGGREGATION) {
			return sdk.NewCoins(), errors.Wrapf(types.ErrInvalidAggregationHook, "nested aggregation hook is not allowed: %s", childHookId.String())
		}

		handler, err := a.k.coreKeeper.PostDispatchRouter().GetModule(childHookId)
		if err != nil {
			return sdk.NewCoins(), err
		}

		childQuote, err := (*handler).QuoteDispatch(ctx, mailboxId, childHookId, metadata, message)
		if err != nil {
			return sdk.NewCoins(), err
		}
		quotedCoins = quotedCoins.Add(childQuote...)
	}

	return quotedCoins, nil
}

func subtractCoins(coins sdk.Coins, toSubtract sdk.Coins) (sdk.Coins, error) {
	remainingCoins, neg := coins.SafeSub(toSubtract...)
	if neg {
		return sdk.NewCoins(), fmt.Errorf("remaining coins cannot be negative")
	}
	return remainingCoins, nil
}
