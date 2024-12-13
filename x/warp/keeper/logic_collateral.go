package keeper

import (
	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RemoteTransferCollateral(ctx sdk.Context, token types.HypToken, cosmosSender string, externalRecipient string, amount math.Int) (messageId util.HexAddress, err error) {

	senderAcc, err := sdk.AccAddressFromBech32(cosmosSender)
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, types.ModuleName, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount)))
	if err != nil {
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
	)
	if err != nil {
		return util.HexAddress{}, err
	}

	return dispatchMsg, nil
}

func (k Keeper) RemoteReceiveCollateral(ctx sdk.Context, token types.HypToken, payload types.WarpPayload) error {

	account := sdk.AccAddress(payload.Recipient()[12:32])

	// TODO track balance for each token
	err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		account,
		sdk.NewCoins(sdk.NewCoin(token.OriginDenom, math.NewIntFromBigInt(payload.Amount()))),
	)
	if err != nil {
		return err
	}

	return nil
}
