package keeper

import (
	"bytes"
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (ms msgServer) AnnounceValidator(ctx context.Context, req *types.MsgAnnounceValidator) (*types.MsgAnnounceValidatorResponse, error) {
	validatorKey, err := util.DecodeEthHex(req.Validator)
	if err != nil {
		return nil, err
	}

	// Ensure that validator hasn't already announced storage location.
	prefixedId := util.CreateValidatorStorageKey(validatorKey)

	exists, err := ms.k.Validators.Has(ctx, prefixedId.Bytes())
	if err != nil {
		return nil, err
	}

	var validator types.Validator

	if exists {
		validator, err = ms.k.Validators.Get(ctx, prefixedId.Bytes())
		if err != nil {
			return nil, err
		}

		for _, location := range validator.StorageLocations {
			if location == req.StorageLocation {
				return nil, fmt.Errorf("validator %s already announced storage location %s", req.Validator, req.StorageLocation)
			}
		}

		validator.StorageLocations = append(validator.StorageLocations, req.StorageLocation)
	} else {
		validator = types.Validator{
			Address:          util.EncodeEthHex(validatorKey),
			StorageLocations: []string{req.StorageLocation},
		}
	}

	sig, err := util.DecodeEthHex(req.Signature)
	if err != nil {
		return nil, err
	}

	mailboxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, err
	}

	localDomain, err := ms.k.LocalDomain(ctx)
	if err != nil {
		return nil, err
	}

	announcementDigest := types.GetAnnouncementDigest(req.StorageLocation, localDomain, mailboxId.Bytes())
	ethSigningHash := util.GetEthSigningHash(announcementDigest[:])

	recoveredPubKey, err := util.RecoverEthSignature(ethSigningHash[:], sig)
	if err != nil {
		return nil, err
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)

	if !bytes.Equal(recoveredAddress[:], validatorKey) {
		return nil, fmt.Errorf("validator %s doesn't match signature. recovered address: %s", util.EncodeEthHex(validatorKey), util.EncodeEthHex(recoveredAddress[:]))
	}

	if err = ms.k.Validators.Set(ctx, prefixedId.Bytes(), validator); err != nil {
		return nil, err
	}

	return &types.MsgAnnounceValidatorResponse{}, nil
}
