package keeper

import (
	"context"
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

	// TODO: Find a better way to store destination configs
	for i := range igp.DestinationGasConfigs {
		if igp.DestinationGasConfigs[i].RemoteDomain == req.DestinationGasConfig.RemoteDomain {
			igp.DestinationGasConfigs[i] = &updatedDestinationGasConfig

			if err = ms.k.Igp.Set(ctx, igpId.Bytes(), igp); err != nil {
				return nil, err
			}

			return &types.MsgSetDestinationGasConfigResponse{}, nil
		}
	}

	if len(igp.DestinationGasConfigs) >= types.MaxDestinationGasConfigs {
		return nil, fmt.Errorf("max DestinationGasConfigs are reached: %v", types.MaxDestinationGasConfigs)
	}

	igp.DestinationGasConfigs = append(igp.DestinationGasConfigs, &updatedDestinationGasConfig)

	if err = ms.k.Igp.Set(ctx, igpId.Bytes(), igp); err != nil {
		return nil, err
	}

	return &types.MsgSetDestinationGasConfigResponse{}, nil
}
