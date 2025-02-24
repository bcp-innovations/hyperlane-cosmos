package types

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

func TreeFromProto(tree *Tree) (*util.MerkleTree, error) {
	if len(tree.Branch) != 32 {
		return nil, fmt.Errorf("invalid branch length: %d", len(tree.Branch))
	}

	var branch [32][32]byte
	for i := range tree.Branch {
		if len(tree.Branch[i]) != 32 {
			return nil, fmt.Errorf("invalid branch element length at index %d: %d", i, len(tree.Branch[i]))
		}
		copy(branch[i][:], tree.Branch[i])
	}

	return &util.MerkleTree{
		Branch: branch,
		Count:  tree.Count,
	}, nil
}

func ProtoFromTree(tree *util.MerkleTree) *Tree {
	branch := make([][]byte, len(tree.Branch))

	for i, v := range tree.Branch {
		branch[i] = v[:]
	}

	return &Tree{
		Branch: branch,
		Count:  tree.Count,
	}
}
