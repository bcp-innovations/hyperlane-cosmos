package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

func (qs queryServer) RecipientIsm(ctx context.Context, request *types.RecipientIsmRequest) (*types.RecipientIsmResponse, error) {

	address, err := util.DecodeHexAddress(request.Recipient)
	if err != nil {
		return nil, err
	}

	get, err := qs.k.ReceiverIsmMapping.Get(ctx, address.Bytes())
	if err != nil {
		return nil, err
	}

	return &types.RecipientIsmResponse{IsmId: util.HexAddress(get).String()}, nil
}

func (qs queryServer) Mailboxes(ctx context.Context, _ *types.QueryMailboxesRequest) (*types.QueryMailboxesResponse, error) {
	it, err := qs.k.Mailboxes.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	mailboxes, err := it.Values()
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxesResponse{
		Mailbox: mailboxes,
	}, nil
}

func (qs queryServer) Mailbox(ctx context.Context, req *types.QueryMailboxRequest) (*types.QueryMailboxResponse, error) {
	id, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}
	mailbox, err := qs.k.Mailboxes.Get(ctx, id.Bytes())

	return &types.QueryMailboxResponse{
		Mailbox: mailbox,
	}, nil
}

func (qs queryServer) Count(ctx context.Context, req *types.QueryCountRequest) (*types.QueryCountResponse, error) {
	id, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}
	mailbox, err := qs.k.Mailboxes.Get(ctx, id.Bytes())

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return nil, err
	}

	return &types.QueryCountResponse{
		Count: tree.GetCount(),
	}, nil
}

func (qs queryServer) Root(ctx context.Context, req *types.QueryRootRequest) (*types.QueryRootResponse, error) {
	id, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}
	mailbox, err := qs.k.Mailboxes.Get(ctx, id.Bytes())

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return nil, err
	}

	root := tree.GetRoot()

	return &types.QueryRootResponse{
		Root: root[:],
	}, nil
}

func (qs queryServer) LatestCheckpoint(ctx context.Context, req *types.QueryLatestCheckpointRequest) (*types.QueryLatestCheckpointResponse, error) {
	id, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}
	mailbox, err := qs.k.Mailboxes.Get(ctx, id.Bytes())

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return nil, err
	}

	root, count := tree.GetLatestCheckpoint()

	return &types.QueryLatestCheckpointResponse{
		Root:  root[:],
		Count: count,
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
