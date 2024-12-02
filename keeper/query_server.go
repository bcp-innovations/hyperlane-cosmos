package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KYVENetwork/mailbox"
	"github.com/cosmos/cosmos-sdk/types/query"
)

var _ mailbox.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k Keeper) mailbox.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k Keeper
}

// Counter defines the handler for the Query/Counter RPC method.
func (qs queryServer) Counter(ctx context.Context, req *mailbox.QueryCounterRequest) (*mailbox.QueryCounterResponse, error) {
	if _, err := qs.k.addressCodec.StringToBytes(req.Address); err != nil {
		return nil, fmt.Errorf("invalid sender address: %w", err)
	}

	counter, err := qs.k.Counter.Get(ctx, req.Address)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &mailbox.QueryCounterResponse{Counter: 0}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &mailbox.QueryCounterResponse{Counter: counter}, nil
}

// Counters defines the handler for the Query/Counters RPC method.
func (qs queryServer) Counters(ctx context.Context, req *mailbox.QueryCountersRequest) (*mailbox.QueryCountersResponse, error) {
	counters, pageRes, err := query.CollectionPaginate(
		ctx,
		qs.k.Counter,
		req.Pagination,
		func(key string, value uint64) (*mailbox.Counter, error) {
			return &mailbox.Counter{
				Address: key,
				Count:   value,
			}, nil
		})
	if err != nil {
		return nil, err
	}

	return &mailbox.QueryCountersResponse{Counters: counters, Pagination: pageRes}, nil
}

// Params defines the handler for the Query/Params RPC method.
func (qs queryServer) Params(ctx context.Context, req *mailbox.QueryParamsRequest) (*mailbox.QueryParamsResponse, error) {
	params, err := qs.k.Params.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return &mailbox.QueryParamsResponse{Params: mailbox.Params{}}, nil
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &mailbox.QueryParamsResponse{Params: params}, nil
}
