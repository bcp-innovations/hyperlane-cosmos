package keeper

import (
	"context"
	"errors"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	k Keeper
}

func (ms msgServer) CreateSyntheticToken(ctx context.Context, msg *types.MsgCreateSyntheticToken) (*types.MsgCreateSyntheticTokenResponse, error) {
	next, err := ms.k.HypTokensCount.Next(ctx)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeHexAddress(msg.OriginMailbox)
	if err != nil {
		return nil, err
	}

	has, err := ms.k.mailboxKeeper.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("mailbox not found")
	}

	receiverContract, err := util.DecodeHexAddress(msg.ReceiverContract)
	if err != nil {
		return nil, err
	}

	tokenId := util.CreateHexAddress(types.ModuleName, int64(next))

	ismId, err := util.DecodeHexAddress(msg.IsmId)
	if err != nil {
		return nil, err
	}
	err = ms.k.mailboxKeeper.RegisterReceiverIsm(ctx, tokenId, ismId)
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:               tokenId.Bytes(),
		Creator:          msg.Creator,
		TokenType:        types.HYP_TOKEN_SYNTHETIC,
		OriginMailbox:    mailboxId.Bytes(),
		OriginDenom:      fmt.Sprintf("hyperlane/%s", tokenId.String()),
		ReceiverDomain:   msg.ReceiverDomain,
		ReceiverContract: receiverContract.Bytes(),
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.Bytes(), newToken); err != nil {
		return nil, err
	}
	return &types.MsgCreateSyntheticTokenResponse{}, nil
}

func (ms msgServer) CreateCollateralToken(ctx context.Context, msg *types.MsgCreateCollateralToken) (*types.MsgCreateCollateralTokenResponse, error) {

	next, err := ms.k.HypTokensCount.Next(ctx)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeHexAddress(msg.OriginMailbox)
	if err != nil {
		return nil, err
	}

	has, err := ms.k.mailboxKeeper.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("mailbox not found")
	}

	receiverContract, err := util.DecodeHexAddress(msg.ReceiverContract)
	if err != nil {
		return nil, err
	}

	tokenId := util.CreateHexAddress(types.ModuleName, int64(next))

	ismId, err := util.DecodeHexAddress(msg.IsmId)
	if err != nil {
		return nil, err
	}
	err = ms.k.mailboxKeeper.RegisterReceiverIsm(ctx, tokenId, ismId)
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:               tokenId.Bytes(),
		Creator:          msg.Creator,
		TokenType:        types.HYP_TOKEN_COLLATERAL,
		OriginMailbox:    mailboxId.Bytes(),
		OriginDenom:      msg.OriginDenom,
		ReceiverDomain:   msg.ReceiverDomain,
		ReceiverContract: receiverContract.Bytes(),
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.Bytes(), newToken); err != nil {
		return nil, err
	}
	return &types.MsgCreateCollateralTokenResponse{}, nil
}

func (ms msgServer) RemoteTransfer(ctx context.Context, msg *types.MsgRemoteTransfer) (*types.MsgRemoteTransferResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, err
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.Bytes())
	if err != nil {
		return nil, err
	}

	var messageResultId string
	if token.TokenType == types.HYP_TOKEN_COLLATERAL {
		result, err := ms.k.RemoteTransferCollateral(goCtx, token, msg.Sender, msg.Recipient, msg.Amount)
		if err != nil {
			return nil, err
		}
		messageResultId = result.String()
	} else if token.TokenType == types.HYP_TOKEN_SYNTHETIC {
		result, err := ms.k.RemoteTransferSynthetic(goCtx, token, msg.Sender, msg.Recipient, msg.Amount)
		if err != nil {
			return nil, err
		}
		messageResultId = result.String()
	} else {
		return nil, errors.New("invalid token type")
	}

	return &types.MsgRemoteTransferResponse{
		MessageId: messageResultId,
	}, nil
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}
