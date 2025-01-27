package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

func (ms msgServer) CreateMultisigIsm(ctx context.Context, req *types.MsgCreateMultisigIsm) (*types.MsgCreateMultisigIsmResponse, error) {
	ismCount, err := ms.k.IsmsSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName+"/ism", int64(ismCount))

	ism := types.MultiSigIsm{
		ValidatorPubKeys: req.MultiSig.ValidatorPubKeys,
		Threshold:        req.MultiSig.Threshold,
	}

	newIsm := types.Ism{
		Id:      prefixedId.String(),
		IsmType: types.MESSAGE_ID_MULTISIG,
		Creator: req.Creator,
		Ism:     &types.Ism_MultiSig{MultiSig: &ism},
	}

	if err = ms.k.Isms.Set(ctx, prefixedId.Bytes(), newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateMultisigIsmResponse{}, nil
}

func (ms msgServer) CreateNoopIsm(ctx context.Context, req *types.MsgCreateNoopIsm) (*types.MsgCreateNoopIsmResponse, error) {
	ismCount, err := ms.k.IsmsSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName+"/ism", int64(ismCount))

	newIsm := types.Ism{
		Id:      prefixedId.String(),
		IsmType: types.UNUSED,
		Creator: req.Creator,
		Ism:     &types.Ism_Noop{Noop: &types.NoopIsm{}},
	}

	if err = ms.k.Isms.Set(ctx, prefixedId.Bytes(), newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateNoopIsmResponse{}, nil
}
