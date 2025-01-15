package types

import (
	"encoding/binary"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/ethereum/go-ethereum/crypto"
)

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

type Metadata struct {
	raw []byte
}

// validates and creates a new metadata object
func NewMetadata(metadata []byte) (Metadata, error) {
	if len(metadata) < SIGNATURES_OFFSET {
		return Metadata{}, fmt.Errorf("invalid metadata length: got %v, expected at least %v bytes", len(metadata), SIGNATURES_OFFSET)
	}

	signatures := len(metadata) - SIGNATURES_OFFSET

	if signatures%SIGNATURE_LENGTH != 0 {
		return Metadata{}, fmt.Errorf("invalid signatures length in metadata")
	}
	return Metadata{raw: metadata}, nil
}

func (m *Metadata) SignatureAt(index uint32) ([]byte, error) {
	if index > m.SignatureCount() {
		return []byte{}, fmt.Errorf("signature index out of bounce: got index %v with signature count %v", index, m.SignatureCount())
	}
	start := SIGNATURES_OFFSET + (index * SIGNATURE_LENGTH)
	return m.raw[start : start+SIGNATURE_LENGTH], nil
}

func (m *Metadata) SignatureCount() uint32 {
	signatures := len(m.raw) - SIGNATURES_OFFSET
	return uint32(signatures / SIGNATURE_LENGTH)
}

func (m *Metadata) MessageIndex() uint32 {
	return binary.BigEndian.Uint32((m.raw)[MESSAGE_INDEX_OFFSET:])
}

func (m *Metadata) SignedIndex() uint32 {
	return binary.BigEndian.Uint32((m.raw)[SIGNED_INDEX_OFFSET:]) // index is in big endian format
}

func (m *Metadata) SignedMessageId() [32]byte {
	var messageId [32]byte
	copy(messageId[:], (m.raw)[MESSAGE_ID_OFFSET:MESSAGE_ID_OFFSET+32])
	return messageId
}

func (m *Metadata) MerkleTreeHook() [32]byte {
	var hook [32]byte
	copy(hook[:], (m.raw)[:32])
	return hook
}

func (m *Metadata) Proof() [32][32]byte {
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
