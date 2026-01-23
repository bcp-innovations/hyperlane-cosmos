package types

import (
	"context"
	"encoding/binary"
	"fmt"
	"slices"

	storetypes "cosmossdk.io/store/types"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ HyperlaneInterchainSecurityModule = &MessageIdMultisigISM{}

func (m *MessageIdMultisigISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

func (m *MessageIdMultisigISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG
}

func (m *MessageIdMultisigISM) Verify(ctx context.Context, rawMetadata []byte, message util.HyperlaneMessage) (bool, error) {
	metadata, err := NewMessageIdMultisigMetadata(rawMetadata)
	if err != nil {
		return false, err
	}

	digest := metadata.Digest(&message)

	// Charge Cosmos gas for each signature verification.
	sdk.UnwrapSDKContext(ctx).GasMeter().ConsumeGas(storetypes.Gas(1000*len(metadata.Signatures)), "ism signature verification")

	return VerifyMultisig(m.Validators, m.Threshold, metadata.Signatures, digest)
}

func (m *MessageIdMultisigISM) GetThreshold() uint32 {
	return m.Threshold
}

func (m *MessageIdMultisigISM) GetValidators() []string {
	return m.Validators
}

func (m *MessageIdMultisigISM) Validate() error {
	return ValidateNewMultisig(m)
}

type MessageIdMultisigMetadata struct {
	MerkleTreeHook [32]byte
	MerkleRoot     [32]byte
	MerkleIndex    uint32
	SignatureCount uint32
	Signatures     [][]byte
}

// NewMessageIdMultisigMetadata validates and creates a new metadata object
func NewMessageIdMultisigMetadata(metadata []byte) (MessageIdMultisigMetadata, error) {
	/*
	 * Format of metadata:
	 * [   0:  32] Origin merkle tree address
	 * [  32:  64] Signed checkpoint root
	 * [  64:  68] Signed checkpoint index
	 * [  68:????] Validator signatures (length := threshold * 65)
	 */
	// originMerkleTreeOffset := 0
	merkleRootOffset := 32
	merkleIndexOffset := 64
	signaturesOffset := 68
	signatureLength := 65

	if len(metadata) < signaturesOffset {
		return MessageIdMultisigMetadata{}, fmt.Errorf("invalid metadata length: got %v, expected at least %v bytes", len(metadata), signaturesOffset)
	}

	signaturesLen := len(metadata) - signaturesOffset
	signatureCount := uint32(signaturesLen / signatureLength)

	if signaturesLen%signatureLength != 0 {
		return MessageIdMultisigMetadata{}, fmt.Errorf("invalid signatures length in metadata")
	}

	var signatures [][]byte
	for i := 0; i < int(signatureCount); i++ {
		start := signaturesOffset + (i * signatureLength)
		sig := make([]byte, signatureLength)
		copy(sig, metadata[start:start+signatureLength])
		signatures = append(signatures, sig)
	}

	var merkleTreeHook [32]byte
	copy(merkleTreeHook[:], (metadata)[:32])

	var merkleRoot [32]byte
	copy(merkleRoot[:], (metadata)[merkleRootOffset:merkleRootOffset+32])

	return MessageIdMultisigMetadata{
		MerkleTreeHook: merkleTreeHook,
		MerkleRoot:     merkleRoot,
		MerkleIndex:    binary.BigEndian.Uint32((metadata)[merkleIndexOffset:]),
		SignatureCount: uint32(signaturesLen / signaturesOffset),
		Signatures:     signatures,
	}, nil
}

func (m *MessageIdMultisigMetadata) Bytes() []byte {
	merkleIndex := make([]byte, 4)
	binary.BigEndian.PutUint32(merkleIndex, m.MerkleIndex)

	var signaturesBytes []byte
	for _, sig := range m.Signatures {
		signaturesBytes = append(signaturesBytes, sig...)
	}

	return slices.Concat(
		m.MerkleTreeHook[:],
		m.MerkleRoot[:],
		merkleIndex,
		signaturesBytes,
	)
}

func (m *MessageIdMultisigMetadata) Digest(message *util.HyperlaneMessage) [32]byte {
	return checkpointDigest(
		message.Origin,
		m.MerkleTreeHook,
		m.MerkleRoot,
		m.MerkleIndex,
		message.Id(),
	)
}
