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

func (qs queryServer) AnnouncedStorageLocations(ctx context.Context, req *types.QueryAnnouncedStorageLocationsRequest) (*types.QueryAnnouncedStorageLocationsResponse, error) {
	validatorAddress, err := util.DecodeEthHex(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateValidatorStorageKey(validatorAddress)

	validator, err := qs.k.Validators.Get(ctx, prefixedId.Bytes())
	if err != nil {
		return nil, err
	}

	return &types.QueryAnnouncedStorageLocationsResponse{
		StorageLocations: validator.StorageLocations,
	}, nil
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
		Mailboxes: mailboxes,
	}, nil
}

func (qs queryServer) Mailbox(ctx context.Context, req *types.QueryMailboxRequest) (*types.QueryMailboxResponse, error) {
	mailboxId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	mailbox, err := qs.k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find Mailbox with Id: %v", mailboxId.String())
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
		return nil, fmt.Errorf("failed to find Mailbox with Id: %v", mailboxId.String())
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
		return nil, fmt.Errorf("failed to find Mailbox with Id: %v", mailboxId.String())
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
		return nil, fmt.Errorf("failed to find Mailbox with Id: %v", mailboxId.String())
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

func (qs queryServer) Validators(ctx context.Context, _ *types.QueryValidatorsRequest) (*types.QueryValidatorsResponse, error) {
	it, err := qs.k.Validators.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	validators, err := it.Values()
	if err != nil {
		return nil, err
	}

	return &types.QueryValidatorsResponse{
		Validators: validators,
	}, nil
}

// IGP

func (qs queryServer) QuoteGasPayment(ctx context.Context, req *types.QueryQuoteGasPaymentRequest) (*types.QueryQuoteGasPaymentResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	destinationDomain, err := strconv.ParseUint(req.DestinationDomain, 10, 32)
	if err != nil {
		return nil, err
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

func (qs queryServer) Igps(ctx context.Context, _ *types.QueryIgpsRequest) (*types.QueryIgpsResponse, error) {
	it, err := qs.k.Igp.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	igps, err := it.Values()
	if err != nil {
		return nil, err
	}

	return &types.QueryIgpsResponse{
		Igps: igps,
	}, nil
}

func (qs queryServer) Igp(ctx context.Context, req *types.QueryIgpRequest) (*types.QueryIgpResponse, error) {
	igpId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	igp, err := qs.k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find IGP with Id: %v", igpId.String())
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

// ISM

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

func (qs queryServer) Ism(ctx context.Context, req *types.QueryIsmRequest) (*types.QueryIsmResponse, error) {
	ismId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	ism, err := qs.k.Isms.Get(ctx, ismId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to find ISM with Id: %v", ismId.String())
	}

	return &types.QueryIsmResponse{
		Ism: ism,
	}, nil
}

func (qs queryServer) VerifyDryRun(ctx context.Context, req *types.QueryVerifyDryRunRequest) (*types.QueryVerifyDryRunResponse, error) {
	rawMessage, err := util.DecodeEthHex(req.Message)
	if err != nil {
		return nil, err
	}

	message, err := types.ParseHyperlaneMessage(rawMessage)
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
