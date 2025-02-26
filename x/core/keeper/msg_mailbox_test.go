package keeper_test

import (
	"fmt"

	ismtypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/types"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_mailbox.go

* CreateMailbox (invalid) without default ISM and without IGP
* CreateMailbox (invalid) with invalid default ISM and without IGP
* CreateMailbox (invalid) with non-existing default ISM and without IGP
* CreateMailbox (invalid) with valid default ISM (Noop) and invalid IGP
* CreateMailbox (invalid) with valid default ISM (Multisig) and invalid IGP
* CreateMailbox (invalid) with valid default ISM (Noop) and non-existent IGP
* CreateMailbox (invalid) with valid default ISM (Multisig) and non-existent IGP
* CreateMailbox (valid) with NoopISM and required IGP
* CreateMailbox (valid) with MultisigISM and required IGP
* CreateMailbox (valid) with NoopISM and optional IGP
* CreateMailbox (valid) with MultisigISM and optional IGP
* DispatchMessage (invalid) with empty body
* DispatchMessage (invalid) with invalid body
* DispatchMessage (invalid) with invalid Mailbox ID
* DispatchMessage (invalid) with empty sender
* DispatchMessage (invalid) with invalid sender
* DispatchMessage (invalid) with empty recipient
* DispatchMessage (invalid) with invalid recipient
* DispatchMessage (valid) with optional IGP
* DispatchMessage (valid) with optional and no specified IGP
* DispatchMessage (valid) with required IGP
* ProcessMessage (invalid) with invalid Mailbox ID
* ProcessMessage (invalid) with empty message
* ProcessMessage (invalid) with invalid non-hex message
* ProcessMessage (invalid) with invalid metadata

*/

var _ = Describe("msg_mailbox.go", Ordered, func() {
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

	// CreateMailbox
	// invalid ISM
	It("CreateMailbox (invalid) without default ISM and without IGP", func() {
		// Arrange
		defaultIsm := ""

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: defaultIsm,
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: invalid hex address length", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with invalid default ISM and without IGP", func() {
		// Arrange
		defaultIsm := "0x1234"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: defaultIsm,
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: invalid hex address length", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with non-existing default ISM and without IGP", func() {
		// Arrange
		defaultIsm := "0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: defaultIsm,
			Igp:        nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	// invalid IGP
	It("CreateMailbox (invalid) with valid default ISM (Noop) and invalid IGP", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId := "0x1234"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp id %s is invalid: invalid hex address length", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Multisig) and invalid IGP", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId := "0x1234"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp id %s is invalid: invalid hex address length", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Noop) and non-existent IGP", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Multisig) and non-existent IGP", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId,
				Required: false,
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	// Mailbox valid cases
	It("CreateMailbox (valid) with NoopISM and required IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId.String(),
				Required: true,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId.String(), ismId.String(), true)
	})

	It("CreateMailbox (valid) with MultisigISM and required IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createMultisigIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId.String(),
				Required: true,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId.String(), ismId.String(), true)
	})

	It("CreateMailbox (valid) with NoopISM and optional IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId.String(),
				Required: false,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId.String(), ismId.String(), false)
	})

	It("CreateMailbox (valid) with MultisigISM and optional IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		ismId := createMultisigIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Creator:    creator.Address,
			DefaultIsm: ismId.String(),
			Igp: &types.InterchainGasPaymaster{
				Id:       igpId.String(),
				Required: false,
			},
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewMailbox(s, res, creator.Address, igpId.String(), ismId.String(), false)
	})

	// DispatchMessage
	It("DispatchMessage (invalid) with empty body", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid body: empty hex string"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (invalid) with invalid body", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "12345",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid body: hex string without 0x prefix"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (invalid) with invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String() + "test",
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: invalid hex address length"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (invalid) with empty sender", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      "",
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid sender: empty address string is not allowed"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (invalid) with invalid sender", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      "hyperlane01234567889",
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid sender: decoding bech32 failed: invalid checksum (expected ca6a9q got 567889)"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (invalid) with empty recipient", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid recipient: invalid hex address length"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (invalid) with invalid recipient", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e7gzc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid recipient: encoding/hex: invalid byte: U+0067 'g'"))

		verifyDispatch(s, mailboxId, 0)
	})

	It("DispatchMessage (valid) with optional IGP", func() {
		// Arrange
		mailboxId, igpId, _ := createValidMailbox(s, creator.Address, "noop", false, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       igpId.String(),
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err).To(BeNil())

		verifyDispatch(s, mailboxId, 1)
	})

	It("DispatchMessage (valid) with optional and no specified IGP", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", false, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err).To(BeNil())

		verifyDispatch(s, mailboxId, 1)
	})

	It("DispatchMessage (valid) with required IGP", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgDispatchMessage{
			MailboxId:   mailboxId.String(),
			Sender:      sender.Address,
			Destination: 1,
			Recipient:   "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647",
			Body:        "0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
			IgpId:       "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      math.NewInt(1000000),
		})

		// Assert
		Expect(err).To(BeNil())

		verifyDispatch(s, mailboxId, 1)
	})

	// ProcessMessage() tests (only with Noop ISM)
	It("ProcessMessage (invalid) with invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		invalidMailboxId := mailboxId.String() + "test"

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: invalidMailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   "",
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: invalid hex address length"))
	})

	It("ProcessMessage (invalid) with empty message", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   "",
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid message"))
	})

	It("ProcessMessage (invalid) with invalid non-hex message", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   "test123",
		})

		// Assert
		Expect(err.Error()).To(Equal("failed to decode message"))
	})

	It("ProcessMessage (invalid) with invalid metadata (Noop ISM)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "xxx",
			Message:   "0xe81bf6f262305f49f318d68f33b04866f092ffdb2ecf9c98469b4a8b829f65e4",
		})

		// Assert
		Expect(err.Error()).To(Equal("failed to decode metadata"))
	})

	It("ProcessMessage (valid) (Noop ISM)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		recipient := util.CreateHexAddress("recipient", 0)

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       1,
			Origin:      0,
			Sender:      util.CreateHexAddress("sender", 0),
			Destination: 1,
			Recipient:   recipient,
			Body:        nil,
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err).To(BeNil())
	})

	// TODO: ProcessMessage (valid) (Multisig ISM)
})

// Utils
func createIgp(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&types.MsgCreateIgp{
		Owner: creator,
		Denom: "acoin",
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	igpId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return igpId
}

func createValidMailbox(s *i.KeeperTestSuite, creator string, ism string, igpRequired bool, destinationDomain uint32) (util.HexAddress, util.HexAddress, util.HexAddress) {
	var ismId util.HexAddress
	switch ism {
	case "noop":
		ismId = createNoopIsm(s, creator)
	case "multisig":
		ismId = createMultisigIsm(s, creator)
	}

	igpId := createIgp(s, creator)

	err := setDestinationGasConfig(s, creator, igpId.String(), destinationDomain)
	Expect(err).To(BeNil())

	res, err := s.RunTx(&types.MsgCreateMailbox{
		Creator:    creator,
		DefaultIsm: ismId.String(),
		Igp: &types.InterchainGasPaymaster{
			Id:       igpId.String(),
			Required: igpRequired,
		},
	})
	Expect(err).To(BeNil())

	return verifyNewMailbox(s, res, creator, igpId.String(), ismId.String(), igpRequired), igpId, ismId
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismtypes.MsgCreateMerkleRootMultisigIsm{
		Creator: creator,
		Validators: []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
		},
		Threshold: 2,
	})
	Expect(err).To(BeNil())

	var response ismtypes.MsgCreateMerkleRootMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	ismId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return ismId
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismtypes.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response ismtypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	ismId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return ismId
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId string, domain uint32) error {
	_, err := s.RunTx(&types.MsgSetDestinationGasConfig{
		Owner: creator,
		IgpId: igpId,
		DestinationGasConfig: &types.DestinationGasConfig{
			RemoteDomain: 1,
			GasOracle: &types.GasOracle{
				TokenExchangeRate: math.NewInt(1e10),
				GasPrice:          math.NewInt(1),
			},
			GasOverhead: math.NewInt(200000),
		},
	})

	return err
}

func verifyNewMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, igpId, ismId string, igpRequired bool) util.HexAddress {
	var response types.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.Bytes())
	Expect(err).To(BeNil())
	Expect(mailbox.Creator).To(Equal(creator))
	Expect(mailbox.Igp.Id).To(Equal(igpId))
	Expect(mailbox.DefaultIsm).To(Equal(ismId))
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

	if igpRequired {
		Expect(mailbox.Igp.Required).To(BeTrue())
	} else {
		Expect(mailbox.Igp.Required).To(BeFalse())
	}

	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Creator).To(Equal(creator))

	return mailboxId
}

func verifyInvalidMailboxCreation(s *i.KeeperTestSuite) {
	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(0))
}

func verifyDispatch(s *i.KeeperTestSuite, mailboxId util.HexAddress, messageSent uint32) {
	mailbox, _ := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.Bytes())
	Expect(mailbox.MessageSent).To(Equal(messageSent))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))
	Expect(mailbox.Tree.Count).To(Equal(messageSent))

	if messageSent == 0 {
		_, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).LatestCheckpoint(s.Ctx(), &types.QueryLatestCheckpointRequest{Id: mailboxId.String()})
		Expect(err.Error()).To(Equal("no leaf inserted yet"))
	} else {
		latestCheckpoint, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).LatestCheckpoint(s.Ctx(), &types.QueryLatestCheckpointRequest{Id: mailboxId.String()})
		Expect(err).To(BeNil())

		Expect(latestCheckpoint.Count).To(Equal(messageSent - 1))
	}

	// TODO: Check claimable fees of IGP
}
