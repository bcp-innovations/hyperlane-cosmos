package util

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	TreeDepth = 32
	MaxLeaves = (1 << TreeDepth) - 1
)

// MerkleTree represents an incremental merkle tree.
// Contains current branch and the number of inserted leaves in the tree.
type MerkleTree struct {
	Branch [TreeDepth][32]byte
	Count  uint32
}

// ZeroHashes represents an array of TREE_DEPTH zero hashes
var ZeroHashes = [TreeDepth][32]byte{
	hexToBytes("0000000000000000000000000000000000000000000000000000000000000000"),
	hexToBytes("ad3228b676f7d3cd4284a5443f17f1962b36e491b30a40b2405849e597ba5fb5"),
	hexToBytes("b4c11951957c6f8f642c4af61cd6b24640fec6dc7fc607ee8206a99e92410d30"),
	hexToBytes("21ddb9a356815c3fac1026b6dec5df3124afbadb485c9ba5a3e3398a04b7ba85"),
	hexToBytes("e58769b32a1beaf1ea27375a44095a0d1fb664ce2dd358e7fcbfb78c26a19344"),
	hexToBytes("0eb01ebfc9ed27500cd4dfc979272d1f0913cc9f66540d7e8005811109e1cf2d"),
	hexToBytes("887c22bd8750d34016ac3c66b5ff102dacdd73f6b014e710b51e8022af9a1968"),
	hexToBytes("ffd70157e48063fc33c97a050f7f640233bf646cc98d9524c6b92bcf3ab56f83"),
	hexToBytes("9867cc5f7f196b93bae1e27e6320742445d290f2263827498b54fec539f756af"),
	hexToBytes("cefad4e508c098b9a7e1d8feb19955fb02ba9675585078710969d3440f5054e0"),
	hexToBytes("f9dc3e7fe016e050eff260334f18a5d4fe391d82092319f5964f2e2eb7c1c3a5"),
	hexToBytes("f8b13a49e282f609c317a833fb8d976d11517c571d1221a265d25af778ecf892"),
	hexToBytes("3490c6ceeb450aecdc82e28293031d10c7d73bf85e57bf041a97360aa2c5d99c"),
	hexToBytes("c1df82d9c4b87413eae2ef048f94b4d3554cea73d92b0f7af96e0271c691e2bb"),
	hexToBytes("5c67add7c6caf302256adedf7ab114da0acfe870d449a3a489f781d659e8becc"),
	hexToBytes("da7bce9f4e8618b6bd2f4132ce798cdc7a60e7e1460a7299e3c6342a579626d2"),
	hexToBytes("2733e50f526ec2fa19a22b31e8ed50f23cd1fdf94c9154ed3a7609a2f1ff981f"),
	hexToBytes("e1d3b5c807b281e4683cc6d6315cf95b9ade8641defcb32372f1c126e398ef7a"),
	hexToBytes("5a2dce0a8a7f68bb74560f8f71837c2c2ebbcbf7fffb42ae1896f13f7c7479a0"),
	hexToBytes("b46a28b6f55540f89444f63de0378e3d121be09e06cc9ded1c20e65876d36aa0"),
	hexToBytes("c65e9645644786b620e2dd2ad648ddfcbf4a7e5b1a3a4ecfe7f64667a3f0b7e2"),
	hexToBytes("f4418588ed35a2458cffeb39b93d26f18d2ab13bdce6aee58e7b99359ec2dfd9"),
	hexToBytes("5a9c16dc00d6ef18b7933a6f8dc65ccb55667138776f7dea101070dc8796e377"),
	hexToBytes("4df84f40ae0c8229d0d6069e5c8f39a7c299677a09d367fc7b05e3bc380ee652"),
	hexToBytes("cdc72595f74c7b1043d0e1ffbab734648c838dfb0527d971b602bc216c9619ef"),
	hexToBytes("0abf5ac974a1ed57f4050aa510dd9c74f508277b39d7973bb2dfccc5eeb0618d"),
	hexToBytes("b8cd74046ff337f0a7bf2c8e03e10f642c1886798d71806ab1e888d9e5ee87d0"),
	hexToBytes("838c5655cb21c6cb83313b5a631175dff4963772cce9108188b34ac87c81c41e"),
	hexToBytes("662ee4dd2dd7b2bc707961b1e646c4047669dcb6584f0d8d770daf5d7e7deb2e"),
	hexToBytes("388ab20e2573d171a88108e79d820e98f26c0b84aa8b2f4aa4968dbb818ea322"),
	hexToBytes("93237c50ba75ee485f4c22adf2f741400bdf8d6a9cc7df7ecae576221665d735"),
	hexToBytes("8448818bb4ae4562849e949e17ac16e0be16688e156b5cf15e098c627c0056a9"),
}

// hexToBytes converts a hex string to a byte array.
func hexToBytes(hexStr string) [32]byte {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(fmt.Sprintf("failed to decode hex string: %v", err))
	}
	var result [32]byte
	copy(result[:], bytes)
	return result
}

func NewTree(branch [32][32]byte, count uint32) *MerkleTree {
	tree := &MerkleTree{}
	copy(tree.Branch[:], branch[:])
	tree.Count = count
	return tree
}

func (tree *MerkleTree) GetCount() uint32 {
	return tree.Count
}

func (tree *MerkleTree) GetRoot() [32]byte {
	return tree.root()
}

func (tree *MerkleTree) GetLatestCheckpoint() ([32]byte, uint32, error) {
	if tree.Count == 0 {
		return [32]byte{}, 0, fmt.Errorf("no leaf inserted yet")
	}
	return tree.root(), tree.Count - 1, nil
}

// Insert inserts node into merkle tree
// Reverts if tree is full
func (tree *MerkleTree) Insert(node [32]byte) error {
	if tree.Count >= MaxLeaves {
		return fmt.Errorf("merkle tree full")
	}

	tree.Count++
	size := tree.Count
	for i := uint64(0); i < TreeDepth; i++ {
		if (size & 1) == 1 {
			tree.Branch[i] = node
			return nil
		}
		node = [32]byte(crypto.Keccak256(tree.Branch[i][:], node[:]))
		size /= 2
	}
	// As the loop should always end prematurely with the return statement,
	// this code should be unreachable. We panic just to be safe.
	panic("unreachable code")
}

// rootWithCtx calculates and returns tree's current root given array of zero hashes
func (tree *MerkleTree) rootWithCtx(zeroes [TreeDepth][32]byte) [32]byte {
	current := [32]byte{} // zero initialized 32 byte long hash
	index := tree.Count

	for i := uint64(0); i < TreeDepth; i++ {
		ithBit := (index >> i) & 0x01
		next := tree.Branch[i]
		if ithBit == 1 {
			current = crypto.Keccak256Hash(next[:], current[:])
		} else {
			current = crypto.Keccak256Hash(current[:], zeroes[i][:])
		}
	}
	return current
}

// root calculates and returns tree's current root
func (tree *MerkleTree) root() [32]byte {
	return tree.rootWithCtx(ZeroHashes)
}

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
