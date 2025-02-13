package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"

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

	tokenId := util.CreateHexAddress(types.ModuleName, int64(next))

	err = ms.k.mailboxKeeper.RegisterReceiverIsm(ctx, tokenId, &mailboxId, msg.IsmId)
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:            tokenId.Bytes(),
		Creator:       msg.Creator,
		TokenType:     types.HYP_TOKEN_TYPE_SYNTHETIC,
		OriginMailbox: mailboxId.Bytes(),
		OriginDenom:   fmt.Sprintf("hyperlane/%s", tokenId.String()),
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.Bytes(), newToken); err != nil {
		return nil, err
	}
	return &types.MsgCreateSyntheticTokenResponse{Id: tokenId.String()}, nil
}

// CreateCollateralToken ...
// TODO: setGasRouter tx
func (ms msgServer) CreateCollateralToken(ctx context.Context, msg *types.MsgCreateCollateralToken) (*types.MsgCreateCollateralTokenResponse, error) {
	next, err := ms.k.HypTokensCount.Next(ctx)
	if err != nil {
		return nil, err
	}

	err = sdk.ValidateDenom(msg.OriginDenom)
	if err != nil {
		return nil, fmt.Errorf("origin denom %s is invalid", msg.OriginDenom)
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

	tokenId := util.CreateHexAddress(types.ModuleName, int64(next))

	err = ms.k.mailboxKeeper.RegisterReceiverIsm(ctx, tokenId, &mailboxId, msg.IsmId)
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:            tokenId.Bytes(),
		Creator:       msg.Creator,
		TokenType:     types.HYP_TOKEN_TYPE_COLLATERAL,
		OriginMailbox: mailboxId.Bytes(),
		OriginDenom:   msg.OriginDenom,
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.Bytes(), newToken); err != nil {
		return nil, err
	}
	return &types.MsgCreateCollateralTokenResponse{Id: tokenId.String()}, nil
}

func (ms msgServer) EnrollRemoteRouter(ctx context.Context, msg *types.MsgEnrollRemoteRouter) (*types.MsgEnrollRemoteRouterResponse, error) {
	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("token with id %s not found", tokenId.String())
	}

	if token.Creator != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	if msg.RemoteRouter == nil {
		return nil, fmt.Errorf("invalid remote router")
	}

	if err = ms.k.EnrolledRouters.Set(ctx, collections.Join(tokenId.Bytes(), msg.RemoteRouter.ReceiverDomain), *msg.RemoteRouter); err != nil {
		return nil, err
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.Bytes(), token); err != nil {
		return nil, err
	}

	return &types.MsgEnrollRemoteRouterResponse{}, nil
}

func (ms msgServer) UnrollRemoteRouter(ctx context.Context, msg *types.MsgUnrollRemoteRouter) (*types.MsgUnrollRemoteRouterResponse, error) {
	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("token with id %s not found", tokenId.String())
	}

	if token.Creator != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	exists, err := ms.k.EnrolledRouters.Has(ctx, collections.Join(tokenId.Bytes(), msg.ReceiverDomain))
	if err != nil || !exists {
		return nil, fmt.Errorf("failed to find remote router for domain %v", msg.ReceiverDomain)
	}

	if err = ms.k.EnrolledRouters.Remove(ctx, collections.Join(tokenId.Bytes(), msg.ReceiverDomain)); err != nil {
		return nil, err
	}

	return &types.MsgUnrollRemoteRouterResponse{}, nil
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
	if token.TokenType == types.HYP_TOKEN_TYPE_COLLATERAL {
		result, err := ms.k.RemoteTransferCollateral(goCtx, token, msg.Sender, msg.DestinationDomain, msg.Recipient, msg.Amount, msg.IgpId, msg.GasLimit, msg.MaxFee)
		if err != nil {
			return nil, err
		}
		messageResultId = result.String()
	} else if token.TokenType == types.HYP_TOKEN_TYPE_SYNTHETIC {
		result, err := ms.k.RemoteTransferSynthetic(goCtx, token, msg.Sender, msg.DestinationDomain, msg.Recipient, msg.Amount, msg.IgpId, msg.GasLimit, msg.MaxFee)
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

func (ms msgServer) SetIsm(ctx context.Context, msg *types.MsgSetIsm) (*types.MsgSetIsmResponse, error) {
	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, err
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.Bytes())
	if err != nil {
		return nil, err
	}

	if token.Creator != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	err = ms.k.mailboxKeeper.RegisterReceiverIsm(ctx, tokenId, nil, msg.IsmId)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetIsmResponse{}, nil
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}
