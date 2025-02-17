package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func (ms msgServer) AnnounceValidator(ctx context.Context, req *types.MsgAnnounceValidator) (*types.MsgAnnounceValidatorResponse, error) {
	if req.Validator == "" {
		return nil, fmt.Errorf("validator cannot be empty")
	}

	if req.StorageLocation == "" {
		return nil, fmt.Errorf("storage location cannot be empty")
	}

	if req.Signature == "" {
		return nil, fmt.Errorf("signature cannot be empty")
	}

	sig, err := util.DecodeEthHex(req.Signature)
	if err != nil {
		return nil, fmt.Errorf("invalid signature")
	}

	mailboxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, fmt.Errorf("invalid mailbox id")
	}

	found, err := ms.k.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil || !found {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", mailboxId.String())
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

	validatorAddress, err := util.DecodeEthHex(req.Validator)
	if err != nil {
		return nil, fmt.Errorf("invalid validator address")
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)

	if !bytes.Equal(recoveredAddress[:], validatorAddress) {
		return nil, fmt.Errorf("validator %s doesn't match signature. recovered address: %s", util.EncodeEthHex(validatorAddress), util.EncodeEthHex(recoveredAddress[:]))
	}

	exists, err := ms.k.Validators.Has(ctx, validatorAddress)
	if err != nil {
		return nil, err
	}

	var newStorageLocation types.StorageLocation
	if exists {
		rng := collections.NewPrefixedPairRange[[]byte, uint64](validatorAddress)

		iter, err := ms.k.StorageLocations.Iterate(ctx, rng)
		if err != nil {
			return nil, err
		}

		storageLocations, err := iter.Values()
		if err != nil {
			return nil, err
		}

		// It is assumed that a validator announces a reasonable amount of storage locations.
		// Otherwise, one would need to store the hash in a separate lookup table which adds more complexity.
		for _, location := range storageLocations {
			if location.Location == req.StorageLocation {
				return nil, fmt.Errorf("validator %s already announced storage location %s", req.Validator, req.StorageLocation)
			}
		}

		newStorageLocation = types.StorageLocation{
			Location: req.StorageLocation,
			Id:       uint64(len(storageLocations)),
		}
	} else {
		validator := types.Validator{
			Address: util.EncodeEthHex(validatorAddress),
		}

		if err = ms.k.Validators.Set(ctx, validatorAddress, validator); err != nil {
			return nil, err
		}

		newStorageLocation = types.StorageLocation{
			Location: req.StorageLocation,
			Id:       0,
		}
	}

	if err = ms.k.StorageLocations.Set(ctx, collections.Join(validatorAddress, newStorageLocation.Id), newStorageLocation); err != nil {
		return nil, err
	}

	return &types.MsgAnnounceValidatorResponse{}, nil
}
