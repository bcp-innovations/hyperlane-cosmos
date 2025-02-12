package keeper

import (
	"context"
	"fmt"

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
		Validators: req.MultiSig.Validators,
		Threshold:  req.MultiSig.Threshold,
	}

	if err = validateMultisigIsm(ism); err != nil {
		return nil, err
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

	return &types.MsgCreateMultisigIsmResponse{Id: prefixedId.String()}, nil
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

	return &types.MsgCreateNoopIsmResponse{Id: prefixedId.String()}, nil
}

func validateMultisigIsm(ism types.MultiSigIsm) error {
	if ism.Threshold == 0 {
		return fmt.Errorf("threshold must be greater than zero")
	}

	if len(ism.Validators) < int(ism.Threshold) {
		return fmt.Errorf("validator address less than threshold")
	}

	for _, validatorAddress := range ism.Validators {
		bytes, err := util.DecodeEthHex(validatorAddress)
		if err != nil {
			return fmt.Errorf("invalid validator address: %s", validatorAddress)
		}

		// ensure that the address is an eth address with 20 bytes
		if len(bytes) != 20 {
			return fmt.Errorf("invalid validator address: must be ethereum address (20 byts)")
		}
	}

	// check for duplications
	count := map[string]int{}
	for _, address := range ism.Validators {
		count[address]++
		if count[address] > 1 {
			return fmt.Errorf("duplicate validator address: %v", address)
		}
	}
	return nil
}
