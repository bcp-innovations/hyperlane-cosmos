package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InterchainGasPaymasterHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = InterchainGasPaymasterHookHandler{}

func (i InterchainGasPaymasterHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER
}

func (i InterchainGasPaymasterHookHandler) SupportsMetadata(metadata []byte) bool {
	// TODO implement me
	panic("implement me")
}

func (i InterchainGasPaymasterHookHandler) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, rawMetadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	metadata, err := util.ParseStandardHookMetadata(rawMetadata)
	if err != nil {
		return nil, err
	}

	return i.PayForGas(ctx, hookId, metadata.Address.String(), message.Id().String(), message.Destination, metadata.GasLimit, maxFee)
}

func (i InterchainGasPaymasterHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := i.k.Igps.Has(ctx, hookId.GetInternalId())
	if err != nil {
		return false, err
	}
	return has, nil
}

// PayForGas executes an InterchainGasPayment using `QuoteGasPayment` beforehand.
// It uses the IGP's  DestinationGasConfig to determine the required payment and
// returns the remaining fees to pay gas for following hooks.
func (i InterchainGasPaymasterHookHandler) PayForGas(ctx context.Context, hookId util.HexAddress, sender string, messageId string, destinationDomain uint32, gasLimit math.Int, maxFee sdk.Coins) (sdk.Coins, error) {
	if maxFee.Empty() {
		return sdk.NewCoins(), fmt.Errorf("maxFee is required")
	}

	requiredPayment, err := i.QuoteGasPayment(ctx, hookId, destinationDomain, gasLimit)
	if err != nil {
		return sdk.NewCoins(), err
	}

	if requiredPayment.IsAllGT(maxFee) {
		return sdk.NewCoins(), fmt.Errorf("required payment exceeds max hyperlane fee: %v", requiredPayment)
	}

	return requiredPayment, i.PayForGasWithoutQuote(ctx, hookId, sender, messageId, destinationDomain, gasLimit, requiredPayment)
}

// PayForGasWithoutQuote executes an InterchainGasPayment without using `QuoteGasPayment`.
// This is used in the `MsgPayForGas` transaction, as the main purpose is paying an exact
// amount for e.g. re-funding a certain message-id as the first payment wasn't enough.
func (i InterchainGasPaymasterHookHandler) PayForGasWithoutQuote(ctx context.Context, hookId util.HexAddress, sender string, messageId string, destinationDomain uint32, gasLimit math.Int, amount sdk.Coins) error {
	igp, err := i.k.Igps.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return fmt.Errorf("igp does not exist: %s", hookId.String())
	}

	if amount.IsZero() {
		return fmt.Errorf("amount must be greater than zero")
	}

	if messageId == "" {
		return fmt.Errorf("message id cannot be empty")
	}

	senderAcc, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	// TODO use core-types module name or create sub-account
	err = i.k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, "hyperlane", amount)
	if err != nil {
		return err
	}

	igp.ClaimableFees = igp.ClaimableFees.Add(amount...)

	err = i.k.Igps.Set(ctx, igp.InternalId, igp)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.GasPayment{
		MessageId:   messageId,
		Destination: destinationDomain,
		GasAmount:   gasLimit.String(),
		Payment:     amount.String(),
		IgpId:       hookId.String(),
	})

	return nil
}

func (i InterchainGasPaymasterHookHandler) QuoteDispatch(ctx context.Context, _, hookId util.HexAddress, rawMetadata []byte, message util.HyperlaneMessage) (sdk.Coins, error) {
	metadata, err := util.ParseStandardHookMetadata(rawMetadata)
	if err != nil {
		return sdk.NewCoins(), err
	}

	return i.QuoteGasPayment(ctx, hookId, message.Destination, metadata.GasLimit)
}

func (i InterchainGasPaymasterHookHandler) QuoteGasPayment(ctx context.Context, hookId util.HexAddress, destinationDomain uint32, gasLimit math.Int) (sdk.Coins, error) {
	igp, err := i.k.Igps.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return sdk.NewCoins(), fmt.Errorf("igp does not exist: %s", hookId.String())
	}

	destinationGasConfig, err := i.k.IgpDestinationGasConfigs.Get(ctx, collections.Join(igp.InternalId, destinationDomain))
	if err != nil {
		return sdk.NewCoins(), fmt.Errorf("remote domain %v is not supported", destinationDomain)
	}

	gasLimit = gasLimit.Add(destinationGasConfig.GasOverhead)

	destinationCost := gasLimit.Mul(destinationGasConfig.GasOracle.GasPrice)

	amount := (destinationCost.Mul(destinationGasConfig.GasOracle.TokenExchangeRate)).Quo(types.TokenExchangeRateScale)

	coin := sdk.Coin{
		Denom:  igp.Denom,
		Amount: amount,
	}

	if err = coin.Validate(); err != nil {
		return sdk.NewCoins(), err
	}

	return sdk.NewCoins(coin), nil
}
