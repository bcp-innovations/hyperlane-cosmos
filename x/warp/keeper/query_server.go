package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (qs queryServer) RemoteRouters(ctx context.Context, request *types.QueryRemoteRoutersRequest) (*types.QueryRemoteRoutersResponse, error) {
	tokenId, err := util.DecodeHexAddress(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	routers, page, err := util.GetPaginatedPrefixFromMap(ctx, qs.k.EnrolledRouters, request.Pagination, tokenId.GetInternalId())
	if err != nil {
		return &types.QueryRemoteRoutersResponse{}, err
	}

	remoteRouters := make([]*types.RemoteRouter, len(routers))
	for i := range routers {
		remoteRouters[i] = &routers[i]
	}

	return &types.QueryRemoteRoutersResponse{
		RemoteRouters: remoteRouters,
		Pagination:    page,
	}, nil
}

func (qs queryServer) Params(ctx context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryParamsResponse{Params: types.Params{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

func (qs queryServer) BridgedSupply(ctx context.Context, request *types.QueryBridgedSupplyRequest) (*types.QueryBridgedSupplyResponse, error) {
	tokenId, err := util.DecodeHexAddress(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := qs.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, err
	}

	var amount math.Int
	switch token.TokenType {
	case types.HYP_TOKEN_TYPE_COLLATERAL:
		amount = token.CollateralBalance
	case types.HYP_TOKEN_TYPE_SYNTHETIC:
		amount = qs.k.bankKeeper.GetSupply(ctx, token.OriginDenom).Amount
	default:
		return nil, fmt.Errorf("query doesn't support token type: %s", token.TokenType)
	}

	bridgedSupply := sdk.Coin{
		Amount: amount,
		Denom:  token.OriginDenom,
	}

	err = bridgedSupply.Validate()
	if err != nil {
		return nil, err
	}

	return &types.QueryBridgedSupplyResponse{BridgedSupply: bridgedSupply}, nil
}

func (qs queryServer) QuoteRemoteTransfer(ctx context.Context, request *types.QueryQuoteRemoteTransferRequest) (*types.QueryQuoteRemoteTransferResponse, error) {
	tokenId, err := util.DecodeHexAddress(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := qs.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, err
	}

	destinationDomain, err := strconv.ParseUint(request.DestinationDomain, 10, 32)
	if err != nil {
		return nil, err
	}

	remoteRouter, err := qs.k.EnrolledRouters.Get(ctx, collections.Join(tokenId.GetInternalId(), uint32(destinationDomain)))
	if err != nil {
		return nil, fmt.Errorf("failed to get remote router for destination domain %v", request.DestinationDomain)
	}

	metadata := util.StandardHookMetadata{
		GasLimit:           remoteRouter.Gas,
		Address:            sdk.AccAddress{},
		CustomHookMetadata: []byte{},
	}

	requiredPayment, err := qs.k.coreKeeper.QuoteDispatch(ctx, util.HexAddress(token.OriginMailbox), util.NewZeroAddress(), metadata, util.HyperlaneMessage{Destination: uint32(destinationDomain)})
	if err != nil {
		return nil, err
	}

	return &types.QueryQuoteRemoteTransferResponse{GasPayment: requiredPayment}, nil
}

func (qs queryServer) Tokens(ctx context.Context, req *types.QueryTokensRequest) (*types.QueryTokensResponse, error) {
	tokens, page, err := util.GetPaginatedFromMap(ctx, qs.k.HypTokens, req.Pagination)
	if err != nil {
		return nil, err
	}

	response := make([]types.WrappedHypToken, 0, len(tokens))
	for _, t := range tokens {
		response = append(response, *parseTokenResponse(t))
	}

	return &types.QueryTokensResponse{
		Tokens:     response,
		Pagination: page,
	}, nil
}

func (qs queryServer) Token(ctx context.Context, request *types.QueryTokenRequest) (*types.QueryTokenResponse, error) {
	tokenId, err := util.DecodeHexAddress(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	get, err := qs.k.HypTokens.Get(ctx, tokenId.GetInternalId())
	if err != nil {
		return nil, err
	}

	return &types.QueryTokenResponse{
		Token: parseTokenResponse(get),
	}, nil
}

func parseTokenResponse(get types.HypToken) *types.WrappedHypToken {
	return &types.WrappedHypToken{
		Id:        get.Id,
		Owner:     get.Owner,
		TokenType: get.TokenType,

		OriginMailbox: util.HexAddress(get.OriginMailbox).String(),
		OriginDenom:   get.OriginDenom,

		IsmId: get.IsmId,
	}
}
