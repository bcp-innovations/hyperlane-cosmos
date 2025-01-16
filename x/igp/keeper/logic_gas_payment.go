package keeper

import (
	"context"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) PayForGas(ctx context.Context, sender string, igpId util.HexAddress, messageId string, destinationDomain uint32, gasLimit uint64) error {
	igp, err := k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return err
	}

	requiredPayment, err := k.QuoteGasPayment(ctx, igpId, destinationDomain, gasLimit)

	senderAcc, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	beneficiaryAcc, err := sdk.AccAddressFromBech32(igp.Beneficiary)
	if err != nil {
		return err
	}

	// TODO: Get denom from global params
	balance := k.bankKeeper.GetBalance(ctx, senderAcc, types.Denom).Amount.Uint64()

	if balance < requiredPayment {
		return fmt.Errorf("insufficient balance to pay gas; got %v required %v", balance, requiredPayment)
	}

	// TODO: Get denom from global params
	coins := sdk.NewCoins(sdk.NewInt64Coin(types.Denom, int64(requiredPayment)))

	err = k.bankKeeper.SendCoins(ctx, senderAcc, beneficiaryAcc, coins)
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
	igp, err := k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return 0, err
	}

	var destinationGasConfig *types.DestinationGasConfig

	for _, c := range igp.DestinationGasConfigs {
		if c.RemoteDomain == destinationDomain {
			destinationGasConfig = c
		}
	}

	destinationCost := gasLimit * destinationGasConfig.GasOracle.GasPrice

	return (destinationCost * destinationGasConfig.GasOracle.TokenExchangeRate) / types.TokenExchangeRateScale, nil
}
