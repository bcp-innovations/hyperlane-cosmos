package keeper

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

func InitGenesis(ctx sdk.Context, k Keeper, data *types.GenesisState) {
	if data == nil || data.Isms == nil {
		return
	}

	isms, err := util.UnpackAnys[types.HyperlaneInterchainSecurityModule](data.Isms)
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

	// TODO: save & load router in genesis state
	// if err := k.ismsSequence.Set(ctx, uint64(len(isms))); err != nil {
	// 	panic(err)
	// }
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

	msgs := make([]proto.Message, len(isms))
	for i, ism := range isms {
		msgs[i] = ism
	}
	ismsAny, err := util.PackAnys(msgs)
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Isms: ismsAny,
	}
}
