package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	isms         collections.Map[uint64, types.HyperlaneInterchainSecurityModule]
	ismsSequence collections.Sequence
	schema       collections.Schema

	hexAddressFactory util.HexAddressFactory
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	factory, err := util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER)
	println("FACTORY CREATED: ISM")
	if err != nil {
		panic(err)
	}

	k := Keeper{
		isms:              collections.NewMap(sb, types.IsmsKey, "isms", collections.Uint64Key, codec.CollInterfaceValue[types.HyperlaneInterchainSecurityModule](cdc)),
		ismsSequence:      collections.NewSequence(sb, types.IsmsSequenceKey, "isms_sequence"),
		hexAddressFactory: factory,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k Keeper) Verify(ctx sdk.Context, ismId util.HexAddress, message util.HyperlaneMessage, metadata coreTypes.Metadata) (bool, error) {
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
