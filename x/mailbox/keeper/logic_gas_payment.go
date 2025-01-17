package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Claim(ctx context.Context, sender string, igpId util.HexAddress) error {
	igp, err := k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return err
	}

	if sender != igp.Owner {
		return fmt.Errorf("failed to claim: %s is not permitted to claim", sender)
	}

	ownerAcc, err := sdk.AccAddressFromBech32(igp.Owner)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(igp.Denom, int64(igp.ClaimableFees)))

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, ownerAcc, coins)
	if err != nil {
		return err
	}

	igp.ClaimableFees = 0

	err = k.Igp.Set(ctx, igpId.Bytes(), igp)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) PayForGas(ctx context.Context, sender string, igpId util.HexAddress, messageId string, destinationDomain uint32, gasLimit uint64, maxFee uint64) error {
	igp, err := k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return err
	}

	requiredPayment, err := k.QuoteGasPayment(ctx, igpId, destinationDomain, gasLimit)
	if err != nil {
		return err
	}

	if requiredPayment > maxFee {
		return fmt.Errorf("required payment exceeds max hyperlane fee: %v", requiredPayment)
	}

	senderAcc, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	coins := sdk.NewCoins(sdk.NewInt64Coin(igp.Denom, int64(requiredPayment)))

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, types.ModuleName, coins)
	if err != nil {
		return err
	}

	igp.ClaimableFees += requiredPayment

	err = k.Igp.Set(ctx, igpId.Bytes(), igp)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.GasPayment{
		MessageId:   messageId,
		Destination: destinationDomain,
		GasAmount:   gasLimit,
		Payment:     requiredPayment,
		IgpId:       igpId.String(),
	})

	return nil
}

func (k Keeper) QuoteGasPayment(ctx context.Context, igpId util.HexAddress, destinationDomain uint32, gasLimit uint64) (uint64, error) {
	destinationGasConfig, err := k.IgpDestinationGasConfigMap.Get(ctx, collections.Join(igpId.Bytes(), destinationDomain))
	if err != nil {
		return 0, fmt.Errorf("remote domain %v is not supported: %e", destinationDomain, err)
	}

	destinationCost := gasLimit * destinationGasConfig.GasOracle.GasPrice

	return (destinationCost * destinationGasConfig.GasOracle.TokenExchangeRate) / types.TokenExchangeRateScale, nil
}
