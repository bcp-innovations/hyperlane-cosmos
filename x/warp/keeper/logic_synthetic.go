package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) RemoteTransferSynthetic(ctx sdk.Context, token types.HypToken, cosmosSender string, destinationDomain uint32, externalRecipient string, amount math.Int, customIgpId string, gasLimit math.Int, maxFee sdk.Coin) (messageId util.HexAddress, err error) {
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

	customPostDispatchHookId := util.NewZeroAddress()
	if customIgpId != "" {
		customPostDispatchHookId, err = util.DecodeHexAddress(customIgpId)
		if err != nil {
			return util.HexAddress{}, err
		}
	}

	// Token destinationDomain, recipientAddress
	dispatchMsg, err := k.mailboxKeeper.DispatchMessage(
		ctx,
		util.HexAddress(token.OriginMailbox),
		k.GetAddressFromToken(token),
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
		customPostDispatchHookId,
	)
	if err != nil {
		return util.HexAddress{}, err
	}

	return dispatchMsg, nil
}

func (k *Keeper) RemoteReceiveSynthetic(ctx sdk.Context, token types.HypToken, payload types.WarpPayload) error {
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

func ValidateTokenMetadata(metadata *types.TokenMetadata) error {
	if metadata == nil {
		return fmt.Errorf("metadata is required")
	}

	if metadata.Symbol == "" {
		return fmt.Errorf("token symbol is required")
	}

	if metadata.Name == "" {
		return fmt.Errorf("token name is required")
	}

	if metadata.Decimals == 0 {
		return fmt.Errorf("token decimals cannot be zero")
	}

	if metadata.TotalSupply.IsZero() {
		return fmt.Errorf("token total supply cannot be zero")
	}
	return nil
}
