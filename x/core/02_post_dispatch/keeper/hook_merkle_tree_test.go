package keeper_test

import (
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_igp.go

* MerkleTreeHook PostDispatch Example
* MerkleTreeHook QuoteDispatch
* MerkleTreeHook HookType
* MerkleTreeHook (valid) exists
* MerkleTreeHook (invalid) exists

*/

var _ = Describe("msg_igp.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	var mailboxId util.HexAddress
	var hookId util.HexAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		mailboxId, err = createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())

		hookId, err = createDummyMerkleTreeHook(s, creator.Address, mailboxId)
		Expect(err).To(BeNil())
	})

	It("MerkleTreeHook PostDispatch Example", func() {
		// Arrange

		// Implementation of the test found here:
		// https://github.com/hyperlane-xyz/hyperlane-monorepo/blob/e90ae5a0d372ab3423fe025f01eea368a6a8c120/solidity/test/MerkleTreeHook.t.sol

		recipient, err := util.DecodeHexAddress("0x00000000000000000000000000000000000000000000000000000000deadbeef")
		Expect(err).To(BeNil())

		// sender: 0x7fa9385be102ac3eac297483dd6233d62b3e1496
		sender, err := util.DecodeHexAddress("0x0000000000000000000000007fa9385be102ac3eac297483dd6233d62b3e1496")
		Expect(err).To(BeNil())

		expectedRoots := []string{
			"0x10df2f89cb24ed6078fc3949b4870e94a7e32e40e8d8c6b7bd74ccc2c933d760",
			"0x080ef1c2cd394de78363ecb0a466c934b57de4abb5604a0684e571990eb7b073",
			"0xbf78ad252da524f1e733aa6b83514dd83225676b5828f888f01487108f8f7cc7",
		}

		// test for 3 messages
		for k := 0; k < 3; k++ {

			// adjust message body
			body := make([]byte, 32)
			body[31] = byte(k)

			message := util.HyperlaneMessage{
				Version:     3,
				Nonce:       0, // nonce is not updated in this test-case
				Origin:      11,
				Sender:      sender,
				Destination: 22,
				Recipient:   recipient,
				Body:        body,
			}

			// Act
			fee, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, hookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
			Expect(err).To(BeNil())
			qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)
			hook, err := qs.MerkleTreeHook(s.Ctx(), &types.QueryMerkleTreeHook{Id: hookId.String()})
			Expect(err).To(BeNil())

			Expect(hook.MerkleTreeHook.MerkleTree.Count).To(Equal(uint32(k + 1)))
			Expect(util.HexAddress(hook.MerkleTreeHook.MerkleTree.Root).String()).To(Equal(expectedRoots[k]))

			// Assert
			Expect(err).To(BeNil())
			Expect(fee).To(Equal(sdk.NewCoins()))
		}
	})

	It("MerkleTreeHook QuoteDispatch", func() {
		// Merkle Tree Hook never charges coins for execution

		// Arrange
		message := util.HyperlaneMessage{}

		// Act
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(hookId)
		Expect(err).To(BeNil())
		fee, err := (*handler).QuoteDispatch(s.Ctx(), mailboxId, hookId, util.StandardHookMetadata{}, message)

		// Assert
		Expect(err).To(BeNil())
		Expect(fee).To(Equal(sdk.NewCoins()))
	})

	It("MerkleTreeHook HookType", func() {
		// Arrange

		// Act
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(hookId)
		Expect(err).To(BeNil())
		hookType := (*handler).HookType()

		// Assert
		Expect(hookType).To(Equal(uint8(3)))
	})

	It("MerkleTreeHook (valid) exists", func() {
		// Arrange

		// Act
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(hookId)
		Expect(err).To(BeNil())
		exists, err := (*handler).Exists(s.Ctx(), hookId)

		// Assert
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
	})

	It("MerkleTreeHook (invalid) exists", func() {
		// Arrange
		hookId[31] = byte(10)

		// Act
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(hookId)
		Expect(err).To(BeNil())
		exists, err := (*handler).Exists(s.Ctx(), hookId)

		// Assert
		Expect(err).To(BeNil())
		Expect(exists).To(BeFalse())
	})
})
