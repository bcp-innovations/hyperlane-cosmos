package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

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

	delivered, err := qs.k.Messages.Has(ctx, messageId)
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
	// values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.Mailboxes, req.Pagination, qs.k.MailboxesSequence)
	// if err != nil {
	// 	return nil, err
	// }

	// return &types.QueryMailboxesResponse{
	// 	Mailboxes:  values,
	// 	Pagination: pagination,
	// }, nil
	panic("TODO: refactor")
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

func (qs queryServer) Count(ctx context.Context, req *types.QueryCountRequest) (*types.QueryCountResponse, error) {
	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return nil, err
	}

	return &types.QueryCountResponse{
		Count: tree.GetCount(),
	}, nil
}

func (qs queryServer) Root(ctx context.Context, req *types.QueryRootRequest) (*types.QueryRootResponse, error) {
	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

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
	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return nil, err
	}

	root, count, err := tree.GetLatestCheckpoint()
	if err != nil {
		return nil, err
	}

	return &types.QueryLatestCheckpointResponse{
		Root:  root[:],
		Count: count,
	}, nil
}

// IGP
func (qs queryServer) QuoteGasPayment(ctx context.Context, req *types.QueryQuoteGasPaymentRequest) (*types.QueryQuoteGasPaymentResponse, error) {
	if len(req.IgpId) == 0 {
		return nil, errors.New("parameter 'igp_id' is required")
	}

	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	if len(req.DestinationDomain) == 0 {
		return nil, errors.New("parameter 'destination_domain' is required")
	}

	destinationDomain, err := strconv.ParseUint(req.DestinationDomain, 10, 32)
	if err != nil {
		return nil, err
	}

	if len(req.GasLimit) == 0 {
		return nil, errors.New("parameter 'gas_limit' is required")
	}

	gasLimit, ok := math.NewIntFromString(req.GasLimit)
	if !ok {
		return nil, fmt.Errorf("failed to convert gasLimit to math.Int")
	}

	payment, err := qs.k.QuoteGasPayment(ctx, igpId, uint32(destinationDomain), gasLimit)
	if err != nil {
		return nil, err
	}

	return &types.QueryQuoteGasPaymentResponse{GasPayment: payment.String()}, nil
}

func (qs queryServer) Igps(ctx context.Context, req *types.QueryIgpsRequest) (*types.QueryIgpsResponse, error) {
	// values, pagination, err := GetPaginatedFromMap(ctx, qs.k.Igp, req.Pagination)
	// if err != nil {
	// 	return nil, err
	// }

	// return &types.QueryIgpsResponse{
	// 	Igps:       values,
	// 	Pagination: pagination,
	// }, nil
	panic("TODO: refactor")
}

func (qs queryServer) Igp(ctx context.Context, req *types.QueryIgpRequest) (*types.QueryIgpResponse, error) {
	igpId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	igp, err := qs.k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find igp with id: %v", igpId.String())
	}

	return &types.QueryIgpResponse{
		Igp: igp,
	}, nil
}

func (qs queryServer) DestinationGasConfigs(ctx context.Context, req *types.QueryDestinationGasConfigsRequest) (*types.QueryDestinationGasConfigsResponse, error) {
	igpId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	rng := collections.NewPrefixedPairRange[[]byte, uint32](igpId.Bytes())

	iter, err := qs.k.IgpDestinationGasConfigMap.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}

	destinationGasConfigs, err := iter.Values()
	if err != nil {
		return nil, err
	}

	configs := make([]*types.DestinationGasConfig, len(destinationGasConfigs))
	for i := range destinationGasConfigs {
		configs[i] = &destinationGasConfigs[i]
	}

	return &types.QueryDestinationGasConfigsResponse{
		DestinationGasConfigs: configs,
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
