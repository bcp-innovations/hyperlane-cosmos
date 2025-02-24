package keeper_test

//
//import (
//	"crypto/ecdsa"
//	"fmt"
//
//	"github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
//
//	"cosmossdk.io/collections"
//
//	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
//	"github.com/bcp-innovations/hyperlane-cosmos/util"
//	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
//	"github.com/ethereum/go-ethereum/crypto"
//	. "github.com/onsi/ginkgo/v2"
//	. "github.com/onsi/gomega"
//)
//
///*
//
//TEST CASES - msg_announce.go
//
//* AnnounceValidator (invalid) with empty validator
//* AnnounceValidator (invalid) with invalid validator
//* AnnounceValidator (invalid) with empty storage location
//* AnnounceValidator (invalid) with empty signature
//* AnnounceValidator (invalid) with invalid signature
//* AnnounceValidator (invalid) with invalid signature recovery id
//* AnnounceValidator (invalid) same storage location for validator (replay protection)
//* AnnounceValidator (invalid) for non-existing Mailbox ID
//* AnnounceValidator (invalid) for invalid Mailbox ID
//* AnnounceValidator (invalid) for non-matching signature validator pair
//* AnnounceValidator (valid)
//* AnnounceValidator (valid) add another storage location
//
//*/
//
//var _ = Describe("msg_announce.go", Ordered, func() {
//	var s *i.KeeperTestSuite
//	var creator i.TestValidatorAddress
//
//	BeforeEach(func() {
//		s = i.NewCleanChain()
//		creator = i.GenerateTestValidatorAddress("Creator")
//		err := s.MintBaseCoins(creator.Address, 1_000_000)
//		Expect(err).To(BeNil())
//	})
//
//	It("AnnounceValidator (invalid) with empty validator", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       "",
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("validator cannot be empty"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) with invalid validator", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		invalidValidatorAddress := "0x0b1caf89d1edb9ee161093test94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       invalidValidatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("invalid validator address"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) with empty storage location", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: "",
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("storage location cannot be empty"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) with empty signature", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		storageLocation := "aws://key.pub"
//
//		// Act
//		_, err := s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       "",
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("signature cannot be empty"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) with invalid signature", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		storageLocation := "aws://key.pub"
//
//		// Act
//		_, err := s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       "0x0b1caf89d1edb9ee161093b1ec94ca75611dtest",
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("invalid signature"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) with invalid signature recovery id", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, true)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("invalid signature recovery id"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) same storage location for validator (replay protection)", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//		Expect(err).To(BeNil())
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal(fmt.Sprintf("validator %s already announced storage location %s", validatorAddress, storageLocation)))
//		validateAnnouncement(s, validatorAddress, []types.StorageLocation{{Location: storageLocation, Id: 0}})
//	})
//
//	It("AnnounceValidator (invalid) for non-existing Mailbox ID", func() {
//		// Arrange
//		nonExistingMailboxId, err := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
//		Expect(err).To(BeNil())
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, nonExistingMailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       nonExistingMailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId.String())))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) for invalid Mailbox ID", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String() + "test",
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal("invalid mailbox id"))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (invalid) for non-matching signature validator pair", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		correctValidatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		wrongValidatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db392"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       wrongValidatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err.Error()).To(Equal(fmt.Sprintf("validator %s doesn't match signature. recovered address: %s", wrongValidatorAddress, correctValidatorAddress)))
//		validateAnnouncement(s, "", []types.StorageLocation{})
//	})
//
//	It("AnnounceValidator (valid)", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err).To(BeNil())
//		validateAnnouncement(s, validatorAddress, []types.StorageLocation{{Location: storageLocation, Id: 0}})
//	})
//
//	It("AnnounceValidator (valid) add another storage location", func() {
//		// Arrange
//		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop", true, 1)
//
//		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
//		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
//		storageLocation := "aws://key.pub"
//
//		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx())
//		Expect(err).To(BeNil())
//
//		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)
//
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//		Expect(err).To(BeNil())
//
//		storageLocation2 := "aws://key2.pub"
//		signature = announce(validatorPrivKey, storageLocation2, mailboxId, localDomain, false)
//
//		// Act
//		_, err = s.RunTx(&types.MsgAnnounceValidator{
//			Validator:       validatorAddress,
//			StorageLocation: storageLocation2,
//			Signature:       signature,
//			MailboxId:       mailboxId.String(),
//			Creator:         creator.Address,
//		})
//
//		// Assert
//		Expect(err).To(BeNil())
//		validateAnnouncement(s, validatorAddress, []types.StorageLocation{
//			{Location: storageLocation, Id: 0},
//			{Location: storageLocation2, Id: 1},
//		})
//	})
//})
//
//func announce(privKey, storageLocation string, mailboxId util.HexAddress, localDomain uint32, skipRecoveryId bool) string {
//	announcementDigest := types.GetAnnouncementDigest(storageLocation, localDomain, mailboxId.Bytes())
//
//	ethDigest := util.GetEthSigningHash(announcementDigest[:])
//
//	privateKey, err := crypto.HexToECDSA(privKey)
//	Expect(err).To(BeNil())
//
//	publicKey := privateKey.Public()
//	_, ok := publicKey.(*ecdsa.PublicKey)
//	Expect(ok).To(BeTrue())
//
//	signedAnnouncement, err := crypto.Sign(ethDigest[:], privateKey)
//	Expect(err).To(BeNil())
//
//	if !skipRecoveryId {
//		// Required for recovery ID
//		// https://eips.ethereum.org/EIPS/eip-155
//		signedAnnouncement[64] += 27
//	}
//
//	return util.EncodeEthHex(signedAnnouncement)
//}
//
//func validateAnnouncement(s *i.KeeperTestSuite, validatorAddress string, storageLocations []types.StorageLocation) {
//	if validatorAddress == "" {
//		it, err := s.App().HyperlaneKeeper.Validators.Iterate(s.Ctx(), nil)
//		Expect(err).To(BeNil())
//
//		validators, err := it.Values()
//		Expect(err).To(BeNil())
//		Expect(validators).To(HaveLen(0))
//
//		rng := collections.NewPrefixedPairRange[[]byte, uint64](nil)
//
//		iter, err := s.App().HyperlaneKeeper.StorageLocations.Iterate(s.Ctx(), rng)
//		Expect(err).To(BeNil())
//
//		announcedStorageLocations, err := iter.Values()
//		Expect(err).To(BeNil())
//		Expect(announcedStorageLocations).To(HaveLen(0))
//	} else {
//		validatorAddressBytes, err := util.DecodeEthHex(validatorAddress)
//		Expect(err).To(BeNil())
//
//		_, err = s.App().HyperlaneKeeper.Validators.Get(s.Ctx(), validatorAddressBytes)
//		Expect(err).To(BeNil())
//
//		rng := collections.NewPrefixedPairRange[[]byte, uint64](validatorAddressBytes)
//
//		iter, err := s.App().HyperlaneKeeper.StorageLocations.Iterate(s.Ctx(), rng)
//		Expect(err).To(BeNil())
//
//		announcedStorageLocations, err := iter.Values()
//		Expect(err).To(BeNil())
//		Expect(storageLocations).To(Equal(announcedStorageLocations))
//
//		latestAnnouncedStorageLocation, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).LatestAnnouncedStorageLocation(s.Ctx(), &types.QueryLatestAnnouncedStorageLocationRequest{ValidatorAddress: validatorAddress})
//		Expect(err).To(BeNil())
//		Expect(latestAnnouncedStorageLocation.StorageLocation).To(Equal(storageLocations[len(storageLocations)-1].Location))
//	}
//}
