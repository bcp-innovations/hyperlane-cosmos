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

func (k *Keeper) RemoteTransferSynthetic(ctx sdk.Context, token types.HypToken, cosmosSender string, destinationDomain uint32, recipient util.HexAddress, amount math.Int, customIgpId *util.HexAddress, gasLimit math.Int, maxFee sdk.Coin) (messageId util.HexAddress, err error) {
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
			Variant:  1,
			Value:    maxFee.Amount,
			GasLimit: gas,
			Address:  senderAcc,
		}.Bytes(),
		customIgpId,
	)
	if err != nil {
		return util.HexAddress{}, err
	}

	return dispatchMsg, nil
}

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
