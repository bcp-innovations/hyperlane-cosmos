package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
)

func (ms msgServer) SetDestinationGasConfig(ctx context.Context, req *types.MsgSetDestinationGasConfig) (*types.MsgSetDestinationGasConfigResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	igp, err := ms.k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return nil, err
	}

	if igp.Owner != req.Owner {
		return nil, fmt.Errorf("failed to set DestinationGasConfigs: %s is not the owner of IGP with id %s", req.Owner, igpId.String())
	}

	updatedDestinationGasConfig := types.DestinationGasConfig{
		RemoteDomain: req.DestinationGasConfig.RemoteDomain,
		GasOracle:    req.DestinationGasConfig.GasOracle,
		GasOverhead:  req.DestinationGasConfig.GasOverhead,
	}

	key := collections.Join(igpId.Bytes(), req.DestinationGasConfig.RemoteDomain)

	err = ms.k.IgpDestinationGasConfigMap.Set(ctx, key, updatedDestinationGasConfig)
	if err != nil {
		return nil, err
	}

	return &types.MsgSetDestinationGasConfigResponse{}, nil
}
