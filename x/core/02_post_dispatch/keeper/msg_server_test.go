package keeper_test

import (
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_test.go

* Create (invalid) Merkle Tree Hook (invalid mailbox)
* Create (invalid) Merkle Tree Hook (invalid mailbox)
* Create (invalid) Merkle Tree Hook (non-existing mailbox)

*/

var _ = Describe("msg_server_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("Create (invalid) Merkle Tree Hook (invalid mailbox)", func() {
		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleTreeHook{
			Owner:     creator.Address,
			MailboxId: "0x1234",
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id 0x1234: mailbox does not exist"))
	})

	It("Create (invalid) Merkle Tree Hook (non-existing mailbox)", func() {
		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleTreeHook{
			Owner:     creator.Address,
			MailboxId: "0x68797065726c616e650000000000000000000000000000000000000000000000",
		})

		// Assert
		Expect(err.Error()).To(Equal("0x68797065726c616e650000000000000000000000000000000000000000000000: mailbox does not exist"))
	})

	It("Create (valid) Merkle Tree Hook", func() {
		// Arrange
		mailboxId, err := createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())
		println(mailboxId.String())

		// Act
		_, err = s.RunTx(&types.MsgCreateMerkleTreeHook{
			Owner:     creator.Address,
			MailboxId: mailboxId.String(),
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("Create (valid) Noop Hook", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateNoopHook{
			Owner: creator.Address,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateNoopHookResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		// Assert
		Expect(err).To(BeNil())

		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)

		hooks, err := qs.NoopHooks(s.Ctx(), &types.QueryNoopHooksRequest{})
		Expect(err).To(BeNil())
		Expect(hooks.NoopHooks).To(HaveLen(1))
		Expect(hooks.NoopHooks[0].Owner).To(Equal(creator.Address))
		Expect(hooks.NoopHooks[0].Id.String()).To(Equal(response.Id))
	})
})
