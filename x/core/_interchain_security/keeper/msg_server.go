package keeper

import (
	"bytes"
	"context"
	"fmt"

	"cosmossdk.io/collections"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (m msgServer) AnnounceValidator(ctx context.Context, req *types.MsgAnnounceValidator) (*types.MsgAnnounceValidatorResponse, error) {
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

	found, err := m.k.coreKeeper.MailboxIdExists(ctx, mailboxId)
	if err != nil || !found {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", mailboxId.String())
	}

	localDomain, err := m.k.coreKeeper.LocalDomain(ctx)
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

	// Check if validator already exists.
	exists, err := m.k.storageLocations.Has(ctx, collections.Join3(mailboxId.Bytes(), validatorAddress, uint64(0)))
	if err != nil {
		return nil, err
	}

	var storageLocationIndex uint64 = 0
	if exists {
		rng := collections.NewSuperPrefixedTripleRange[[]byte, []byte, uint64](mailboxId.Bytes(), validatorAddress)

		iter, err := m.k.storageLocations.Iterate(ctx, rng)
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
			if location == req.StorageLocation {
				return nil, fmt.Errorf("validator %s already announced storage location %s", req.Validator, req.StorageLocation)
			}
		}
		storageLocationIndex = uint64(len(storageLocations))
	}

	if err = m.k.storageLocations.Set(ctx, collections.Join3(mailboxId.Bytes(), validatorAddress, storageLocationIndex), req.StorageLocation); err != nil {
		return nil, err
	}

	return &types.MsgAnnounceValidatorResponse{}, nil
}

func (m msgServer) CreateMessageIdMultisigIsm(ctx context.Context, req *types.MsgCreateMessageIdMultisigIsm) (*types.MsgCreateMessageIdMultisigIsmResponse, error) {
	ismId, err := m.k.router.GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TPYE_MESSAGE_ID_MULTISIG)
	if err != nil {
		return nil, err
	}

	newIsm := types.MessageIdMultisigISM{
		Id:         ismId,
		Owner:      req.Creator,
		Validators: req.Validators,
		Threshold:  req.Threshold,
	}

	if err = newIsm.Validate(); err != nil {
		return nil, err
	}

	if err = m.k.isms.Set(ctx, ismId, &newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateMessageIdMultisigIsmResponse{Id: fmt.Sprintf("%v", ismId)}, nil
}

func (m msgServer) CreateMerkleRootMultisigIsm(ctx context.Context, req *types.MsgCreateMerkleRootMultisigIsm) (*types.MsgCreateMerkleRootMultisigIsmResponse, error) {
	ismId, err := m.k.router.GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG)

	if err != nil {
		return nil, err
	}

	newIsm := types.MerkleRootMultisigISM{
		Id:         ismId,
		Owner:      req.Creator,
		Validators: req.Validators,
		Threshold:  req.Threshold,
	}

	if err = newIsm.Validate(); err != nil {
		return nil, err
	}

	if err = m.k.isms.Set(ctx, ismId, &newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateMerkleRootMultisigIsmResponse{Id: fmt.Sprintf("%v", ismId)}, nil
}

func (m msgServer) CreateNoopIsm(ctx context.Context, ism *types.MsgCreateNoopIsm) (*types.MsgCreateNoopIsmResponse, error) {
	ismId, err := m.k.router.GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	newIsm := types.NoopISM{
		Id:    ismId,
		Owner: ism.Creator,
	}

	if err = m.k.isms.Set(ctx, ismId, &newIsm); err != nil {
		return nil, err
	}

	return &types.MsgCreateNoopIsmResponse{Id: fmt.Sprintf("%v", ismId)}, nil
}
