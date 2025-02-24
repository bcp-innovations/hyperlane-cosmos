package types

import (
	"encoding/binary"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
)

// TODO: Remove
func GetAnnouncementDigest(storageLocation string, domainId uint32, mailbox []byte) [32]byte {
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

	return crypto.Keccak256Hash(announcementDigestBytes)
}
