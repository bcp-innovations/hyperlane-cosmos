package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

func (ms msgServer) Claim(ctx context.Context, req *types.MsgClaim) (*types.MsgClaimResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, fmt.Errorf("ism id %s is invalid: %s", req.IgpId, err.Error())
	}

	return &types.MsgClaimResponse{}, ms.k.Claim(ctx, req.Sender, igpId)
}

func (ms msgServer) CreateIgp(ctx context.Context, req *types.MsgCreateIgp) (*types.MsgCreateIgpResponse, error) {
	err := sdk.ValidateDenom(req.Denom)
	if err != nil {
		return nil, fmt.Errorf("denom %s is invalid", req.Denom)
	}

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

	return &types.MsgCreateIgpResponse{Id: prefixedId.String()}, nil
}

// PayForGas executes an InterchainGasPayment without for the specified payment amount.
func (ms msgServer) PayForGas(ctx context.Context, req *types.MsgPayForGas) (*types.MsgPayForGasResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, fmt.Errorf("ism id %s is invalid: %s", req.IgpId, err.Error())
	}

	return &types.MsgPayForGasResponse{}, ms.k.PayForGasWithoutQuote(ctx, req.Sender, igpId, req.MessageId, req.DestinationDomain, req.GasLimit, req.Amount)
}

func (ms msgServer) SetDestinationGasConfig(ctx context.Context, req *types.MsgSetDestinationGasConfig) (*types.MsgSetDestinationGasConfigResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, fmt.Errorf("ism id %s is invalid: %s", req.IgpId, err.Error())
	}

	igp, err := ms.k.Igp.Get(ctx, igpId.Bytes())
	if err != nil {
		return nil, fmt.Errorf("igp does not exist: %s", igpId.String())
	}

	if igp.Owner != req.Owner {
		return nil, fmt.Errorf("failed to set DestinationGasConfigs: %s is not the owner of igp with id %s", req.Owner, igpId.String())
	}

	if req.DestinationGasConfig.GasOracle == nil {
		return nil, fmt.Errorf("failed to set DestinationGasConfigs: gas Oracle is required")
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
