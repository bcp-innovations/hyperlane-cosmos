package util

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func encodeCollectionKey[K any](key K, codec codec.KeyCodec[K]) ([]byte, error) {
	buffer := make([]byte, codec.Size(key))
	_, err := codec.Encode(buffer, key)
	return buffer, err
}

func decodeCollectionKey[K any](buffer []byte, codec codec.KeyCodec[K]) (K, error) {
	_, key, err := codec.Decode(buffer)
	return key, err
}

func GetPaginatedFromMap[T any](ctx context.Context, collection collections.Map[uint64, T], pagination *query.PageRequest, sequence collections.Sequence) ([]T, *query.PageResponse, error) {
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
		return nil, nil, fmt.Errorf("invalid request: either offset or key is expected, got both")
	}

	// get the current end of the collection
	current_sequence, err := sequence.Peek(ctx)
	if err != nil {
		return nil, nil, err
	}
	if pagination.CountTotal {
		// the current sequence is equal to the total available items
		pageResponse.Total = current_sequence
	}

	// startKey is either the key or the offset
	startKey := uint64(0)
	if key != nil {
		startKey, err = decodeCollectionKey(key, collection.KeyCodec())
		if err != nil {
			return nil, nil, err
		}
	} else {
		// start key equals to the offset
		startKey = offset
	}

	// end is the start + limit
	endKey := startKey + limit

	rng := (&collections.Range[uint64]{}).StartInclusive(startKey).EndExclusive(endKey)
	if reverse {
		rng = rng.Descending()
	}

	it, err := collection.Iterate(ctx, rng)
	if err != nil {
		return nil, nil, err
	}

	defer it.Close()

	values, err := it.Values()
	if err != nil {
		return nil, nil, err
	}

	// items are less than the limit meaning we reached the end
	if len(values) < int(limit) || endKey >= current_sequence {
		pageResponse.NextKey = nil
		return values, &pageResponse, nil
	}

	// items are equal to the limit meaning we need to paginate
	encodedKey, err := encodeCollectionKey(endKey, collection.KeyCodec())
	if err != nil {
		return nil, nil, err
	}
	pageResponse.NextKey = encodedKey

	return values, &pageResponse, nil
}
