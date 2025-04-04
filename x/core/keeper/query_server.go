package keeper

import (
	"context"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"

	storetypes "cosmossdk.io/store/types"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

func (qs queryServer) RecipientIsm(ctx context.Context, req *types.QueryRecipientIsmRequest) (*types.QueryRecipientIsmResponse, error) {
	recipient, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, err
	}

	get, err := qs.k.ReceiverIsmId(ctx, recipient)
	if err != nil {
		return nil, err
	}

	return &types.QueryRecipientIsmResponse{IsmId: get.String()}, nil
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
	limit := uint64(200000)
	if req.GasLimit != "" {
		parsed, err := strconv.ParseUint(req.GasLimit, 10, 32)
		if err != nil {
			return nil, err
		}

		limit = parsed
	}

	// explicitly set a GasMeter to not run into stack overflows as ISM might call themselves
	sdkCtx := sdk.UnwrapSDKContext(ctx).WithGasMeter(storetypes.NewGasMeter(limit))

	ismId, err := util.DecodeHexAddress(req.IsmId)
	if err != nil {
		return nil, err
	}

	metadata := []byte(req.Metadata)

	msg, err := util.ParseHyperlaneMessage([]byte(req.Message))
	if err != nil {
		return nil, err
	}

	verified, err := qs.k.Verify(sdkCtx, ismId, metadata, msg)
	return &types.QueryVerifyDryRunResponse{
		Verified: verified,
	}, err
}

func (qs queryServer) RegisteredISMs(_ context.Context, _ *types.QueryRegisteredISMs) (*types.QueryRegisteredISMsResponse, error) {
	return &types.QueryRegisteredISMsResponse{
		Ids: qs.k.IsmRouter().GetModuleIds(),
	}, nil
}

func (qs queryServer) RegisteredHooks(_ context.Context, _ *types.QueryRegisteredHooks) (*types.QueryRegisteredHooksResponse, error) {
	return &types.QueryRegisteredHooksResponse{
		Ids: qs.k.PostDispatchRouter().GetModuleIds(),
	}, nil
}

func (qs queryServer) RegisteredApps(_ context.Context, _ *types.QueryRegisteredApps) (*types.QueryRegisteredAppsResponse, error) {
	return &types.QueryRegisteredAppsResponse{
		Ids: qs.k.AppRouter().GetModuleIds(),
	}, nil
}
