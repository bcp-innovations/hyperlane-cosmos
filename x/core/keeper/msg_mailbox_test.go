package keeper_test

import (
	"fmt"

	pdTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	ismtypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

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
		defaultIsm := util.NewZeroAddress()

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:      creator.Address,
			DefaultIsm: defaultIsm,
		})
		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm.String())))
		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with invalid default ISM and without IGP", func() {
		// Arrange
		defaultIsm, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:      creator.Address,
			DefaultIsm: defaultIsm,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with non-existing default ISM and without IGP", func() {
		// Arrange
		defaultIsm, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:      creator.Address,
			DefaultIsm: defaultIsm,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	// invalid IGP
	It("CreateMailbox (invalid) with valid default ISM (Noop) and invalid IGP", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Multisig) and invalid IGP", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Noop) and non-existent IGP", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Multisig) and non-existent IGP", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	// Mailbox valid cases
	It("CreateMailbox (valid) with NoopISM and required IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		noopHookId := createNoopHook(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:        creator.Address,
			DefaultIsm:   ismId,
			RequiredHook: &igpId,
			DefaultHook:  &noopHookId,
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewSingleMailbox(s, res, creator.Address, ismId.String(), igpId.String(), noopHookId.String())
	})

	It("CreateMailbox (valid) with MultisigISM and required IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		noopHookId := createNoopHook(s, creator.Address)
		ismId := createMultisigIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:        creator.Address,
			DefaultIsm:   ismId,
			RequiredHook: &igpId,
			DefaultHook:  &noopHookId,
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewSingleMailbox(s, res, creator.Address, ismId.String(), igpId.String(), noopHookId.String())
	})

	It("CreateMailbox (valid) with NoopISM and optional IGP", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		noopId := createNoopHook(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:        creator.Address,
			DefaultIsm:   ismId,
			RequiredHook: &igpId,
			DefaultHook:  &noopId,
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewSingleMailbox(s, res, creator.Address, ismId.String(), igpId.String(), noopId.String())
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   igpId.String(),
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1250000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
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
			CustomIgp:   "",
			GasLimit:    math.NewInt(50000),
			MaxFee:      sdk.NewCoin("acoin", math.NewInt(1000000)),
		})

		// Assert
		Expect(err).To(BeNil())

		verifyDispatch(s, mailboxId, 1)
	})

	It("ProcessMessage (invalid) with empty message", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
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
			MailboxId: mailboxId,
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
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "xxx",
			Message:   "0xe81bf6f262305f49f318d68f33b04866f092ffdb2ecf9c98469b4a8b829f65e4",
		})

		// Assert
		Expect(err.Error()).To(Equal("failed to decode metadata"))
	})

	PIt("ProcessMessage (valid) (Noop ISM)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		recipient := util.CreateMockHexAddress("recipient", 0)

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       1,
			Origin:      0,
			Sender:      util.CreateMockHexAddress("sender", 0),
			Destination: 1,
			Recipient:   recipient,
			Body:        nil,
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
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
	res, err := s.RunTx(&pdTypes.MsgCreateIgp{
		Owner: creator,
		Denom: "acoin",
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	igpId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return igpId
}

func createNoopHook(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&pdTypes.MsgCreateNoopHook{
		Owner: creator,
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateNoopHookResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	noopHookId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return noopHookId
}

// TODO fix
func createValidMailbox(s *i.KeeperTestSuite, creator string, ism string, igpIsRequiredHook bool, destinationDomain uint32) (util.HexAddress, util.HexAddress, util.HexAddress) {
	var ismId util.HexAddress
	switch ism {
	case "noop":
		ismId = createNoopIsm(s, creator)
	case "multisig":
		ismId = createMultisigIsm(s, creator)
	}

	igpId := createIgp(s, creator)
	noopId := createNoopHook(s, creator)

	err := setDestinationGasConfig(s, creator, igpId.String(), destinationDomain)
	Expect(err).To(BeNil())

	res, err := s.RunTx(&types.MsgCreateMailbox{
		Owner:        creator,
		DefaultIsm:   ismId,
		DefaultHook:  &noopId,
		RequiredHook: &igpId,
	})
	Expect(err).To(BeNil())

	return verifyNewSingleMailbox(s, res, creator, ismId.String(), igpId.String(), noopId.String()), igpId, ismId
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

	return response.Id
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismtypes.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response ismtypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId string, domain uint32) error {
	_, err := s.RunTx(&pdTypes.MsgSetDestinationGasConfig{
		Owner: creator,
		IgpId: igpId,
		DestinationGasConfig: &pdTypes.DestinationGasConfig{
			RemoteDomain: 1,
			GasOracle: &pdTypes.GasOracle{
				TokenExchangeRate: math.NewInt(1e10),
				GasPrice:          math.NewInt(1),
			},
			GasOverhead: math.NewInt(200000),
		},
	})

	return err
}

func verifyNewSingleMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, ismId, requiredHookId, defaultHookId string) util.HexAddress {
	var response types.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
	Expect(err).To(BeNil())
	Expect(mailbox.Owner).To(Equal(creator))
	Expect(mailbox.DefaultIsm.String()).To(Equal(ismId))
	if defaultHookId != "" {
		Expect(mailbox.DefaultHook.String()).To(Equal(defaultHookId))
	} else {
		Expect(mailbox.DefaultHook).To(BeNil())
	}
	if requiredHookId != "" {
		Expect(mailbox.RequiredHook.String()).To(Equal(requiredHookId))
	} else {
		Expect(mailbox.RequiredHook).To(BeNil())
	}
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Owner).To(Equal(creator))

	return mailboxId
}

func verifyInvalidMailboxCreation(s *i.KeeperTestSuite) {
	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(0))
}

func verifyDispatch(s *i.KeeperTestSuite, mailboxId util.HexAddress, messageSent uint32) {
	mailbox, _ := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
	Expect(mailbox.MessageSent).To(Equal(messageSent))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))
}
