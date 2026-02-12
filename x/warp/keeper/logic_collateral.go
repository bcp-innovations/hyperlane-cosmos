package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RemoteTransferCollateral handles the transfer of collateral token to a remote chain.
// It withdraws the collateral from the sender, updates the token balance, and dispatches a message to the destination.
func (k *Keeper) RemoteTransferCollateral(ctx sdk.Context, token types.HypToken, cosmosSender string, destinationDomain uint32, externalRecipient util.HexAddress, amount math.Int, customHookId *util.HexAddress, gasLimit math.Int, maxFee sdk.Coin, customHookMetadata []byte) (messageId util.HexAddress, err error) {
	senderAcc, err := sdk.AccAddressFromBech32(cosmosSender)
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, types.ModuleName, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)))
	if err != nil {
		return util.HexAddress{}, err
	}

	token.CollateralBalance = token.CollateralBalance.Add(amount)

	if err = k.HypTokens.Set(ctx, token.Id.GetInternalId(), token); err != nil {
		return util.HexAddress{}, err
	}

	remoteRouter, err := k.EnrolledRouters.Get(ctx, collections.Join(token.Id.GetInternalId(), destinationDomain))
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("no enrolled router found for destination domain %d", destinationDomain)
	}

	gas := remoteRouter.Gas
	if !gasLimit.IsZero() {
		gas = gasLimit
	}

	warpPayload, err := types.NewWarpPayload(externalRecipient.Bytes(), *amount.BigInt())
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventSendRemoteTransfer{
		Sender:            cosmosSender,
		TokenId:           token.Id,
		DestinationDomain: destinationDomain,
		Recipient:         externalRecipient,
		Amount:            sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)).String(),
	})

	// Token destinationDomain, recipientAddress
	dispatchMsg, err := k.coreKeeper.DispatchMessage(
		ctx,
		util.HexAddress(token.OriginMailbox),
		token.Id, // sender
		sdk.NewCoins(maxFee),

		remoteRouter.ReceiverDomain,
		remoteRouter.ReceiverContract,

		warpPayload.Bytes(), // message body

		util.StandardHookMetadata{
			GasLimit:           gas,
			Address:            senderAcc,
			CustomHookMetadata: customHookMetadata,
		},
		customHookId,
	)
	if err != nil {
		return util.HexAddress{}, err
	}

	return dispatchMsg, nil
}

// RemoteReceiveCollateral handles the receipt of collateral from a remote chain.
// It updates the token balance, verifies there is enough collateral, and transfers the funds to the recipient's account.
func (k *Keeper) RemoteReceiveCollateral(ctx context.Context, token types.HypToken, payload types.WarpPayload) error {
	account := payload.GetCosmosAccount()

	amount := math.NewIntFromBigInt(payload.Amount())

	token.CollateralBalance = token.CollateralBalance.Sub(amount)
	if token.CollateralBalance.IsNegative() {
		return types.ErrNotEnoughCollateral
	}

	if err := k.HypTokens.Set(ctx, token.Id.GetInternalId(), token); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		account,
		sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)),
	); err != nil {
		return err
	}

	return nil
}
