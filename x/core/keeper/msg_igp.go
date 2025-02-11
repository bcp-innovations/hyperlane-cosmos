package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

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

	return &types.MsgSetDestinationGasConfigResponse{}, ms.k.SetDestinationGasConfig(ctx, igpId, req.Owner, req.DestinationGasConfig)
}
