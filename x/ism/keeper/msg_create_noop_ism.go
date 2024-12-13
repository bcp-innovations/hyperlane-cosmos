package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/ism/types"
)

func (ms msgServer) CreateNoopIsm(ctx context.Context, req *types.MsgCreateNoopIsm) (*types.MsgCreateNoopIsmResponse, error) {
	ismCount, err := ms.k.IsmsSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName, int64(ismCount))

	newIsm := types.Ism{
		Id:      prefixedId.String(),
		IsmType: types.UNUSED,
		Creator: req.Creator,
		Ism:     &types.Ism_Noop{Noop: &types.NoopIsm{}},
	}

	if err = ms.k.Isms.Set(ctx, prefixedId.String(), newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateNoopIsmResponse{}, nil
}
