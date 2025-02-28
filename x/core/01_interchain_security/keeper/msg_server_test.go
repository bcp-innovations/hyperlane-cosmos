package keeper_test

import (
	"crypto/ecdsa"
	"fmt"

	types2 "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	keeper2 "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server.go

* Create (valid) Noop ISM
* Create (invalid) MessageIdMultisig ISM with less addresses
* Create (invalid) MessageIdMultisig ISM with invalid threshold
* Create (invalid) MessageIdMultisig ISM with duplicate validator addresses
* Create (invalid) MessageIdMultisig ISM with invalid validator addresses
* Create (valid) MessageIdMultisig ISM
* Create (invalid) MerkleRootMultisig ISM with less addresses
* Create (invalid) MerkleRootMultisig ISM with invalid threshold
* Create (invalid) MerkleRootMultisig ISM with duplicate validator addresses
* Create (invalid) MerkleRootMultisig ISM with invalid validator addresses
* Create (valid) MerkleRootMultisig ISM
* AnnounceValidator (invalid) with empty validator
* AnnounceValidator (invalid) with invalid validator
* AnnounceValidator (invalid) with empty storage location
* AnnounceValidator (invalid) with empty signature
* AnnounceValidator (invalid) with invalid signature
* AnnounceValidator (invalid) with invalid signature recovery id
* AnnounceValidator (invalid) same storage location for validator (replay protection)
* AnnounceValidator (invalid) for non-existing Mailbox ID
* AnnounceValidator (invalid) for invalid Mailbox ID
* AnnounceValidator (invalid) for non-matching signature validator pair
* AnnounceValidator (valid)
* AnnounceValidator (valid) add another storage location

*/

var _ = Describe("msg_server.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("Create (valid) Noop ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateNoopIsm{
			Creator: creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateNoopIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		queryServer := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper)
		ism, err := queryServer.Ism(s.Ctx(), &types.QueryIsmRequest{Id: response.Id.String()})
		Expect(err).To(BeNil())
		Expect(ism.Ism.TypeUrl).To(Equal("/hyperlane.core.interchain_security.v1.NoopISM"))
		var noopIsm types.NoopISM
		err = proto.Unmarshal(ism.Ism.Value, &noopIsm)
		Expect(err).To(BeNil())
		Expect(noopIsm.Owner).To(Equal(creator.Address))
		Expect(noopIsm.Id.String()).To(Equal(response.Id.String()))
	})

	It("Create (invalid) MessageIdMultisig ISM with less address", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses less than threshold: invalid multisig configuration"))
	})

	It("Create (invalid) MessageIdMultisig ISM with invalid threshold", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  0,
		})

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero: invalid multisig configuration"))
	})

	It("Create (invalid) MessageIdMultisig ISM with duplicate validator addresses", func() {
		// Arrange
		invalidAddress := []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
		}

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Validators: invalidAddress,
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicate validator address: %v: invalid multisig configuration", invalidAddress[0])))
	})

	It("Create (invalid) MessageIdMultisig ISM with invalid validator address", func() {
		// Arrange
		validValidatorAddress := "0xb05b6a0aa112b61a7aa16c19cac27d970692995e"
		invalidAddress := []string{
			// one character less
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995",
			// one character more
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995ef",
			// invalid character included (`t`)
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995t",
		}

		for _, invalidKey := range invalidAddress {
			// Act
			_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
				Creator: creator.Address,
				Validators: []string{
					validValidatorAddress,
					invalidKey,
				},
				Threshold: 2,
			})

			// Assert
			Expect(err.Error()).To(Equal(fmt.Sprintf("invalid validator address: %s: invalid multisig configuration", invalidKey)))
		}
	})

	It("Create (valid) MessageIdMultisig ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateMessageIdMultisigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		// ismId, err := util.DecodeHexAddress(response.Id)
		// Expect(err).To(BeNil())

		// ism, err := keeper.NewQueryServerImpl(&s.App().IsmKeeper).Ism(s.Ctx(), &types.QueryIsmRequest{Id: ismId.String()})
		// Expect(ism.Ism).To(BeAssignableToTypeOf(&types.MerkleRootMultisigISM{}))
		// Expect(ism).To(Equal(creator.Address))
		// Expect(ism.IsmType).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TPYE_MESSAGE_ID_MULTISIG))
	})

	It("Create (invalid) MerkleRootMultisig ISM with less address", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses less than threshold: invalid multisig configuration"))
	})

	It("Create (invalid) MerkleRootMultisig ISM with invalid threshold", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  0,
		})

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero: invalid multisig configuration"))
	})

	It("Create (invalid) MerkleRootMultisig ISM with duplicate validator addresses", func() {
		// Arrange
		invalidAddress := []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
		}

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Validators: invalidAddress,
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicate validator address: %v: invalid multisig configuration", invalidAddress[0])))
	})

	It("Create (invalid) MerkleRootMultisig ISM with invalid validator address", func() {
		// Arrange
		validValidatorAddress := "0xb05b6a0aa112b61a7aa16c19cac27d970692995e"
		invalidAddress := []string{
			// one character less
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995",
			// one character more
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995ef",
			// invalid character included (`t`)
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995t",
		}

		for _, invalidKey := range invalidAddress {
			// Act
			_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
				Creator: creator.Address,
				Validators: []string{
					validValidatorAddress,
					invalidKey,
				},
				Threshold: 2,
			})

			// Assert
			Expect(err.Error()).To(Equal(fmt.Sprintf("invalid validator address: %s: invalid multisig configuration", invalidKey)))
		}
	})

	It("Create (valid) MerkleRootMultisig ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateMerkleRootMultisigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		// ismId, err := util.DecodeHexAddress(response.Id)
		// Expect(err).To(BeNil())

		// ism, err := keeper.NewQueryServerImpl(&s.App().IsmKeeper).Ism(s.Ctx(), &types.QueryIsmRequest{Id: ismId.String()})
		// Expect(ism.Ism).To(BeAssignableToTypeOf(&types.MerkleRootMultisigISM{}))
		// Expect(ism).To(Equal(creator.Address))
		// Expect(ism.IsmType).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG))
	})

	It("AnnounceValidator (invalid) with empty validator", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       "",
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator cannot be empty"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with invalid validator", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		invalidValidatorAddress := "0x0b1caf89d1edb9ee161093test94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       invalidValidatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid validator address"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with empty storage location", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: "",
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("storage location cannot be empty"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with empty signature", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		storageLocation := "aws://key.pub"

		// Act
		_, err := s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       "",
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("signature cannot be empty"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with invalid signature", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		storageLocation := "aws://key.pub"

		// Act
		_, err := s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       "0x0b1caf89d1edb9ee161093b1ec94ca75611dtest",
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid signature"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with invalid signature recovery id", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, true)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid signature recovery id"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) same storage location for validator (replay protection)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("validator %s already announced storage location %s", validatorAddress, storageLocation)))
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{storageLocation})
	})

	It("AnnounceValidator (invalid) for non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, err := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		Expect(err).To(BeNil())

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, nonExistingMailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       nonExistingMailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId.String())))
	})

	It("AnnounceValidator (invalid) for invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String() + "test",
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) for non-matching signature validator pair", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		correctValidatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		wrongValidatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db392"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       wrongValidatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("validator %s doesn't match signature. recovered address: %s", wrongValidatorAddress, correctValidatorAddress)))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{
			storageLocation,
		})
	})

	It("AnnounceValidator (valid) add another storage location", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})
		Expect(err).To(BeNil())

		storageLocation2 := "aws://key2.pub"
		signature = announce(validatorPrivKey, storageLocation2, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation2,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{
			storageLocation,
			storageLocation2,
		})
	})
})

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

	res, err := s.RunTx(&types2.MsgCreateMailbox{
		Owner:      creator,
		DefaultIsm: ismId.String(),
		// TODO fix
		//Igp: &types.InterchainGasPaymaster{
		//	Id:       igpId,
		//	Required: false,
		//},
	})
	_ = ismId
	_ = igpId
	Expect(err).To(BeNil())

	return verifyNewMailbox(s, res, creator, igpId.String(), ismId.String(), igpRequired), igpId, ismId
}

func createIgp(s *i.KeeperTestSuite, creator string) util.HexAddress {
	//res, err := s.RunTx(&types2.MsgCreateIgp{
	//	Owner: creator,
	//	Denom: "acoin",
	//})
	//Expect(err).To(BeNil())
	//
	//var response types2.MsgCreateIgpResponse
	//err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	//Expect(err).To(BeNil())
	//
	//igpId, err := util.DecodeHexAddress(response.Id)
	//Expect(err).To(BeNil())
	//
	//return igpId
	return util.HexAddress{} // TODO fix
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
		Creator: creator,
		Validators: []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
		},
		Threshold: 2,
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateMerkleRootMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&types.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId string, domain uint32) error {
	//_, err := s.RunTx(&types2.MsgSetDestinationGasConfig{
	//	Owner: creator,
	//	IgpId: igpId,
	//	DestinationGasConfig: &types2.DestinationGasConfig{
	//		RemoteDomain: 1,
	//		GasOracle: &types2.GasOracle{
	//			TokenExchangeRate: math.NewInt(1e10),
	//			GasPrice:          math.NewInt(1),
	//		},
	//		GasOverhead: math.NewInt(200000),
	//	},
	//})
	//
	//return err
	return nil // TODO fix test
}

func verifyNewMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, igpId, ismId string, igpRequired bool) util.HexAddress {
	var response types2.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.Bytes())
	Expect(err).To(BeNil())
	Expect(mailbox.Owner).To(Equal(creator))
	// Expect(mailbox.Igp.Id).To(Equal(igpId)) // TODO
	Expect(mailbox.DefaultIsm).To(Equal(ismId))
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

	//if igpRequired { // TODO
	//	Expect(mailbox.Igp.Required).To(BeTrue())
	//} else {
	//	Expect(mailbox.Igp.Required).To(BeFalse())
	//}

	mailboxes, err := keeper2.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types2.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Owner).To(Equal(creator))

	return mailboxId
}

func announce(privKey, storageLocation string, mailboxId util.HexAddress, localDomain uint32, skipRecoveryId bool) string {
	announcementDigest := types.GetAnnouncementDigest(storageLocation, localDomain, mailboxId.Bytes())

	ethDigest := util.GetEthSigningHash(announcementDigest[:])

	privateKey, err := crypto.HexToECDSA(privKey)
	Expect(err).To(BeNil())

	publicKey := privateKey.Public()
	_, ok := publicKey.(*ecdsa.PublicKey)
	Expect(ok).To(BeTrue())

	signedAnnouncement, err := crypto.Sign(ethDigest[:], privateKey)
	Expect(err).To(BeNil())

	if !skipRecoveryId {
		// Required for recovery ID
		// https://eips.ethereum.org/EIPS/eip-155
		signedAnnouncement[64] += 27
	}

	return util.EncodeEthHex(signedAnnouncement)
}

func validateAnnouncement(s *i.KeeperTestSuite, mailboxId, validatorAddress string, storageLocations []string) {
	if validatorAddress == "" {
		announcedStorageLocations, err := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper).AnnouncedStorageLocations(s.Ctx(), &types.QueryAnnouncedStorageLocationsRequest{MailboxId: mailboxId, ValidatorAddress: validatorAddress})
		Expect(err).To(BeNil())

		Expect(announcedStorageLocations.StorageLocations).To(HaveLen(0))
	} else {
		announcedStorageLocations, err := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper).AnnouncedStorageLocations(s.Ctx(), &types.QueryAnnouncedStorageLocationsRequest{MailboxId: mailboxId, ValidatorAddress: validatorAddress})
		Expect(err).To(BeNil())

		Expect(announcedStorageLocations.StorageLocations).To(Equal(storageLocations))

		latestAnnouncedStorageLocation, err := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper).LatestAnnouncedStorageLocation(s.Ctx(), &types.QueryLatestAnnouncedStorageLocationRequest{MailboxId: mailboxId, ValidatorAddress: validatorAddress})
		Expect(err).To(BeNil())
		Expect(latestAnnouncedStorageLocation.StorageLocation).To(Equal(storageLocations[len(storageLocations)-1]))
	}
}
