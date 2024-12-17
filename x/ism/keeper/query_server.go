package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/ism/types"
	mailboxTypes "github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
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

func (qs queryServer) Isms(ctx context.Context, _ *types.QueryIsmsRequest) (*types.QueryIsmsResponse, error) {
	it, err := qs.k.Isms.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	isms, err := it.Values()
	if err != nil {
		return nil, err
	}

	return &types.QueryIsmsResponse{
		Isms: isms,
	}, nil
}

func (qs queryServer) VerifyDryRun(ctx context.Context, req *types.QueryVerifyDryRunRequest) (*types.QueryVerifyDryRunResponse, error) {
	rawMessage, err := util.DecodeEthHex(req.Message)
	if err != nil {
		return nil, err
	}

	message, err := mailboxTypes.ParseHyperlaneMessage(rawMessage)
	if err != nil {
		return nil, err
	}

	metadata, err := util.DecodeEthHex(req.Metadata)
	if err != nil {
		return nil, err
	}

	ismId, err := util.DecodeHexAddress(req.IsmId)
	if err != nil {
		return nil, err
	}

	verified, err := qs.k.Verify(ctx, ismId, metadata, message)
	if err != nil {
		return nil, err
	}

	return &types.QueryVerifyDryRunResponse{
		Verified: verified,
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
