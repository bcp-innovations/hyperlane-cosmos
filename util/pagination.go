package util

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/types/query"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_QUERY_LIMIT = 1000
)

type Ranger[K any] struct {
	start *collections.RangeKey[K]
	end   *collections.RangeKey[K]
	order collections.Order
}

func (r Ranger[K]) RangeValues() (start, end *collections.RangeKey[K], order collections.Order, err error) {
	return r.start, r.end, r.order, nil
}

// GetPaginatedPrefixFromMap retrieves paginated entries from a collections.Map where keys are pairs [K1, K2]
// and values are of type T. It filters entries by a prefix of type K1 and returns them according to
// pagination parameters.
//
// Parameters:
// - ctx: Context for the operation
// - collection: The collections.Map containing the data, with Pair[K1, K2] keys and T values
// - pagination: Query pagination parameters
// - prefix: The K1 type prefix to filter entries by
//
// Returns:
// - []T: Slice of paginated values matching the prefix
// - *query.PageResponse: Pagination metadata including the next key if more results exist
// - error: Any error encountered during the operation
//
// The function handles pagination in the following ways:
// 1. If no pagination is provided, defaults to counting total with default limit
// 2. key-based pagination can be used, offset is not supported
// 3. For key-based pagination, decodes the provided key and uses it as start/end boundary
// 4. In reverse order, results are returned in descending order
// 5. Returns nextKey in PageResponse if more results exist beyond the current page
//
// Note: This implementation assumes the collection's KeyCodec can properly handle Pair[K1, K2] keys
func GetPaginatedPrefixFromMap[T any, K1 any, K2 any](ctx context.Context, collection collections.Map[collections.Pair[K1, K2], T], pagination *query.PageRequest, prefix K1) ([]T, *query.PageResponse, error) {
	// Parse basic pagination
	if pagination == nil {
		pagination = &query.PageRequest{}
	}

	key := pagination.Key
	limit := pagination.Limit
	reverse := pagination.Reverse

	if limit == 0 {
		limit = query.DefaultLimit
	}

	if limit > MAX_QUERY_LIMIT {
		return nil, nil, status.Errorf(codes.InvalidArgument, "max limit of %v exceeded", MAX_QUERY_LIMIT)
	}

	if pagination.Offset != 0 {
		return nil, nil, status.Error(codes.InvalidArgument, "offset is not supported")
	}

	pageResponse := query.PageResponse{}

	ordering := collections.OrderAscending
	start := collections.RangeKeyExact(collections.PairPrefix[K1, K2](prefix))
	end := collections.RangeKeyPrefixEnd(collections.PairPrefix[K1, K2](prefix))

	if reverse {
		ordering = collections.OrderDescending
	}

	if len(key) != 0 {
		// decode the prefixed key in the pagination
		codec := collection.KeyCodec()
		_, decodedKey, err := codec.Decode(key)
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, "failed to decode pagination key: %v", err)
		}

		// if the query is reverse we want to only get the items before the key (key becomes end)
		// otherwise we want to get items after the key (key becomes start)
		if reverse {
			end = collections.RangeKeyExact(decodedKey)
		} else {
			start = collections.RangeKeyExact(decodedKey)
		}

	}

	rng := Ranger[collections.Pair[K1, K2]]{
		start: start,
		end:   end,
		order: ordering,
	}

	it, err := collection.Iterate(ctx, rng)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create iterator: %v", err)
	}

	defer it.Close()

	data := make([]T, 0, limit)
	var keyValues []collections.KeyValue[collections.Pair[K1, K2], T]

	index := uint64(0)

	// if the query uses `key` instead of `offset` we can just fetch enough values until the limit is hit
	for ; index < limit && it.Valid(); index++ {
		keyValue, err := it.KeyValue()
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, "failed to retrieve item: %v", err)
		}
		keyValues = append(keyValues, keyValue)
		data = append(data, keyValue.Value)
		it.Next()
	}

	if it.Valid() {
		var key collections.Pair[K1, K2]
		// when the query is in reverse we want to pass the chronological last element as the next key
		// the last element will be at index 0 in that case because the order is descending
		if reverse {
			key = keyValues[len(keyValues)-1].Key
		} else {
			// the current key is the last key that is not in the response, meaning it would be the first key in the next page
			currentKey, err := it.Key()
			if err != nil {
				return nil, nil, status.Errorf(codes.Internal, "failed to retrieve key: %v", err)
			}
			key = currentKey
		}
		codec := collection.KeyCodec()
		buffer := make([]byte, codec.Size(key))
		_, err := codec.Encode(buffer, key)
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, "failed to encode next key: %v", err)
		}
		pageResponse.NextKey = buffer
	}

	return data, &pageResponse, nil
}

// GetPaginatedFromMap retrieves paginated entries from a collections.Map where keys are of type K
// and values are of type T. Returns them according to pagination parameters.
//
// Parameters:
// - ctx: Context for the operation
// - collection: The collections.Map containing the data, with K keys and T values
// - pagination: Query pagination parameters
//
// Returns:
// - []T: Slice of paginated values
// - *query.PageResponse: Pagination metadata including the next key if more results exist
// - error: Any error encountered during the operation
//
// The function handles pagination in the following ways:
// 1. If no pagination is provided, defaults to counting total with default limit
// 2. key-based pagination can be used, offset is not supported
// 3. For key-based pagination, uses provided key as start/end boundary
// 4. In reverse order, results are returned in descending order
// 5. Returns nextKey in PageResponse if more results exist beyond the current page
//
// Note: This is a simpler version of GetPaginatedPrefixFromMap that works with any key type K
func GetPaginatedFromMap[T any, K any](ctx context.Context, collection collections.Map[K, T], pagination *query.PageRequest) ([]T, *query.PageResponse, error) {
	if pagination == nil {
		pagination = &query.PageRequest{}
	}

	key := pagination.Key
	limit := pagination.Limit
	reverse := pagination.Reverse

	if limit == 0 {
		limit = query.DefaultLimit
	}

	if limit > MAX_QUERY_LIMIT {
		return nil, nil, status.Errorf(codes.InvalidArgument, "max limit of %v exceeded", MAX_QUERY_LIMIT)
	}

	if pagination.Offset != 0 {
		return nil, nil, status.Error(codes.InvalidArgument, "offset is not supported")
	}

	pageResponse := query.PageResponse{}

	ordering := collections.OrderAscending
	var end []byte = nil

	// if the query is reverse we want to only get the items before the key (key becomes end)
	// otherwise we want to get items after the key (key becomes start)
	if reverse {
		ordering = collections.OrderDescending
		end = key
		key = nil
	}

	it, err := collection.IterateRaw(ctx, key, end, ordering)
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, "failed to create iterator: %v", err)
	}

	defer it.Close()

	data := make([]T, 0, limit)
	var keyValues []collections.KeyValue[K, T]

	index := uint64(0)

	// if the query uses `key` instead of `offset` we can just fetch enough values until the limit is hit
	for ; index < limit && it.Valid(); index++ {
		keyValue, err := it.KeyValue()
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, "failed to retrieve item: %v", err)
		}
		keyValues = append(keyValues, keyValue)
		data = append(data, keyValue.Value)
		it.Next()
	}

	if it.Valid() {
		var key K
		// when the query is in reverse we want to pass the chronological last element as the next key
		// the last element will be at index 0 in that case because the order is descending
		if reverse {
			key = keyValues[len(keyValues)-1].Key
		} else {
			// the current key is the last key that is not in the response, meaning it would be the first key in the next page
			currentKey, err := it.Key()
			if err != nil {
				return nil, nil, status.Errorf(codes.Internal, "failed to retrieve key: %v", err)
			}
			key = currentKey
		}

		codec := collection.KeyCodec()
		buffer := make([]byte, codec.Size(key))
		_, err := codec.Encode(buffer, key)
		if err != nil {
			return nil, nil, status.Errorf(codes.Internal, "failed to encode next key: %v", err)
		}
		pageResponse.NextKey = buffer
	}

	return data, &pageResponse, nil
}
