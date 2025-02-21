package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

func (q queryServer) Igps(ctx context.Context, request *types.QueryIgpsRequest) (*types.QueryIgpsResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (q queryServer) Igp(ctx context.Context, request *types.QueryIgpRequest) (*types.QueryIgpResponse, error) {
	// TODO implement me
	panic("implement me")
}
