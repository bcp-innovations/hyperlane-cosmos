package types

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var _ HyperlaneInterchainSecurityModule = &MerkleRootMultiSigISM{}

func (m *MerkleRootMultiSigISM) GetId() uint64 {
	return m.Id
}

func (m *MerkleRootMultiSigISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG
}

func (m *MerkleRootMultiSigISM) Verify(_ sdk.Context, rawMetadata []byte, message util.HyperlaneMessage) (bool, error) {
	metadata, err := NewMetadata(rawMetadata)
	if err != nil {
		return false, err
	}

	if metadata.SignedIndex() > metadata.MessageIndex() {
		return false, fmt.Errorf("invalid signed index")
	}

	digest := multiSigDigest(&metadata, &message)

	if m.Threshold == 0 {
		return false, fmt.Errorf("invalid ism. no threshold present")
	}

	signatures, validSignatures := metadata.SignatureCount(), uint32(0)
	threshold := m.Threshold

	// Early return if we can't possibly meet the threshold
	if signatures < m.Threshold {
		return false, nil
	}

	// Get MultiSig ISM validator public keys
	validatorAddresses := make(map[string]bool, len(m.Validators))
	for _, address := range m.Validators {
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
}

func (m *MerkleRootMultiSigISM) Validate() error {
	if m.Threshold == 0 {
		return fmt.Errorf("threshold must be greater than zero")
	}

	if len(m.Validators) < int(m.Threshold) {
		return fmt.Errorf("validator addresses less than threshold")
	}

	for _, validatorAddress := range m.Validators {
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
	for _, address := range m.Validators {
		count[address]++
		if count[address] > 1 {
			return fmt.Errorf("duplicate validator address: %v", address)
		}
	}
	return nil
}

/**
 * Format of m:
 * [   0:  32] Origin merkle tree address
 * [  32:  36] Index of message ID in merkle tree
 * [  36:  68] Signed checkpoint message ID
 * [  68:1092] Merkle proof
 * [1092:1096] Signed checkpoint index (computed from proof and index)
 * [1096:????] Validator signatures (length := threshold * 65)
 */
const (
	ORIGIN_MERKLE_TREE_OFFSET = 0
	MESSAGE_INDEX_OFFSET      = 32
	MESSAGE_ID_OFFSET         = 36
	MERKLE_PROOF_OFFSET       = 68
	MERKLE_PROOF_LENGTH       = 32 * 32
	SIGNED_INDEX_OFFSET       = 1092
	SIGNATURES_OFFSET         = 1096
	SIGNATURE_LENGTH          = 65
)

type MerkleRootMultiSigMetadata struct {
	raw []byte
}

// NewMetadata validates and creates a new metadata object
func NewMetadata(metadata []byte) (MerkleRootMultiSigMetadata, error) {
	if len(metadata) < SIGNATURES_OFFSET {
		return MerkleRootMultiSigMetadata{}, fmt.Errorf("invalid metadata length: got %v, expected at least %v bytes", len(metadata), SIGNATURES_OFFSET)
	}

	signatures := len(metadata) - SIGNATURES_OFFSET

	if signatures%SIGNATURE_LENGTH != 0 {
		return MerkleRootMultiSigMetadata{}, fmt.Errorf("invalid signatures length in metadata")
	}
	return MerkleRootMultiSigMetadata{raw: metadata}, nil
}

func (m *MerkleRootMultiSigMetadata) SignatureAt(index uint32) ([]byte, error) {
	if index > m.SignatureCount() {
		return []byte{}, fmt.Errorf("signature index out of bounce: got index %v with signature count %v", index, m.SignatureCount())
	}
	start := SIGNATURES_OFFSET + (index * SIGNATURE_LENGTH)
	return m.raw[start : start+SIGNATURE_LENGTH], nil
}

func (m *MerkleRootMultiSigMetadata) SignatureCount() uint32 {
	signatures := len(m.raw) - SIGNATURES_OFFSET
	return uint32(signatures / SIGNATURE_LENGTH)
}

func (m *MerkleRootMultiSigMetadata) MessageIndex() uint32 {
	return binary.BigEndian.Uint32((m.raw)[MESSAGE_INDEX_OFFSET:])
}

func (m *MerkleRootMultiSigMetadata) SignedIndex() uint32 {
	return binary.BigEndian.Uint32((m.raw)[SIGNED_INDEX_OFFSET:]) // index is in big endian format
}

func (m *MerkleRootMultiSigMetadata) SignedMessageId() [32]byte {
	var messageId [32]byte
	copy(messageId[:], (m.raw)[MESSAGE_ID_OFFSET:MESSAGE_ID_OFFSET+32])
	return messageId
}

func (m *MerkleRootMultiSigMetadata) MerkleTreeHook() [32]byte {
	var hook [32]byte
	copy(hook[:], (m.raw)[:32])
	return hook
}

func (m *MerkleRootMultiSigMetadata) Proof() [32][32]byte {
	proof := (m.raw)[MERKLE_PROOF_OFFSET : MERKLE_PROOF_OFFSET+MERKLE_PROOF_LENGTH]
	// proof is a 32 element long array of hashes encoded as 32 byte long arrays
	var decodedProof [32][32]byte
	for i := 0; i < 32; i++ {
		copy(decodedProof[i][:], proof[i*32:(i+1)*32])
	}
	return decodedProof
}

func CheckpointDigest(origin uint32, merkleTreeHook, checkpointRoot [32]byte, checkpointIndex uint32, messageId [32]byte) [32]byte {
	domainHash := DomainHash(origin, merkleTreeHook)

	bytes := make([]byte, 0, 32+32+4+32)
	bytes = append(bytes, domainHash[:]...)
	bytes = append(bytes, checkpointRoot[:]...)
	bytes = binary.BigEndian.AppendUint32(bytes, checkpointIndex)
	bytes = append(bytes, messageId[:]...)

	return util.GetEthSigningHash(crypto.Keccak256(bytes))
}

func DomainHash(origin uint32, merkleTreeHook [32]byte) [32]byte {
	bytes := make([]byte, 0, 46)

	bytes = binary.BigEndian.AppendUint32(bytes, origin)
	bytes = append(bytes, merkleTreeHook[:]...)
	bytes = append(bytes, []byte("HYPERLANE")...)

	return crypto.Keccak256Hash(bytes)
}

func multiSigDigest(metadata *MerkleRootMultiSigMetadata, message *util.HyperlaneMessage) [32]byte {
	messageId := message.Id()
	signedRoot := util.BranchRoot(messageId, metadata.Proof(), metadata.MessageIndex())

	return CheckpointDigest(
		message.Origin,
		metadata.MerkleTreeHook(),
		signedRoot,
		metadata.SignedIndex(),
		metadata.SignedMessageId(),
	)
}
