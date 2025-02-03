package keeper

import (
	"context"
	"errors"

	"cosmossdk.io/collections"

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

func (qs queryServer) Params(ctx context.Context, request *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryParamsResponse{Params: types.Params{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}

func (qs queryServer) Tokens(ctx context.Context, request *types.QueryTokensRequest) (*types.QueryTokensResponse, error) {
	it, err := qs.k.HypTokens.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	tokens, err := it.Values()
	if err != nil {
		return nil, err
	}

	responseTokens := make([]types.QueryTokenResponse, len(tokens))
	for i, t := range tokens {
		responseTokens[i] = parseTokenResponse(t)
	}

	return &types.QueryTokensResponse{
		Tokens: responseTokens,
	}, nil
}

func (qs queryServer) Token(ctx context.Context, request *types.QueryMailboxRequest) (*types.QueryTokenResponse, error) {
	tokenId, err := util.DecodeHexAddress(request.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	get, err := qs.k.HypTokens.Get(ctx, tokenId.Bytes())
	if err != nil {
		return nil, err
	}

	response := parseTokenResponse(get)
	return &response, nil
}

func parseTokenResponse(get types.HypToken) types.QueryTokenResponse {
	return types.QueryTokenResponse{
		Id:        util.HexAddress(get.Id).String(),
		Creator:   get.Creator,
		TokenType: get.TokenType,

		OriginMailbox:    util.HexAddress(get.OriginMailbox).String(),
		OriginDenom:      get.OriginDenom,
		ReceiverDomain:   get.ReceiverDomain,
		ReceiverContract: util.HexAddress(get.ReceiverContract).String(),
	}
}
