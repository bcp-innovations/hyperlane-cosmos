package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/collections"

	"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
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
	mailboxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, err
	}

	validatorAddress, err := util.DecodeEthHex(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	rng := collections.NewSuperPrefixedTripleRange[[]byte, []byte, uint64](mailboxId.Bytes(), validatorAddress)

	iter, err := qs.k.storageLocations.Iterate(ctx, rng)
	if err != nil {
		return nil, err
	}

	storageLocations, err := iter.Values()
	if err != nil {
		return nil, err
	}

	return &types.QueryAnnouncedStorageLocationsResponse{
		StorageLocations: storageLocations,
	}, nil
}

func (qs queryServer) LatestAnnouncedStorageLocation(ctx context.Context, req *types.QueryLatestAnnouncedStorageLocationRequest) (*types.QueryLatestAnnouncedStorageLocationResponse, error) {
	mailboxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, err
	}

	validatorAddress, err := util.DecodeEthHex(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	// encode the prefix key so we can set the order by ourselves
	key := collections.TripleSuperPrefix[[]byte, []byte, uint64](mailboxId.Bytes(), validatorAddress)
	codec := qs.k.storageLocations.KeyCodec()
	start := make([]byte, codec.Size(key))
	_, err = codec.Encode(start, key)
	if err != nil {
		return nil, err
	}

	// create a new iterator that is in reverse order
	// meaning that the first item will be the latest location
	iter, err := qs.k.storageLocations.IterateRaw(ctx, start, nil, collections.OrderDescending)
	if err != nil {
		return nil, err
	}

	location, err := iter.Value()

	return &types.QueryLatestAnnouncedStorageLocationResponse{
		StorageLocation: location,
	}, err
}

// Isms returns all ism which are registered in this module.
func (qs queryServer) Isms(ctx context.Context, req *types.QueryIsmsRequest) (*types.QueryIsmsResponse, error) {
	values, pagination, err := util.GetPaginatedFromMap(ctx, qs.k.isms, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	msgs := make([]proto.Message, len(values))
	for i, value := range values {
		msgs[i] = value
	}

	isms, err := util.PackAnys(msgs)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryIsmsResponse{
		Isms:       isms,
		Pagination: pagination,
	}, nil
}

// Ism returns a single ISM given by its id
func (qs queryServer) Ism(ctx context.Context, req *types.QueryIsmRequest) (*types.QueryIsmResponse, error) {
	ismId, err := util.DecodeHexAddress(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid hex address %s, %s", req.Id, err.Error())
	}

	ism, err := qs.k.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ism %s not found", req.Id)
	}

	toAny, err := util.PackAny(ism)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryIsmResponse{
		Ism: *toAny,
	}, nil
}
