package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	isms         collections.Map[uint64, types.HyperlaneInterchainSecurityModule]
	ismsSequence collections.Sequence
	// Key: Mailbox ID, Validator address, Storage Location index
	storageLocations collections.Map[collections.Triple[[]byte, []byte, uint64], string]
	schema           collections.Schema

	coreKeeper types.CoreKeeper

	hexAddressFactory util.HexAddressFactory
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	factory, err := util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER)
	if err != nil {
		panic(err)
	}

	k := Keeper{
		isms:              collections.NewMap(sb, types.IsmsKey, "isms", collections.Uint64Key, codec.CollInterfaceValue[types.HyperlaneInterchainSecurityModule](cdc)),
		ismsSequence:      collections.NewSequence(sb, types.IsmsSequenceKey, "isms_sequence"),
		storageLocations:  collections.NewMap(sb, types.StorageLocationsKey, "storage_locations", collections.TripleKeyCodec(collections.BytesKey, collections.BytesKey, collections.Uint64Key), collections.StringValue),
		coreKeeper:        nil,
		hexAddressFactory: factory,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k *Keeper) SetCoreKeeper(coreKeeper types.CoreKeeper) {
	k.coreKeeper = coreKeeper
}

func (k Keeper) Verify(ctx sdk.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	// Global Conventions
	// - Address must be unique
	// - Hook must check if id exists (and correct recipient)
	// module_name / class / type / custom

	if !k.hexAddressFactory.IsClassMember(ismId) {
		return false, nil
	}

	ism, err := k.isms.Get(ctx, 0)
	if err != nil {
		return false, err
	}

	return ism.Verify(ctx, metadata, message)
}

func (k Keeper) IsmIdExists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	exists, err := k.isms.Has(ctx, ismId.GetInternalId())
	if err != nil {
		return false, err
	}
	return exists, nil
}

// TODO outsource to utils class, once migrated
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
