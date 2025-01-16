package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
)

type msgServer struct {
	k Keeper
}

func (ms msgServer) SetBeneficiary(ctx context.Context, req *types.MsgSetBeneficiary) (*types.MsgSetBeneficiaryResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}
