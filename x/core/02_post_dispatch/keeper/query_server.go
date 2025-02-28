package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

var _ types.QueryServer = queryServer{}

// NewQueryServerImpl returns an implementation of the module QueryServer.
func NewQueryServerImpl(k *Keeper) types.QueryServer {
	return queryServer{k}
}

type queryServer struct {
	k *Keeper
}

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

	igpHandler := InterchainGasPaymasterHookHandler{*qs.k}

	payment, err := igpHandler.QuoteGasPayment(ctx, util.NewZeroAddress(), igpId, uint32(destinationDomain), gasLimit)
	if err != nil {
		return nil, err
	}

	return &types.QueryQuoteGasPaymentResponse{GasPayment: payment.String()}, nil
}

//
// Interchain Gas Paymaster

func (qs queryServer) Igps(ctx context.Context, req *types.QueryIgpsRequest) (*types.QueryIgpsResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.Igps, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryIgpsResponse{
		Igps:       values,
		Pagination: pagination,
	}, nil
}

func (qs queryServer) Igp(ctx context.Context, req *types.QueryIgpRequest) (*types.QueryIgpResponse, error) {
	igpId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	igp, err := qs.k.Igps.Get(ctx, igpId.GetInternalId())
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

	rng := collections.NewPrefixedPairRange[uint64, uint32](igpId.GetInternalId())

	iter, err := qs.k.IgpDestinationGasConfigs.Iterate(ctx, rng)
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

// Merkle Tree Hook
func (qs queryServer) MerkleTreeHooks(ctx context.Context, req *types.QueryMerkleTreeHooks) (*types.QueryMerkleTreeHooksResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.merkleTreeHooks, req.Pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]*types.QueryMerkleTreeHookResponse, len(values))
	for i := 0; i < len(values); i++ {
		merkleTreeHook := values[i]
		tree, err := types.TreeFromProto(merkleTreeHook.Tree)
		if err != nil {
			return nil, err
		}

		root := tree.GetRoot()
		responses[i] = &types.QueryMerkleTreeHookResponse{
			Root:           root[:],
			Count:          tree.Count,
			MerkleTreeHook: &merkleTreeHook,
		}
	}

	return &types.QueryMerkleTreeHooksResponse{
		MerkleTreeHooks: responses,
		Pagination:      pagination,
	}, nil
}

func (qs queryServer) MerkleTreeHook(ctx context.Context, req *types.QueryMerkleTreeHook) (*types.QueryMerkleTreeHookResponse, error) {
	merkleTreeHooksId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	merkleTreeHook, err := qs.k.merkleTreeHooks.Get(ctx, merkleTreeHooksId.GetInternalId())
	if err != nil {
		return nil, err
	}

	tree, err := types.TreeFromProto(merkleTreeHook.Tree)
	if err != nil {
		return nil, err
	}

	root := tree.GetRoot()

	return &types.QueryMerkleTreeHookResponse{
		Root:           root[:],
		Count:          tree.Count,
		MerkleTreeHook: &merkleTreeHook,
	}, nil
}

//
// Noop Hook

func (qs queryServer) NoopHook(ctx context.Context, req *types.QueryNoopHook) (*types.QueryNoopHookResponse, error) {
	hookId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	noopHook, err := qs.k.noopHooks.Get(ctx, hookId.Bytes())
	if err != nil {
		return nil, err
	}

	return &types.QueryNoopHookResponse{
		NoopHook: &noopHook,
	}, nil
}
