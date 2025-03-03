package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

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

	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	delivered, err := qs.k.Messages.Has(ctx, collections.Join(mailboxId.GetInternalId(), messageId))
	if err != nil {
		return nil, err
	}

	return &types.QueryDeliveredResponse{Delivered: delivered}, nil
}

func (qs queryServer) RecipientIsm(ctx context.Context, req *types.RecipientIsmRequest) (*types.RecipientIsmResponse, error) {
	recipient, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, err
	}

	get, err := qs.k.ReceiverIsmId(ctx, recipient)
	if err != nil {
		return nil, err
	}

	return &types.RecipientIsmResponse{IsmId: get.String()}, nil
}

func (qs queryServer) Mailboxes(ctx context.Context, req *types.QueryMailboxesRequest) (*types.QueryMailboxesResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.Mailboxes, req.Pagination)
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

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

	return &types.QueryMailboxResponse{
		Mailbox: mailbox,
	}, nil
}

func (qs queryServer) VerifyDryRun(ctx context.Context, req *types.QueryVerifyDryRunRequest) (*types.QueryVerifyDryRunResponse, error) {
	ismId, err := util.DecodeHexAddress(req.IsmId)
	if err != nil {
		return nil, err
	}

	metadata := []byte(req.Metadata)

	msg, err := util.ParseHyperlaneMessage([]byte(req.Message))
	if err != nil {
		return nil, err
	}

	verified, err := qs.k.Verify(ctx, ismId, metadata, msg)
	return &types.QueryVerifyDryRunResponse{
		Verified: verified,
	}, err
}
