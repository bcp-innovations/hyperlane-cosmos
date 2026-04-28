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
		return nil, errors.New("failed to convert gasLimit to math.Int")
	}

	igpHandler := InterchainGasPaymasterHookHandler{*qs.k}

	payment, _, err := igpHandler.QuoteGasPayment(ctx, igpId, uint32(destinationDomain), gasLimit)
	if err != nil {
		return nil, err
	}

	return &types.QueryQuoteGasPaymentResponse{GasPayment: payment}, nil
}

//
// Merkle Tree Hook

func (qs queryServer) MerkleTreeHooks(ctx context.Context, req *types.QueryMerkleTreeHooksRequest) (*types.QueryMerkleTreeHooksResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.merkleTreeHooks, req.Pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]types.WrappedMerkleTreeHookResponse, len(values))
	for i := 0; i < len(values); i++ {
		merkleTreeHook := values[i]
		tree, err := types.TreeFromProto(merkleTreeHook.Tree)
		if err != nil {
			return nil, err
		}

		root := tree.GetRoot()

		responses[i] = types.WrappedMerkleTreeHookResponse{
			Id:        merkleTreeHook.Id.String(),
			Owner:     merkleTreeHook.Owner,
			MailboxId: merkleTreeHook.MailboxId.String(),
			MerkleTree: &types.TreeResponse{
				Count: merkleTreeHook.Tree.Count,
				Root:  root[:],
				Leafs: merkleTreeHook.Tree.Branch,
			},
		}
	}

	return &types.QueryMerkleTreeHooksResponse{
		MerkleTreeHooks: responses,
		Pagination:      pagination,
	}, nil
}

func (qs queryServer) MerkleTreeHook(ctx context.Context, req *types.QueryMerkleTreeHookRequest) (*types.QueryMerkleTreeHookResponse, error) {
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

	return &types.QueryMerkleTreeHookResponse{MerkleTreeHook: types.WrappedMerkleTreeHookResponse{
		Id:        merkleTreeHook.Id.String(),
		Owner:     merkleTreeHook.Owner,
		MailboxId: merkleTreeHook.MailboxId.String(),
		MerkleTree: &types.TreeResponse{
			Count: merkleTreeHook.Tree.Count,
			Root:  root[:],
			Leafs: merkleTreeHook.Tree.Branch,
		},
	}}, nil
}

//
// Noop Hook

func (qs queryServer) NoopHook(ctx context.Context, req *types.QueryNoopHookRequest) (*types.QueryNoopHookResponse, error) {
	hookId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	noopHook, err := qs.k.noopHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	return &types.QueryNoopHookResponse{
		NoopHook: &noopHook,
	}, nil
}

func (qs queryServer) NoopHooks(ctx context.Context, req *types.QueryNoopHooksRequest) (*types.QueryNoopHooksResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.noopHooks, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryNoopHooksResponse{
		NoopHooks:  values,
		Pagination: pagination,
	}, nil
}

//
// Aggregation Hook

func (qs queryServer) AggregationHook(ctx context.Context, req *types.QueryAggregationHookRequest) (*types.QueryAggregationHookResponse, error) {
	hookId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	aggregationHook, err := qs.k.aggregationHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregationHookResponse{
		AggregationHook: &aggregationHook,
	}, nil
}

func (qs queryServer) AggregationHooks(ctx context.Context, req *types.QueryAggregationHooksRequest) (*types.QueryAggregationHooksResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.aggregationHooks, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregationHooksResponse{
		AggregationHooks: values,
		Pagination:       pagination,
	}, nil
}

//
// Pausable Hook

func (qs queryServer) PausableHook(ctx context.Context, req *types.QueryPausableHookRequest) (*types.QueryPausableHookResponse, error) {
	hookId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	hook, err := qs.k.pausableHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	return &types.QueryPausableHookResponse{
		PausableHook: &hook,
	}, nil
}

func (qs queryServer) PausableHooks(ctx context.Context, req *types.QueryPausableHooksRequest) (*types.QueryPausableHooksResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.pausableHooks, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryPausableHooksResponse{
		PausableHooks: values,
		Pagination:    pagination,
	}, nil
}

//
// Rate Limited Hook

func (qs queryServer) RateLimitedHook(ctx context.Context, req *types.QueryRateLimitedHookRequest) (*types.QueryRateLimitedHookResponse, error) {
	hookId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, err
	}

	hook, err := qs.k.rateLimitedHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	return &types.QueryRateLimitedHookResponse{
		RateLimitedHook: &hook,
	}, nil
}

func (qs queryServer) RateLimitedHooks(ctx context.Context, req *types.QueryRateLimitedHooksRequest) (*types.QueryRateLimitedHooksResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.rateLimitedHooks, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryRateLimitedHooksResponse{
		RateLimitedHooks: values,
		Pagination:       pagination,
	}, nil
}

func (qs queryServer) TokenRateLimits(ctx context.Context, req *types.QueryTokenRateLimitsRequest) (*types.QueryTokenRateLimitsResponse, error) {
	hookId, err := util.DecodeHexAddress(req.HookId)
	if err != nil {
		return nil, err
	}

	rng := collections.NewPrefixedPairRange[uint64, []byte](hookId.GetInternalId())
	iter, err := qs.k.tokenRateLimits.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}

	tokenRateLimits, err := iter.Values()
	if err != nil {
		return nil, err
	}

	responses := make([]types.EffectiveTokenRateLimit, len(tokenRateLimits))
	for i := range tokenRateLimits {
		responses[i] = qs.effectiveTokenRateLimit(ctx, tokenRateLimits[i])
	}

	return &types.QueryTokenRateLimitsResponse{
		EffectiveTokenRateLimits: responses,
	}, nil
}

func (qs queryServer) TokenRateLimit(ctx context.Context, req *types.QueryTokenRateLimitRequest) (*types.QueryTokenRateLimitResponse, error) {
	hookId, err := util.DecodeHexAddress(req.HookId)
	if err != nil {
		return nil, err
	}

	tokenId, err := util.DecodeHexAddress(req.TokenId)
	if err != nil {
		return nil, err
	}

	tokenRateLimit, err := qs.k.tokenRateLimits.Get(ctx, types.TokenRateLimitKey(hookId, tokenId))
	if err != nil {
		return nil, err
	}

	response := qs.effectiveTokenRateLimit(ctx, tokenRateLimit)
	return &types.QueryTokenRateLimitResponse{
		EffectiveTokenRateLimit: response,
	}, nil
}

func (qs queryServer) effectiveTokenRateLimit(ctx context.Context, tokenRateLimit types.TokenRateLimit) types.EffectiveTokenRateLimit {
	return types.EffectiveTokenRateLimit{
		TokenRateLimit:    tokenRateLimit,
		CurrentLevel:      qs.k.CurrentRateLimitLevel(ctx, tokenRateLimit),
		EffectiveCapacity: tokenRateLimit.EffectiveCapacity(),
	}
}
