package keeper

import (
	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) RemoteTransferCollateral(ctx sdk.Context, token types.HypToken, cosmosSender string, externalRecipient string, amount math.Int, customIgpId string, gasLimit uint64, maxFee uint64) (messageId util.HexAddress, err error) {

	senderAcc, err := sdk.AccAddressFromBech32(cosmosSender)
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, types.ModuleName, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)))
	if err != nil {
		return util.HexAddress{}, err
	}

	token.CollateralBalance = token.CollateralBalance.Add(amount)

	if err = k.HypTokens.Set(ctx, token.Id, token); err != nil {
		return util.HexAddress{}, err
	}

	recipient, err := util.DecodeEthHex(externalRecipient)
	if err != nil {
		return util.HexAddress{}, err
	}

	warpPayload, err := types.NewWarpPayload(recipient, *amount.BigInt())
	if err != nil {
		return util.HexAddress{}, err
	}

	// Token destinationDomain, recipientAddress
	dispatchMsg, err := k.mailboxKeeper.DispatchMessage(
		ctx,
		util.HexAddress(token.OriginMailbox),
		token.ReceiverDomain,
		util.HexAddress(token.ReceiverContract),
		util.HexAddress(token.Id),
		warpPayload.Bytes(),
		cosmosSender,
		customIgpId,
		gasLimit,
		maxFee,
	)
	if err != nil {
		return util.HexAddress{}, err
	}

	return dispatchMsg, nil
}

func (k *Keeper) RemoteReceiveCollateral(ctx sdk.Context, token types.HypToken, payload types.WarpPayload) error {

	account := sdk.AccAddress(payload.Recipient()[12:32])

	amount := math.NewIntFromBigInt(payload.Amount())

	token.CollateralBalance = token.CollateralBalance.Sub(amount)
	if token.CollateralBalance.IsNegative() {
		return types.ErrNotEnoughCollateral
	}

	if err := k.HypTokens.Set(ctx, token.Id, token); err != nil {
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
