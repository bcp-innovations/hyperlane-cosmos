package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

func (qs queryServer) Delivered(ctx context.Context, req *types.QueryDeliveredRequest) (*types.QueryDeliveredResponse, error) {
	messageId, err := util.DecodeEthHex(req.MessageId)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeEthHex(req.Id)
	if err != nil {
		return nil, err
	}

	delivered, err := qs.k.Messages.Has(ctx, collections.Join(mailboxId, messageId))
	if err != nil {
		return nil, err
	}

	return &types.QueryDeliveredResponse{Delivered: delivered}, nil
}

func (qs queryServer) RecipientIsm(ctx context.Context, req *types.RecipientIsmRequest) (*types.RecipientIsmResponse, error) {
	address, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, err
	}

	get, err := qs.k.Hooks().ReceiverIsmId(ctx, address)
	if err != nil {
		return nil, err
	}

	return &types.RecipientIsmResponse{IsmId: get.String()}, nil
}

func (qs queryServer) Mailboxes(ctx context.Context, req *types.QueryMailboxesRequest) (*types.QueryMailboxesResponse, error) {
	values, pagination, err := GetPaginatedFromMap(ctx, qs.k.Mailboxes, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryMailboxesResponse{
		Mailboxes:  values,
		Pagination: pagination,
	}, nil
}

func (qs queryServer) Mailbox(ctx context.Context, req *types.QueryMailboxRequest) (*types.QueryMailboxResponse, error) {
	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

	return &types.QueryMailboxResponse{
		Mailbox: mailbox,
	}, nil
}

// TODO: Remove
func (qs queryServer) VerifyDryRun(ctx context.Context, req *types.QueryVerifyDryRunRequest) (*types.QueryVerifyDryRunResponse, error) {
	panic("Not Implemented")
}

// Params defines the handler for the Query/Params RPC method.
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
