package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"errors"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (qs queryServer) Igps(ctx context.Context, _ *types.QueryIgpsRequest) (*types.QueryIgpsResponse, error) {
	it, err := qs.k.Igp.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	igps, err := it.Values()
	if err != nil {
		return nil, err
	}

	return &types.QueryIgpsResponse{
		Igps: igps,
	}, nil
}

// Params defines the handler for the Query/Params RPC method.
func (qs queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &types.QueryParamsResponse{Params: types.Params{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryParamsResponse{Params: params}, nil
}
