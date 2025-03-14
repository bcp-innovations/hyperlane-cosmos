package util_test

import (
	"encoding/hex"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - hyperlane_message.go

* Create new empty tree

*/

var _ = Describe("merkle_tree.go", Ordered, func() {
	BeforeEach(func() {
	})

	It("Create new empty tree", func() {
		// Arrange
		tree := util.NewTree(util.ZeroHashes, 0)

		// Act
		checkpoint, n, err := tree.GetLatestCheckpoint()

		// Assert
		Expect(tree).ToNot(BeNil())
		Expect(tree.Count).To(Equal(uint32(0)))
		Expect(tree.GetCount()).To(Equal(uint32(0)))

		root := tree.GetRoot()
		Expect(hex.EncodeToString(root[:])).To(Equal("27ae5ba08d7291c96c8cbddcc148bf48a6d68c7974b94356f53754ef6171d757"))
		Expect(n).To(Equal(uint32(0)))
		Expect(err.Error()).To(Equal("no leaf inserted yet"))
		Expect(checkpoint).To(Equal([32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}))
	})
})
