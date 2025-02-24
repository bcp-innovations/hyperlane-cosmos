package keeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func multiSigDigest(metadata *types.Metadata, message *util.HyperlaneMessage) [32]byte {
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

func (k Keeper) Verify(ctx context.Context, ismId util.HexAddress, rawMetadata []byte, message util.HyperlaneMessage) (verified bool, err error) {
	// Retrieve ISM
	ism, err := k.Isms.Get(ctx, ismId.Bytes())
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
			return false, fmt.Errorf("invalid signed index")
		}

		digest := multiSigDigest(&metadata, &message)
		multiSigIsm := v.MultiSig

		if multiSigIsm.Threshold == 0 {
			return false, fmt.Errorf("invalid ism. no threshold present")
		}

		signatures, validSignatures := metadata.SignatureCount(), uint32(0)
		threshold := multiSigIsm.Threshold

		// Early return if we can't possibly meet the threshold
		if signatures < multiSigIsm.Threshold {
			return false, nil
		}

		// Get MultiSig ISM validator public keys
		validatorAddresses := make(map[string]bool, len(multiSigIsm.Validators))
		for _, address := range multiSigIsm.Validators {
			validatorAddresses[strings.ToLower(address)] = true
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

			addressBytes := crypto.PubkeyToAddress(*recoveredPubkey)
			address := util.EncodeEthHex(addressBytes[:])
			if validatorAddresses[address] {
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
