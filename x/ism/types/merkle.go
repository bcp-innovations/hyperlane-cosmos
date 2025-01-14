package types

// TODO: make this a unified package with the mailbox merkle rooot

import (
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	TreeDepth = 32
	MaxLeaves = (1 << TreeDepth) - 1
)

// BranchRoot calculates and returns the merkle root for the given leaf item, a merkle branch, and the index of item in the tree.
func BranchRoot(item [32]byte, branch [TreeDepth][32]byte, index uint32) [32]byte {
	current := item

	for i := uint64(0); i < TreeDepth; i++ {
		ithBit := (index >> i) & 0x01
		next := branch[i]
		if ithBit == 1 {
			current = crypto.Keccak256Hash(next[:], current[:])
		} else {
			current = crypto.Keccak256Hash(current[:], next[:])
		}
	}
	return current
}
