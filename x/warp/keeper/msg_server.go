package keeper

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"cosmossdk.io/collections"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	k Keeper
}

func (ms msgServer) CreateSyntheticToken(ctx context.Context, msg *types.MsgCreateSyntheticToken) (*types.MsgCreateSyntheticTokenResponse, error) {
	if !slices.Contains(ms.k.enabledTokens, int32(types.HYP_TOKEN_TYPE_SYNTHETIC)) {
		return nil, fmt.Errorf("module disabled synthetic tokens")
	}

	next, err := ms.k.HypTokensCount.Next(ctx)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeHexAddress(msg.OriginMailbox)
	if err != nil {
		return nil, fmt.Errorf("invalid mailbox id: %s", err)
	}

	has, err := ms.k.mailboxKeeper.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", mailboxId.String())
	}

	if err = ValidateTokenMetadata(msg.Metadata); err != nil {
		return nil, err
	}

	tokenId := ms.k.hexAddressFactory.GenerateId(uint32(types.HYP_TOKEN_TYPE_SYNTHETIC), next)

	newToken := types.HypToken{
		Id:            next,
		Owner:         msg.Owner,
		TokenType:     types.HYP_TOKEN_TYPE_SYNTHETIC,
		OriginMailbox: mailboxId.Bytes(),
		OriginDenom:   fmt.Sprintf("hyperlane/%s", tokenId.String()),
		Metadata:      msg.Metadata,
	}

	if err = ms.k.HypTokens.Set(ctx, newToken.Id, newToken); err != nil {
		return nil, err
	}

	return &types.MsgCreateSyntheticTokenResponse{Id: tokenId.String()}, nil
}

// CreateCollateralToken ...
func (ms msgServer) CreateCollateralToken(ctx context.Context, msg *types.MsgCreateCollateralToken) (*types.MsgCreateCollateralTokenResponse, error) {
	if !slices.Contains(ms.k.enabledTokens, int32(types.HYP_TOKEN_TYPE_COLLATERAL)) {
		return nil, fmt.Errorf("module disabled collateral tokens")
	}

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
		return nil, fmt.Errorf("invalid mailbox id: %s", err)
	}

	has, err := ms.k.mailboxKeeper.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", mailboxId.String())
	}

	newToken := types.HypToken{
		Id:            next,
		Owner:         msg.Owner,
		TokenType:     types.HYP_TOKEN_TYPE_COLLATERAL,
		OriginMailbox: mailboxId.Bytes(),
		OriginDenom:   msg.OriginDenom,
		Metadata:      nil,
	}

	if err = ms.k.HypTokens.Set(ctx, newToken.Id, newToken); err != nil {
		return nil, err
	}
	return &types.MsgCreateCollateralTokenResponse{Id: ms.k.GetAddressFromToken(newToken).String()}, nil
}

func (ms msgServer) EnrollRemoteRouter(ctx context.Context, msg *types.MsgEnrollRemoteRouter) (*types.MsgEnrollRemoteRouterResponse, error) {
	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("token with id %s not found", tokenId.String())
	}

	if token.Owner != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	if msg.RemoteRouter == nil {
		return nil, fmt.Errorf("invalid remote router")
	}

	exists, err := ms.k.EnrolledRouters.Has(ctx, collections.Join(tokenId.GetInternalId(), msg.RemoteRouter.ReceiverDomain))
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("remote router for domain %v is already enrolled", msg.RemoteRouter.ReceiverDomain)
	}

	if msg.RemoteRouter.ReceiverContract == "" {
		return nil, fmt.Errorf("invalid receiver contract")
	}

	if err = ms.k.EnrolledRouters.Set(ctx, collections.Join(tokenId.GetInternalId(), msg.RemoteRouter.ReceiverDomain), *msg.RemoteRouter); err != nil {
		return nil, err
	}

	return &types.MsgEnrollRemoteRouterResponse{}, nil
}

func (ms msgServer) SetRemoteRouter(ctx context.Context, msg *types.MsgSetRemoteRouter) (*types.MsgSetRemoteRouterResponse, error) {
	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("token with id %s not found", tokenId.String())
	}

	if token.Owner != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	if msg.RemoteRouter == nil {
		return nil, fmt.Errorf("invalid remote router")
	}

	exists, err := ms.k.EnrolledRouters.Has(ctx, collections.Join(tokenId.GetInternalId(), msg.RemoteRouter.ReceiverDomain))
	if err != nil || !exists {
		return nil, fmt.Errorf("failed to find remote router for domain %v", msg.RemoteRouter.ReceiverDomain)
	}

	if msg.RemoteRouter.ReceiverContract == "" {
		return nil, fmt.Errorf("invalid receiver contract")
	}

	if err = ms.k.EnrolledRouters.Set(ctx, collections.Join(tokenId.GetInternalId(), msg.RemoteRouter.ReceiverDomain), *msg.RemoteRouter); err != nil {
		return nil, err
	}

	return &types.MsgSetRemoteRouterResponse{}, nil
}

func (ms msgServer) UnrollRemoteRouter(ctx context.Context, msg *types.MsgUnrollRemoteRouter) (*types.MsgUnrollRemoteRouterResponse, error) {
	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("token with id %s not found", tokenId.String())
	}

	if token.Owner != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	exists, err := ms.k.EnrolledRouters.Has(ctx, collections.Join(tokenId.GetInternalId(), msg.ReceiverDomain))
	if err != nil || !exists {
		return nil, fmt.Errorf("failed to find remote router for domain %v", msg.ReceiverDomain)
	}

	if err = ms.k.EnrolledRouters.Remove(ctx, collections.Join(tokenId.GetInternalId(), msg.ReceiverDomain)); err != nil {
		return nil, err
	}

	return &types.MsgUnrollRemoteRouterResponse{}, nil
}

func (ms msgServer) RemoteTransfer(ctx context.Context, msg *types.MsgRemoteTransfer) (*types.MsgRemoteTransferResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find token with id: %s", tokenId.String())
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

func (ms msgServer) SetInterchainSecurityModule(ctx context.Context, msg *types.MsgSetInterchainSecurityModule) (*types.MsgSetInterchainSecurityModuleResponse, error) {
	if msg.IsmId == "" {
		return nil, fmt.Errorf("ism id cannot be empty")
	}

	tokenId, err := util.DecodeHexAddress(msg.TokenId)
	if err != nil {
		return nil, fmt.Errorf("invalid token id %s", msg.TokenId)
	}

	token, err := ms.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, err
	}

	if token.Owner != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	ismAddress, err := util.DecodeHexAddress(msg.IsmId)
	if err != nil {
		return nil, fmt.Errorf("invalid ism id: %s", err)
	}

	token.IsmId = ismAddress.String()

	err = ms.k.HypTokens.Set(ctx, tokenId.GetInternalId(), token)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetInterchainSecurityModuleResponse{}, nil
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}
