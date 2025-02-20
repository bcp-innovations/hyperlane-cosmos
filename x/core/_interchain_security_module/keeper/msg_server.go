package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security_module/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (m msgServer) CreateMerkleRootMultiSigIsm(ctx context.Context, ism *types.MsgCreateMerkleRootMultiSigIsm) (*types.MsgCreateMerkleRootMultiSigIsmResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) CreateNoopIsm(ctx context.Context, ism *types.MsgCreateNoopIsm) (*types.MsgCreateNoopIsmResponse, error) {
	ismCount, err := m.k.ismsSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	newIsm := types.NoopISM{
		Id:    ismCount,
		Owner: ism.Creator,
	}

	if err = m.k.isms.Set(ctx, ismCount, &newIsm); err != nil {
		return nil, err
	}

	hexAddress := m.k.hexAddressFactory.GenerateId(uint32(types.INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED), ismCount)

	return &types.MsgCreateNoopIsmResponse{Id: hexAddress.String()}, nil
}
