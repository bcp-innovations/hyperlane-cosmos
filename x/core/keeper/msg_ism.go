package keeper

import (
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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

	if len(ism.ValidatorPubKeys) < int(ism.Threshold) {
		return fmt.Errorf("validator pubkeys less than threshold")
	}

	for _, validatorPubKey := range ism.ValidatorPubKeys {
		pubKey, err := util.DecodeEthHex(validatorPubKey)
		if err != nil {
			return fmt.Errorf("invalid validator pub key: %s", validatorPubKey)
		}

		_, err = crypto.UnmarshalPubkey(pubKey)
		if err != nil {
			return fmt.Errorf("invalid validator pub key: %s", validatorPubKey)
		}
	}
	return nil
}
