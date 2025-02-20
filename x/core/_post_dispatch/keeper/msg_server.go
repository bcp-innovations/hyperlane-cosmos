package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (m msgServer) CreateInterchainGasPaymaster(ctx context.Context, req *types.MsgCreateInterchainGasPaymaster) (*types.MsgCreateInterchainGasPaymasterResponse, error) {
	err := sdk.ValidateDenom(req.Denom)
	if err != nil {
		return nil, fmt.Errorf("denom %s is invalid", req.Denom)
	}

	igpCount, err := m.k.interchainGasPaymastersSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	newIgp := types.InterchainGasPaymaster{
		Id:            igpCount,
		Owner:         req.Creator,
		Denom:         req.Denom,
		ClaimableFees: math.NewInt(0),
	}

	if err = m.k.interchainGasPaymasters.Set(ctx, igpCount, newIgp); err != nil {
		return nil, err
	}

	hexAddress := m.k.hexAddressFactory.GenerateId(uint32(types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER), igpCount)

	return &types.MsgCreateInterchainGasPaymasterResponse{Id: hexAddress.String()}, nil
}
