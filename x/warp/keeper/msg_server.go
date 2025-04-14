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

// CreateSyntheticToken ...
func (ms msgServer) CreateSyntheticToken(ctx context.Context, msg *types.MsgCreateSyntheticToken) (*types.MsgCreateSyntheticTokenResponse, error) {
	if !slices.Contains(ms.k.enabledTokens, int32(types.HYP_TOKEN_TYPE_SYNTHETIC)) {
		return nil, fmt.Errorf("module disabled synthetic tokens")
	}

	has, err := ms.k.coreKeeper.MailboxIdExists(ctx, msg.OriginMailbox)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", msg.OriginMailbox.String())
	}

	tokenId, err := ms.k.coreKeeper.AppRouter().GetNextSequence(ctx, uint8(types.HYP_TOKEN_TYPE_SYNTHETIC))
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:            tokenId,
		Owner:         msg.Owner,
		TokenType:     types.HYP_TOKEN_TYPE_SYNTHETIC,
		OriginMailbox: msg.OriginMailbox,
		OriginDenom:   fmt.Sprintf("hyperlane/%s", tokenId.String()),
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.GetInternalId(), newToken); err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateSyntheticToken{
		TokenId:       newToken.Id,
		Owner:         newToken.Owner,
		OriginMailbox: newToken.OriginMailbox,
		OriginDenom:   newToken.OriginDenom,
	})

	return &types.MsgCreateSyntheticTokenResponse{Id: tokenId}, nil
}

// CreateCollateralToken ...
func (ms msgServer) CreateCollateralToken(ctx context.Context, msg *types.MsgCreateCollateralToken) (*types.MsgCreateCollateralTokenResponse, error) {
	if !slices.Contains(ms.k.enabledTokens, int32(types.HYP_TOKEN_TYPE_COLLATERAL)) {
		return nil, fmt.Errorf("module disabled collateral tokens")
	}

	err := sdk.ValidateDenom(msg.OriginDenom)
	if err != nil {
		return nil, fmt.Errorf("origin denom %s is invalid", msg.OriginDenom)
	}

	has, err := ms.k.coreKeeper.MailboxIdExists(ctx, msg.OriginMailbox)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", msg.OriginMailbox.String())
	}

	tokenId, err := ms.k.coreKeeper.AppRouter().GetNextSequence(ctx, uint8(types.HYP_TOKEN_TYPE_COLLATERAL))
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:            tokenId,
		Owner:         msg.Owner,
		TokenType:     types.HYP_TOKEN_TYPE_COLLATERAL,
		OriginMailbox: msg.OriginMailbox,
		OriginDenom:   msg.OriginDenom,
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.GetInternalId(), newToken); err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateCollateralToken{
		TokenId:       newToken.Id,
		Owner:         newToken.Owner,
		OriginMailbox: newToken.OriginMailbox,
		OriginDenom:   newToken.OriginDenom,
	})

	return &types.MsgCreateCollateralTokenResponse{Id: tokenId}, nil
}

// SetToken allows the owner of a token to change its ownership or update its ISM ID.
func (ms msgServer) SetToken(ctx context.Context, msg *types.MsgSetToken) (*types.MsgSetTokenResponse, error) {
	if msg.NewOwner == "" && msg.IsmId == nil && !msg.RenounceOwnership {
		return nil, fmt.Errorf("new owner, renounce ownership or ism id required")
	}

	tokenId := msg.TokenId
	token, err := ms.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find token with id: %s", tokenId.String())
	}

	if token.Owner != msg.Owner {
		return nil, fmt.Errorf("%s does not own token with id %s", msg.Owner, tokenId.String())
	}

	// Only renounce if new owner is empty
	if msg.RenounceOwnership && msg.NewOwner != "" {
		return nil, fmt.Errorf("cannot set new owner and renounce ownership at the same time")
	}

	if msg.NewOwner != "" {
		if _, err := sdk.AccAddressFromBech32(msg.NewOwner); err != nil {
			return nil, fmt.Errorf("invalid new owner")
		}
		token.Owner = msg.NewOwner
	}

	if msg.RenounceOwnership {
		token.Owner = ""
	}

	if msg.IsmId != nil {
		if err := ms.k.coreKeeper.AssertIsmExists(ctx, *msg.IsmId); err != nil {
			return nil, err
		}
		token.IsmId = msg.IsmId
	}

	err = ms.k.HypTokens.Set(ctx, tokenId.GetInternalId(), token)
	if err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetToken{
		TokenId:           tokenId.String(),
		Owner:             msg.Owner,
		IsmId:             msg.IsmId,
		NewOwner:          msg.NewOwner,
		RenounceOwnership: msg.RenounceOwnership,
	})

	return &types.MsgSetTokenResponse{}, nil
}

// EnrollRemoteRouter enrolls a new remote router for a specific token.
func (ms msgServer) EnrollRemoteRouter(ctx context.Context, msg *types.MsgEnrollRemoteRouter) (*types.MsgEnrollRemoteRouterResponse, error) {
	tokenId := msg.TokenId
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

	if msg.RemoteRouter.ReceiverContract == "" {
		return nil, fmt.Errorf("invalid receiver contract")
	}

	if err = ms.k.EnrolledRouters.Set(ctx, collections.Join(tokenId.GetInternalId(), msg.RemoteRouter.ReceiverDomain), *msg.RemoteRouter); err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventEnrollRemoteRouter{
		TokenId:          tokenId.String(),
		Owner:            msg.Owner,
		ReceiverDomain:   msg.RemoteRouter.ReceiverDomain,
		ReceiverContract: msg.RemoteRouter.ReceiverContract,
		Gas:              msg.RemoteRouter.Gas,
	})

	return &types.MsgEnrollRemoteRouterResponse{}, nil
}

// UnrollRemoteRouter removes an existing remote router from a token.
func (ms msgServer) UnrollRemoteRouter(ctx context.Context, msg *types.MsgUnrollRemoteRouter) (*types.MsgUnrollRemoteRouterResponse, error) {
	tokenId := msg.TokenId
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

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventUnrollRemoteRouter{
		TokenId:        tokenId.String(),
		Owner:          msg.Owner,
		ReceiverDomain: msg.ReceiverDomain,
	})

	return &types.MsgUnrollRemoteRouterResponse{}, nil
}

// RemoteTransfer handles the transfer of tokens (collateral or synthetic) to a remote chain.
func (ms msgServer) RemoteTransfer(ctx context.Context, msg *types.MsgRemoteTransfer) (*types.MsgRemoteTransferResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	token, err := ms.k.HypTokens.Get(ctx, msg.TokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find token with id: %s", msg.TokenId.String())
	}

	customHookMetadata, err := util.DecodeEthHex(msg.CustomHookMetadata)
	if err != nil {
		return nil, fmt.Errorf("invalid custom hook metadata")
	}

	var messageResultId util.HexAddress
	if token.TokenType == types.HYP_TOKEN_TYPE_COLLATERAL {
		result, err := ms.k.RemoteTransferCollateral(goCtx, token, msg.Sender, msg.DestinationDomain, msg.Recipient, msg.Amount, msg.CustomHookId, msg.GasLimit, msg.MaxFee, customHookMetadata)
		if err != nil {
			return nil, err
		}
		messageResultId = result
	} else if token.TokenType == types.HYP_TOKEN_TYPE_SYNTHETIC {
		result, err := ms.k.RemoteTransferSynthetic(goCtx, token, msg.Sender, msg.DestinationDomain, msg.Recipient, msg.Amount, msg.CustomHookId, msg.GasLimit, msg.MaxFee, customHookMetadata)
		if err != nil {
			return nil, err
		}
		messageResultId = result
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
