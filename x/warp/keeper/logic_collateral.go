package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) RemoteTransferCollateral(ctx sdk.Context, token types.HypToken, cosmosSender string, destinationDomain uint32, externalRecipient string, amount math.Int, customIgpId string, gasLimit math.Int, maxFee math.Int) (messageId util.HexAddress, err error) {
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

	if externalRecipient == "" {
		return util.HexAddress{}, fmt.Errorf("recipient cannot be empty")
	}

	recipient, err := util.DecodeEthHex(externalRecipient)
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("invalid recipient address")
	}

	remoteRouter, err := k.EnrolledRouters.Get(ctx, collections.Join(token.Id, destinationDomain))
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

	warpPayload, err := types.NewWarpPayload(recipient, *amount.BigInt())
	if err != nil {
		return util.HexAddress{}, err
	}

	igpCustomHookId := util.NewZeroAddress()
	if customIgpId != "" {
		igpCustomHookId, err = util.DecodeHexAddress(customIgpId)
		if err != nil {
			return util.HexAddress{}, err
		}
	}

	// Token destinationDomain, recipientAddress
	dispatchMsg, err := k.mailboxKeeper.DispatchMessage(
		ctx,
		util.HexAddress(token.OriginMailbox),
		k.GetAddressFromToken(token), // sender
		sdk.Coins{},                  // TODO figure out Coins                    // maxFee the user is willing to pay

		remoteRouter.ReceiverDomain,
		receiverContract,

		warpPayload.Bytes(), // message body

		util.StandardHookMetadata{
			Variant:  1,
			Value:    *big.NewInt(0), // TODO figure out usage of maxFee
			GasLimit: *gas.BigInt(),
			Address:  senderAcc,
		}.Bytes(), // metadata for gas payment
		igpCustomHookId, // don't override post dispatch hook
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
