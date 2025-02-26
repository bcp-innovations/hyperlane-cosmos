package keeper

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/query"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
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

	payment, err := igpHandler.QuoteGasPayment(ctx, igpId, uint32(destinationDomain), gasLimit)
	if err != nil {
		return nil, err
	}

	return &types.QueryQuoteGasPaymentResponse{GasPayment: payment.String()}, nil
}

func (qs queryServer) Igps(ctx context.Context, req *types.QueryIgpsRequest) (*types.QueryIgpsResponse, error) {
	values, pagination, err := GetPaginatedFromMap(ctx, qs.k.Igps, req.Pagination)
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

// TODO remove once PR (pagination is merged)
func GetPaginatedFromMap[T any, K any](ctx context.Context, collection collections.Map[K, T], pagination *query.PageRequest) ([]T, *query.PageResponse, error) {
	// Parse basic pagination
	if pagination == nil {
		pagination = &query.PageRequest{CountTotal: true}
	}

	offset := pagination.Offset
	key := pagination.Key
	limit := pagination.Limit
	reverse := pagination.Reverse

	if limit == 0 {
		limit = query.DefaultLimit
	}

	pageResponse := query.PageResponse{}

	// user has to use either offset or key, not both
	if offset > 0 && key != nil {
		return nil, nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	ordering := collections.OrderDescending
	if reverse {
		ordering = collections.OrderAscending
	}

	// TODO: subject to change -> use it as key so we can jump to the offset directly
	it, err := collection.IterateRaw(ctx, key, nil, ordering)
	if err != nil {
		return nil, nil, err
	}

	defer it.Close()

	data := make([]T, 0, limit)
	keyValues, err := it.KeyValues()
	if err != nil {
		return nil, nil, err
	}
	length := uint64(len(keyValues))

	i := uint64(offset)
	for ; i < limit+offset && i < length; i++ {
		data = append(data, keyValues[i].Value)
	}

	if i < length {
		encodedKey := keyValues[i].Key
		codec := collection.KeyCodec()
		buffer := make([]byte, codec.Size(encodedKey))
		_, err := codec.Encode(buffer, encodedKey)
		if err != nil {
			return nil, nil, err
		}
		pageResponse.NextKey = buffer
	}

	return data, &pageResponse, nil
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
