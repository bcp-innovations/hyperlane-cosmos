package keeper

import (
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
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

	nextId, err := ms.k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER)
	if err != nil {
		return nil, err
	}

	newIgp := types.InterchainGasPaymaster{
		InternalId:    nextId.GetInternalId(),
		Id:            nextId.String(),
		Owner:         req.Owner,
		Denom:         req.Denom,
		ClaimableFees: math.NewInt(0),
	}

	if err = ms.k.igps.Set(ctx, newIgp.InternalId, newIgp); err != nil {
		return nil, err
	}

	return &types.MsgCreateIgpResponse{Id: nextId.String()}, nil
}

func (ms msgServer) SetIgpOwner(ctx context.Context, req *types.MsgSetIgpOwner) (*types.MsgSetIgpOwnerResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	igp, err := ms.k.igps.Get(ctx, igpId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find igp with id: %v", igpId.String())
	}

	if igp.Owner != req.Owner {
		return nil, fmt.Errorf("%s does not own igp with id %s", req.Owner, igpId.String())
	}

	// TODO: Verfiy NewOwner

	igp.Owner = req.NewOwner

	if err = ms.k.igps.Set(ctx, igpId.GetInternalId(), igp); err != nil {
		return nil, err
	}

	return &types.MsgSetIgpOwnerResponse{}, nil
}

// PayForGas executes an InterchainGasPayment without for the specified payment amount.
func (ms msgServer) PayForGas(ctx context.Context, req *types.MsgPayForGas) (*types.MsgPayForGasResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, fmt.Errorf("ism id %s is invalid: %s", req.IgpId, err.Error())
	}

	handler := InterchainGasPaymasterHookHandler{*ms.k}

	return &types.MsgPayForGasResponse{}, handler.PayForGasWithoutQuote(ctx, igpId, req.Sender, req.MessageId, req.DestinationDomain, req.GasLimit, req.Amount)
}

func (ms msgServer) SetDestinationGasConfig(ctx context.Context, req *types.MsgSetDestinationGasConfig) (*types.MsgSetDestinationGasConfigResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, fmt.Errorf("ism id %s is invalid: %s", req.IgpId, err.Error())
	}

	return &types.MsgSetDestinationGasConfigResponse{}, ms.k.SetDestinationGasConfig(ctx, igpId, req.Owner, req.DestinationGasConfig)
}
