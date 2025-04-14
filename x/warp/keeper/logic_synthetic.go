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

// RemoteTransferSynthetic handles the transfer of synthetic tokens to a remote chain.
// It burns the specified amount of tokens, then dispatches a message to the destination
// with the required gas limit, fee, and custom metadata.
func (k *Keeper) RemoteTransferSynthetic(ctx sdk.Context, token types.HypToken, cosmosSender string, destinationDomain uint32, recipient util.HexAddress, amount math.Int, customHookId *util.HexAddress, gasLimit math.Int, maxFee sdk.Coin, customHookMetadata []byte) (messageId util.HexAddress, err error) {
	senderAcc, err := sdk.AccAddressFromBech32(cosmosSender)
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, types.ModuleName, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)))
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)))
	if err != nil {
		return util.HexAddress{}, err
	}

	remoteRouter, err := k.EnrolledRouters.Get(ctx, collections.Join(token.Id.GetInternalId(), destinationDomain))
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("no enrolled router found for destination domain %d", destinationDomain)
	}

	receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to decode receiver contract address %s", remoteRouter.ReceiverContract)
	}

	gas := remoteRouter.Gas
	if !gasLimit.IsZero() {
		gas = gasLimit
	}

	warpPayload, err := types.NewWarpPayload(recipient.Bytes(), *amount.BigInt())
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = ctx.EventManager().EmitTypedEvent(&types.EventSendRemoteTransfer{
		Sender:            cosmosSender,
		TokenId:           token.Id,
		DestinationDomain: destinationDomain,
		Recipient:         recipient,
		Amount:            sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)).String(),
	})

	// Token destinationDomain, recipientAddress
	dispatchMsg, err := k.coreKeeper.DispatchMessage(
		ctx,
		util.HexAddress(token.OriginMailbox),
		token.Id,
		sdk.NewCoins(maxFee),

		remoteRouter.ReceiverDomain,
		receiverContract,

		warpPayload.Bytes(),
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

// RemoteReceiveSynthetic handles the receipt of synthetic tokens from a remote chain.
// It mints the synthetic token and transfers it to the recipient's account.
func (k *Keeper) RemoteReceiveSynthetic(ctx context.Context, token types.HypToken, payload types.WarpPayload) error {
	account := payload.GetCosmosAccount()

	shadowToken := sdk.NewCoin(
		token.OriginDenom,
		math.NewIntFromBigInt(payload.Amount()),
	)

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(shadowToken)); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, account, sdk.NewCoins(shadowToken)); err != nil {
		return err
	}

	return nil
}
