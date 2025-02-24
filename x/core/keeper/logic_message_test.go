package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_message.go

* DispatchMessage (invalid) with non-existing Mailbox ID
* ProcessMessage (invalid) with non-existing Mailbox ID
* ProcessMessage (invalid) with invalid hex message
* ProcessMessage (invalid) already processed message (replay protection)
* ProcessMessage (invalid) with invalid message: non-registered recipient

*/

var _ = Describe("logic_message.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var sender i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		sender = i.GenerateTestValidatorAddress("Sender")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("DispatchMessage (invalid) with non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   nonExistingMailboxId,
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))

		verifyDispatch(s, mailboxId, 0)
	})

	It("ProcessMessage (invalid) with non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"
		createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateHexAddress("test", 0)
		recipientHex := util.CreateHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: nonExistingMailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("ProcessMessage (invalid) with invalid hex message", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateHexAddress("test", 0)
		recipientHex := util.CreateHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String()[:util.BodyOffset-1],
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid hyperlane message"))
	})

	It("ProcessMessage (invalid) already processed message (replay protection)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateHexAddress("test", 0)
		recipientHex := util.CreateHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte(""),
		}

		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("already received messsage with id %s", hypMsg.Id())))
	})

	// TODO rework test ->
	PIt("ProcessMessage (invalid) with invalid message: non-registered recipient", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateHexAddress("test", 0)
		recipientHex := util.CreateHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to get receiver ism address for recipient: %s", recipientHex)))
	})
})
