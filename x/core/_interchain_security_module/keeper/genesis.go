package keeper

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security_module/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

func InitGenesis(ctx sdk.Context, k Keeper, data *types.GenesisState) {
	isms, err := unpackAccounts(data.Isms)
	if err != nil {
		panic(err)
	}

	for index, ism := range isms {
		if uint64(index) != ism.GetId() {
			panic(fmt.Sprintf("unexpected ism %d; expected %d", index, ism.GetId()))
		}
		if err := k.isms.Set(ctx, ism.GetId(), ism); err != nil {
			panic(err)
		}
	}

	if err := k.ismsSequence.Set(ctx, uint64(len(isms))); err != nil {
		panic(err)
	}
}

func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	iter, err := k.isms.Iterate(ctx, nil)
	if err != nil {
		panic(err)
	}

	isms, err := iter.Values()
	if err != nil {
		panic(err)
	}

	ismsAny, err := packAccounts(isms)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Isms: ismsAny,
	}
}

func packAccounts(isms []types.HyperlaneInterchainSecurityModule) ([]*codectypes.Any, error) {
	ismsAny := make([]*codectypes.Any, len(isms))
	for i, acc := range isms {
		msg, ok := acc.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("cannot proto marshal %T", acc)
		}
		anyProto, err := codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
		ismsAny[i] = anyProto
	}

	return ismsAny, nil
}

// UnpackAccounts converts Any slice to GenesisAccounts
func unpackAccounts(ismsAny []*codectypes.Any) ([]types.HyperlaneInterchainSecurityModule, error) {
	accounts := make([]types.HyperlaneInterchainSecurityModule, len(ismsAny))
	for i, anyProto := range ismsAny {
		acc, ok := anyProto.GetCachedValue().(types.HyperlaneInterchainSecurityModule)
		if !ok {
			return nil, fmt.Errorf("expected genesis ism")
		}
		accounts[i] = acc
	}

	return accounts, nil
}
