package keeper

import (
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

// ISM
func (qs queryServer) Isms(ctx context.Context, req *types.QueryIsmsRequest) (*types.QueryIsmsResponse, error) {
	values, pagination, err := GetPaginatedFromMap(ctx, qs.k.isms, req.Pagination)
	if err != nil {
		return nil, err
	}

	isms, _ := packAccounts(values)

	return &types.QueryIsmsResponse{
		Isms:       isms,
		Pagination: pagination,
	}, nil
}

func (qs queryServer) Ism(ctx context.Context, req *types.QueryIsmRequest) (*types.QueryIsmResponse, error) {
	ismId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	ism, err := qs.k.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find ism with id: %v", ismId.String())
	}

	toAny, err := fromProtoToAny(ism)
	if err != nil {
		return nil, err
	}

	return &types.QueryIsmResponse{
		Ism: *toAny,
	}, nil
}
