package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/ism/types"
	mailboxTypes "github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func multiSigDigest(metadata *types.Metadata, message *mailboxTypes.HyperlaneMessage) [32]byte {
	messageId := message.Id()
	signedRoot := types.BranchRoot(messageId, metadata.Proof(), metadata.MessageIndex())

	return types.CheckpointDigest(
		message.Origin,
		metadata.MerkleTreeHook(),
		signedRoot,
		metadata.SignedIndex(),
		metadata.SignedMessageId(),
	)
}

func (k Keeper) Verify(ctx context.Context, ismId util.HexAddress, rawMetadata []byte, message mailboxTypes.HyperlaneMessage) (verified bool, err error) {
	// Retrieve ISM
	ism, err := k.Isms.Get(ctx, ismId.String())
	if err != nil {
		return false, err
	}

	switch v := ism.Ism.(type) {
	case *types.Ism_MultiSig:

		metadata, err := types.NewMetadata(rawMetadata)
		if err != nil {
			return false, err
		}

		if metadata.SignedIndex() > metadata.MessageIndex() {
			// TODO: invalid merkle index metadata, have to understand this first!
		}

		digest := multiSigDigest(&metadata, &message)
		multiSigIsm := v.MultiSig

		if multiSigIsm.Threshold == 0 {
			return false, fmt.Errorf("invalid ism. no threshold present")
		}

		// Get MultiSig ISM validator public keys
		validatorPubKeys := make(map[string]bool, len(multiSigIsm.ValidatorPubKeys))
		for _, pubKeyStr := range multiSigIsm.ValidatorPubKeys {
			validatorPubKeys[pubKeyStr] = true
		}

		signatures, validSignatures := metadata.SignatureCount(), uint32(0)
		threshold := uint32(multiSigIsm.Threshold)

		// Early return if we can't possibly meet the threshold
		if signatures < multiSigIsm.Threshold {
			return false, nil
		}

		for i := uint32(0); i < signatures && validSignatures < threshold; i++ {
			signature, err := metadata.SignatureAt(i)
			if err != nil {
				break
			}

			recoveredPubkey, err := util.RecoverEthSignature(digest[:], signature)
			if err != nil {
				continue // Skip invalid signatures
			}

			pubKeyHex := hex.EncodeToString(crypto.FromECDSAPub(recoveredPubkey))
			if validatorPubKeys[pubKeyHex] {
				validSignatures++
			}
		}

		if validSignatures >= threshold {
			return true, nil
		}
		return false, nil
	case *types.Ism_Noop:
		return true, nil
	default:
		return false, fmt.Errorf("ism type not supported: %T", v)
	}
}

func (k Keeper) IsmIdExists(ctx context.Context, ismId string) (bool, error) {
	ism, err := k.Isms.Has(ctx, ismId)
	if err != nil {
		return false, err
	}
	return ism, nil
}
