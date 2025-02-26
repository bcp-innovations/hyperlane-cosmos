package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InterchainGasPaymasterHookHandler struct {
	k Keeper
}

var _ types.PostDispatchHookHandler = InterchainGasPaymasterHookHandler{}

func (i InterchainGasPaymasterHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER
}

func (i InterchainGasPaymasterHookHandler) SupportsMetadata(metadata []byte) bool {
	// TODO implement me
	panic("implement me")
}

func (i InterchainGasPaymasterHookHandler) PostDispatch(ctx sdk.Context, hookId uint64, rawMetadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	metadata, err := util.ParseStandardHookMetadata(rawMetadata)
	if err != nil {
		return nil, err
	}

	err = i.PayForGas(ctx, hookId, metadata.Address.String(), message.Id().String(), message.Destination, math.NewIntFromBigInt(&metadata.GasLimit), math.NewIntFromBigInt(&metadata.Value))
	if err != nil {
		return nil, err
	}

	// TODO substract coins
	return nil, nil
}

// PayForGasWithoutQuote executes an InterchainGasPayment without using `QuoteGasPayment`.
// This is used in the `MsgPayForGas` transaction, as the main purpose is paying an exact
// amount for e.g. re-funding a certain message-id as the first payment wasn't enough.
func (i InterchainGasPaymasterHookHandler) PayForGasWithoutQuote(ctx context.Context, hookId uint64, sender string, messageId string, destinationDomain uint32, gasLimit math.Int, amount math.Int) error {
	igp, err := i.k.igps.Get(ctx, hookId)
	if err != nil {
		return fmt.Errorf("igp does not exist: %d", hookId)
	}

	if amount.Equal(math.ZeroInt()) {
		return fmt.Errorf("amount must be greater than zero")
	}

	if messageId == "" {
		return fmt.Errorf("message id cannot be empty")
	}

	senderAcc, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(igp.Denom, amount.Int64()))

	// TODO use core-types module name or create sub-account
	err = i.k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, "TODO", coins)
	if err != nil {
		return err
	}

	igp.ClaimableFees = igp.ClaimableFees.Add(amount)

	err = i.k.igps.Set(ctx, igp.Id, igp)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.GasPayment{
		MessageId:   messageId,
		Destination: destinationDomain,
		GasAmount:   gasLimit.String(),
		Payment:     amount.String(),
		IgpId:       "", // TODO figure out id usage
	})

	return nil
}

func (i InterchainGasPaymasterHookHandler) QuoteGasPayment(ctx context.Context, hookId uint64, destinationDomain uint32, gasLimit math.Int) (math.Int, error) {
	igp, err := i.k.igps.Get(ctx, hookId)
	if err != nil {
		return math.ZeroInt(), fmt.Errorf("igp does not exist: %d", hookId)
	}

	destinationGasConfig, err := i.k.igpDestinationGasConfigs.Get(ctx, collections.Join(igp.Id, destinationDomain))
	if err != nil {
		return math.Int{}, fmt.Errorf("remote domain %v is not supported", destinationDomain)
	}

	gasLimit = gasLimit.Add(destinationGasConfig.GasOverhead)

	destinationCost := gasLimit.Mul(destinationGasConfig.GasOracle.GasPrice)

	return (destinationCost.Mul(destinationGasConfig.GasOracle.TokenExchangeRate)).Quo(types.TokenExchangeRateScale), nil
}

func (i InterchainGasPaymasterHookHandler) PayForGas(ctx context.Context, hookId uint64, sender string, messageId string, destinationDomain uint32, gasLimit math.Int, maxFee math.Int) error {
	requiredPayment, err := i.QuoteGasPayment(ctx, hookId, destinationDomain, gasLimit)
	if err != nil {
		return err
	}

	if requiredPayment.GT(maxFee) {
		return fmt.Errorf("required payment exceeds max hyperlane fee: %v", requiredPayment)
	}

	return i.PayForGasWithoutQuote(ctx, hookId, sender, messageId, destinationDomain, gasLimit, requiredPayment)
}
