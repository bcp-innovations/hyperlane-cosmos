package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

func (qs queryServer) Mailboxes(ctx context.Context, request *types.QueryMailboxesRequest) (*types.QueryMailboxesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (qs queryServer) Mailbox(ctx context.Context, request *types.QueryMailboxRequest) (*types.QueryMailboxResponse, error) {
	//TODO implement me
	panic("implement me")
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
