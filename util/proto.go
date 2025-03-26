package util

import (
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
