package util

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/gogoproto/proto"
)

func PackAny(item proto.Message) (*codectypes.Any, error) {
	anyProto, err := codectypes.NewAnyWithValue(item)
	if err != nil {
		return nil, err
	}
	return anyProto, nil
}

func UnpackAny[T any](anyProto *codectypes.Any) (*T, error) {
	item, ok := anyProto.GetCachedValue().(T)
	if !ok {
		return nil, fmt.Errorf("cannot cast %T", anyProto)
	}
	return &item, nil
}

func PackAnys(isms []proto.Message) ([]*codectypes.Any, error) {
	ismsAny := make([]*codectypes.Any, len(isms))
	for i, acc := range isms {
		anyProto, err := PackAny(acc)
		if err != nil {
			return nil, err
		}
		ismsAny[i] = anyProto
	}

	return ismsAny, nil
}

func UnpackAnys[T any](anys []*codectypes.Any) ([]T, error) {
	items := make([]T, len(anys))
	for i, anyProto := range anys {
		item, err := UnpackAny[T](anyProto)
		if err != nil {
			return nil, err
		}
		items[i] = *item
	}

	return items, nil
}
