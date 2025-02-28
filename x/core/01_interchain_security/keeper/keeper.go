package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

type Keeper struct {
	// TODO: move to internal id -> uint64
	isms collections.Map[[]byte, types.HyperlaneInterchainSecurityModule]
	// Key: Mailbox ID, Validator address, Storage Location index
	storageLocations collections.Map[collections.Triple[[]byte, []byte, uint64], string]
	schema           collections.Schema

	coreKeeper types.CoreKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		isms:             collections.NewMap(sb, types.IsmsKey, "isms", collections.BytesKey, codec.CollInterfaceValue[types.HyperlaneInterchainSecurityModule](cdc)),
		storageLocations: collections.NewMap(sb, types.StorageLocationsKey, "storage_locations", collections.TripleKeyCodec(collections.BytesKey, collections.BytesKey, collections.Uint64Key), collections.StringValue),
		coreKeeper:       nil,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k *Keeper) SetCoreKeeper(coreKeeper types.CoreKeeper) {
	if k.coreKeeper != nil {
		panic("core keeper already set")
	}

	k.coreKeeper = coreKeeper

	// set the router from the core keeper
	router := coreKeeper.IsmRouter()
	// add default modules
	router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED, NewIsmHandler(k))
	router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG, NewIsmHandler(k))
	router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TPYE_MESSAGE_ID_MULTISIG, NewIsmHandler(k))
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
