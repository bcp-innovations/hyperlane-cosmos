package types

// TODO: Add ORIGIN_MAILBOX_OFFSET, MERKLE_ROOT_OFFSET
const (
	SIGNATURES_OFFSET = 0
	SIGNATURE_LENGTH  = 65
)

func SignatureAt(metadata []byte, index uint32) []byte {
	start := SIGNATURES_OFFSET + (index * SIGNATURE_LENGTH)
	end := start + SIGNATURE_LENGTH
	return metadata[start:end]
}

// TODO add structs for MerkleMultiSigISM and MessageIdISM
