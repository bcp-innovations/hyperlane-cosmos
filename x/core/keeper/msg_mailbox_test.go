package keeper_test

import (
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server_create_funder.go

* Create New (invalid) Mailbox without default ISM and without IGP
* Create New (invalid) Mailbox without default ISM and non-existent IGP
* Create New (valid) Mailbox

*/

var _ = Describe("msg_mailbox.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("Create New (invalid) Mailbox without default ISM and without IGP", func() {
		// Arrange
		// nothing to do

		// Act
		// TODO improve test, does not actually test whats claimed in the message
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: "",
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid hex address length"))
	})

	It("Create New (invalid) Mailbox without default ISM and non-existent IGP", func() {
		// Arrange
		// nothing to do

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: "",
			// TODO improve test, does not actually test whats claimed in the message
			Igp: &types.InterchainGasPaymaster{
				Id:       "0x1234",
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid hex address length"))
	})

	It("Create New (valid) Mailbox", func() {
		// Arrange
		_, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: "uhyp",
		})
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgCreateNoopIsm{Creator: creator.Address})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0",
			Igp: &types.InterchainGasPaymaster{
				Id:       "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
				Required: false,
			},
		})

		// Assert
		Expect(err).To(BeNil())
		mailbox, _ := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), util.CreateHexAddress(types.ModuleName, int64(0)).Bytes())
		Expect(mailbox.Creator).To(Equal(creator.Address))
		Expect(mailbox.Igp.Id).To(Equal("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"))
		Expect(mailbox.DefaultIsm).To(Equal("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0"))
		Expect(mailbox.MessageSent).To(Equal(uint32(0)))
		Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

		mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
		Expect(err).To(BeNil())
		Expect(mailboxes.Mailboxes).To(HaveLen(1))
		Expect(mailboxes.Mailboxes[0].Creator).To(Equal(creator.Address))
	})
})
