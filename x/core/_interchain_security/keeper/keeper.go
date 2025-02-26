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
	isms collections.Map[uint64, types.HyperlaneInterchainSecurityModule]
	// Key: Mailbox ID, Validator address, Storage Location index
	storageLocations collections.Map[collections.Triple[[]byte, []byte, uint64], string]
	schema           collections.Schema

	coreKeeper types.CoreKeeper
	router     *util.Router[util.InterchainSecurityModule]
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		isms:             collections.NewMap(sb, types.IsmsKey, "isms", collections.Uint64Key, codec.CollInterfaceValue[types.HyperlaneInterchainSecurityModule](cdc)),
		storageLocations: collections.NewMap(sb, types.StorageLocationsKey, "storage_locations", collections.TripleKeyCodec(collections.BytesKey, collections.BytesKey, collections.Uint64Key), collections.StringValue),
		coreKeeper:       nil,
		router:           util.NewRouter[util.InterchainSecurityModule](types.RouterKey, sb),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	// add default modules
	k.router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED, NewIsmHandler(&k))
	k.router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG, NewIsmHandler(&k))

	k.schema = schema

	return k
}

func (k *Keeper) SetCoreKeeper(coreKeeper types.CoreKeeper) {
	k.coreKeeper = coreKeeper
}

func (k *Keeper) GetRouter() *util.Router[util.InterchainSecurityModule] {
	return k.router
}

func (k *Keeper) GetIsm(ctx sdk.Context, ismId uint64) (types.HyperlaneInterchainSecurityModule, error) {
	ism, err := k.isms.Get(ctx, ismId)
	if err != nil {
		return nil, err
	}

	return ism, nil
}

func (k Keeper) Verify(ctx sdk.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	id := ismId.GetInternalId() // TODO: change to uint64
	handler, err := k.router.GetModule(ctx, id)
	if err != nil {
		return false, err
	}
	return (*handler).Verify(ctx, id, metadata, message)
}

func (k Keeper) IsmIdExists(ctx context.Context, ismId util.HexAddress) bool {
	return k.router.Exists(ctx, ismId.GetInternalId()) // TODO: change to uint64
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
