package types_test

import (
	"fmt"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - message_id_multisig.go

* Validate (invalid) invalid validators
* Validate (invalid) invalid validator ethereum address
* Validate (invalid) unsorted validators
* Validate (invalid) duplicated validators
* Validate (invalid) too high threshold
* Validate (invalid) zero threshold
* Verify (invalid) empty metadata
+ Verify (invalid) invalid metadata - signatures length
* Verify (invalid) threshold can't be reached
* Verify (invalid) invalid signature
* Verify (invalid) wrong signature
* Verify (invalid) duplicated signature
* Verify (valid) multi-sig signature

*/

var _ = Describe("message_id_multisig.go", Ordered, func() {
	var s *i.KeeperTestSuite

	BeforeEach(func() {
		s = i.NewCleanChain()
	})

	It("Validate (invalid) invalid validators", func() {
		// Arrange
		validators := []string{
			"invalid1",
			"invalid2",
		}

		// Act
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		// Assert
		Expect(messageIdMultisigIsm.Validate().Error()).To(Equal(fmt.Sprintf("invalid validator address: %s", validators[0])))
	})

	It("Validate (invalid) invalid validator ethereum address", func() {
		// Arrange
		invalidValidator := make([]byte, 21)
		validators := []string{
			util.EncodeEthHex(invalidValidator),
		}

		// Act
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  1,
		}

		// Assert
		Expect("invalid validator address: must be 20 bytes").To(Equal(messageIdMultisigIsm.Validate().Error()))
	})

	It("Validate (invalid) unsorted validators", func() {
		// Arrange
		var validators []string
		for i := range PrivateKeys {
			validators = append(validators, PrivateKeys[len(PrivateKeys)-(i+1)].address)
		}

		// Act
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		// Assert
		Expect(messageIdMultisigIsm.Validate().Error()).To(Equal("validator addresses are not sorted correctly in ascending order"))
	})

	It("Validate (invalid) duplicated validators", func() {
		// Arrange
		var validators []string
		for range PrivateKeys {
			validators = append(validators, PrivateKeys[0].address)
		}

		// Act
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		// Assert
		Expect(messageIdMultisigIsm.Validate().Error()).To(Equal(fmt.Sprintf("duplicate validator address: %s", PrivateKeys[0].address)))
	})

	It("Validate (invalid) too high threshold", func() {
		// Arrange
		var validators []string
		for i := range PrivateKeys {
			validators = append(validators, PrivateKeys[len(PrivateKeys)-(i+1)].address)
		}

		// Act
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  4,
		}
		// Assert
		Expect(messageIdMultisigIsm.Validate().Error()).To(Equal("validator addresses less than threshold"))
	})

	It("Validate (invalid) zero threshold", func() {
		// Arrange
		var validators []string
		for i := range PrivateKeys {
			validators = append(validators, PrivateKeys[len(PrivateKeys)-(i+1)].address)
		}

		// Act
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  0,
		}
		// Assert
		Expect(messageIdMultisigIsm.Validate().Error()).To(Equal("threshold must be greater than zero"))
	})

	It("Verify (invalid) empty metadata", func() {
		// Arrange
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{"0x0c60e7eCd06429052223C78452F791AAb5C5CAc6"},
			Threshold:  1,
		}

		metadata := bytesFromHexString("")

		// Act
		verify, err := messageIdMultisigIsm.Verify(s.Ctx(), metadata, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("invalid metadata length: got 0, expected at least 68 bytes"))
		Expect(verify).To(BeFalse())
	})

	It("Verify (invalid) invalid metadata - signatures length", func() {
		// Arrange
		merkleRootMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{"0x0c60e7eCd06429052223C78452F791AAb5C5CAc6"},
			Threshold:  1,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MessageIdMultisigMetadata{
			MerkleTreeHook: [32]byte{},
			MerkleRoot:     [32]byte{},
			MerkleIndex:    0,
			SignatureCount: 0,
		}

		digest := metadata.Digest(&message)

		var signatures [][]byte
		sig := signDigest(digest[:], PrivateKeys[1].privateKey)
		signatures = append(signatures, sig)
		signatures = append(signatures, []byte{0})
		metadata.Signatures = signatures

		// Act
		verify, err := merkleRootMultisigIsm.Verify(s.Ctx(), metadata.Bytes(), message)

		// Assert
		Expect(err.Error()).To(Equal("invalid signatures length in metadata"))
		Expect(verify).To(BeFalse())
	})

	It("Verify (invalid) threshold can't be reached", func() {
		// Arrange
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{PrivateKeys[0].address, PrivateKeys[1].address},
			Threshold:  2,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MessageIdMultisigMetadata{
			MerkleTreeHook: [32]byte{},
			MerkleRoot:     [32]byte{},
			MerkleIndex:    uint32(0),
		}

		digest := metadata.Digest(&message)

		var signatures [][]byte
		sig := signDigest(digest[:], PrivateKeys[1].privateKey)
		signatures = append(signatures, sig)
		metadata.Signatures = signatures

		// Act
		verify, err := messageIdMultisigIsm.Verify(s.Ctx(), metadata.Bytes(), message)

		// Assert
		Expect(err.Error()).To(Equal("threshold can not be reached"))
		Expect(verify).To(BeFalse())
	})

	It("Verify (invalid) invalid signature", func() {
		// Arrange
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{PrivateKeys[0].address, PrivateKeys[1].address},
			Threshold:  1,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MessageIdMultisigMetadata{
			MerkleTreeHook: [32]byte{},
			MerkleRoot:     [32]byte{},
			MerkleIndex:    uint32(0),
		}

		var signatures [][]byte
		invalidSignature := make([]byte, 65)
		signatures = append(signatures, invalidSignature)
		metadata.Signatures = signatures

		// Act
		verify, err := messageIdMultisigIsm.Verify(s.Ctx(), metadata.Bytes(), message)

		// Assert
		Expect(err.Error()).To(Equal("failed to recover validator signature: invalid signature recovery id"))
		Expect(verify).To(BeFalse())
	})

	It("Verify (invalid) wrong signature", func() {
		// Arrange
		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: []string{PrivateKeys[0].address, PrivateKeys[1].address},
			Threshold:  1,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MessageIdMultisigMetadata{
			MerkleTreeHook: [32]byte{},
			MerkleRoot:     [32]byte{},
			MerkleIndex:    uint32(0),
		}

		digest := metadata.Digest(&message)

		var signatures [][]byte
		sig := signDigest(digest[:], PrivateKeys[2].privateKey)
		signatures = append(signatures, sig)
		metadata.Signatures = signatures

		// Act
		verify, err := messageIdMultisigIsm.Verify(s.Ctx(), metadata.Bytes(), message)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeFalse())
	})

	It("Verify (invalid) duplicated signature", func() {
		// Arrange
		var validators []string
		for _, validator := range PrivateKeys {
			validators = append(validators, validator.address)
		}

		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MessageIdMultisigMetadata{
			MerkleTreeHook: [32]byte{},
			MerkleRoot:     [32]byte{},
			MerkleIndex:    uint32(0),
		}

		digest := metadata.Digest(&message)

		var duplicatedSignatures [][]byte
		for range PrivateKeys {
			sig := signDigest(digest[:], PrivateKeys[0].privateKey)
			duplicatedSignatures = append(duplicatedSignatures, sig)
		}
		metadata.Signatures = duplicatedSignatures

		// Act
		verify, err := messageIdMultisigIsm.Verify(s.Ctx(), metadata.Bytes(), message)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeFalse())
	})

	It("Verify (valid) multi-sig signature", func() {
		// Arrange
		var validators []string
		for _, validator := range PrivateKeys {
			validators = append(validators, validator.address)
		}

		messageIdMultisigIsm := types.MessageIdMultisigISM{
			Id:         util.HexAddress{},
			Owner:      "",
			Validators: validators,
			Threshold:  2,
		}

		message := util.HyperlaneMessage{
			Version:     0,
			Nonce:       0,
			Origin:      0,
			Sender:      util.HexAddress{},
			Destination: 0,
			Recipient:   util.HexAddress{},
			Body:        nil,
		}

		metadata := types.MessageIdMultisigMetadata{
			MerkleTreeHook: [32]byte{},
			MerkleRoot:     [32]byte{},
			MerkleIndex:    uint32(0),
		}

		digest := metadata.Digest(&message)

		var signatures [][]byte
		for i := range PrivateKeys {
			sig := signDigest(digest[:], PrivateKeys[i].privateKey)
			signatures = append(signatures, sig)
		}
		metadata.Signatures = signatures

		// Act
		verify, err := messageIdMultisigIsm.Verify(s.Ctx(), metadata.Bytes(), message)

		// Assert
		Expect(err).To(BeNil())
		Expect(verify).To(BeTrue())
	})
})
