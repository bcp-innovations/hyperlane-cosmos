package keeper

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
	"github.com/ethereum/go-ethereum/crypto"
	"slices"
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

		for _, location := range validator.StorageLocation {
			if location == req.StorageLocation {
				return nil, fmt.Errorf("validator %s already announced storage location %s", req.Validator, req.StorageLocation)
			}
		}

		validator.StorageLocation = append(validator.StorageLocation, req.StorageLocation)
	} else {
		validator = types.Validator{
			Address:         util.EncodeEthHex(validatorKey),
			StorageLocation: []string{req.StorageLocation},
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

	announcementDigest := getAnnouncementDigest(req.StorageLocation, ms.k.LocalDomain(), mailboxId.Bytes())

	recoveredPubKey, err := crypto.SigToPub(announcementDigest, sig)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(crypto.FromECDSAPub(recoveredPubKey), validatorKey) {
		return nil, fmt.Errorf("validator %s doesn't match signature", util.EncodeEthHex(validatorKey))
	}

	if err = ms.k.Validators.Set(ctx, prefixedId.Bytes(), validator); err != nil {
		return nil, err
	}

	return &types.MsgAnnounceValidatorResponse{}, nil
}

func getAnnouncementDigest(storageLocation string, domainId uint32, mailbox []byte) []byte {
	var domainHashBytes []byte

	domainIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(domainIdBytes, domainId)

	// TODO: Check if all of them are required
	domainHashBytes = slices.Concat(
		domainIdBytes,
		mailbox,
		[]byte("HYPERLANE_ANNOUNCEMENT"),
	)

	domainHash := crypto.Keccak256Hash(domainHashBytes)

	announcementDigestBytes := slices.Concat(
		domainHash.Bytes(),
		[]byte(storageLocation),
	)

	return crypto.Keccak256(announcementDigestBytes)
}
