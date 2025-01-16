package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
)

func (ms msgServer) CreateIGP(ctx context.Context, req *types.MsgCreateIgp) (*types.MsgCreateIgpResponse, error) {
	igpCount, err := ms.k.IgpSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName, int64(igpCount))

	newIgp := types.Igp{
		Id:                    prefixedId.String(),
		Owner:                 req.Owner,
		Beneficiary:           req.Beneficiary,
		DestinationGasConfigs: []*types.DestinationGasConfig{},
	}

	if err = ms.k.Igp.Set(ctx, prefixedId.Bytes(), newIgp); err != nil {
		return nil, err
	}

	return &types.MsgCreateIgpResponse{}, nil
}
