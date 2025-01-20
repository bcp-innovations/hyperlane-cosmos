package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) Claim(ctx context.Context, req *types.MsgClaim) (*types.MsgClaimResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	return &types.MsgClaimResponse{}, ms.k.Claim(ctx, req.Sender, igpId)
}

func (ms msgServer) CreateIGP(ctx context.Context, req *types.MsgCreateIgp) (*types.MsgCreateIgpResponse, error) {
	igpCount, err := ms.k.IgpSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(fmt.Sprintf(types.ModuleName+"/igp"), int64(igpCount))

	newIgp := types.Igp{
		Id:            prefixedId.String(),
		Owner:         req.Owner,
		Denom:         req.Denom,
		ClaimableFees: math.NewInt(0),
	}

	if err = ms.k.Igp.Set(ctx, prefixedId.Bytes(), newIgp); err != nil {
		return nil, err
	}

	return &types.MsgCreateIgpResponse{}, nil
}

func (ms msgServer) PayForGas(ctx context.Context, req *types.MsgPayForGas) (*types.MsgPayForGasResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	return &types.MsgPayForGasResponse{}, ms.k.PayForGas(ctx, req.Sender, igpId, req.MessageId, req.DestinationDomain, *req.GasLimit, *req.MaxFee)
}

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
