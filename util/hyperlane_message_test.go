package util_test

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/ethereum/go-ethereum/common/hexutil"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - hex_address.go

*

*/

var _ = Describe("hyperlane_message.go", Ordered, func() {

	BeforeEach(func() {
	})

	validWarpMessage := "0x" +
		"03" + // Version
		"00000003" + // Nonce
		"00000001" + // Origin
		"b32677d8121a50c7b960b8561ead86278a7d75ec786807983e1eebfcbc2d9cfc" + // Sender
		"00000507" + // Destination Domain
		"000000000000000000000000f254e1ce6b468e5c118214d13faa263011046715" // Recipient

	It("Decode (valid) Empty Hyperlane Message", func() {
		// Arrange

		// Act
		message, err := util.ParseHyperlaneMessage(hexutil.MustDecode(validWarpMessage))

		// Assert
		Expect(err).To(BeNil())
		Expect(message.Version).To(Equal(uint8(3)))
		Expect(message.Nonce).To(Equal(uint32(3)))
		Expect(message.Origin).To(Equal(uint32(1)))
		Expect(message.Sender.String()).To(Equal("0xb32677d8121a50c7b960b8561ead86278a7d75ec786807983e1eebfcbc2d9cfc"))
		Expect(message.Destination).To(Equal(uint32(1287)))
		Expect(message.Recipient.String()).To(Equal("0x000000000000000000000000f254e1ce6b468e5c118214d13faa263011046715"))

		Expect(len(message.Body)).To(Equal(0))

		// Keccak-256 calculated with different tool
		Expect(message.Id().String()).To(Equal("0xec919468b164c02e27b82c3f4c59ca10c401c0e99af2b333f47a80752cf4d481"))

		Expect(message.String()).To(Equal(validWarpMessage))
	})

	It("Decode (valid) Hyperlane Warp Message", func() {
		// Arrange
		rawMessage := validWarpMessage +
			"0000000000000000000000000c60e7ecd06429052223c78452f791aab5c5cac6" + // WarpPayload: Recipient
			"000000000000000000000000000000000000000000000000000000000000000b" // WarpPayload: amount

		// Act
		message, err := util.ParseHyperlaneMessage(hexutil.MustDecode(rawMessage))

		// Assert
		Expect(err).To(BeNil())
		Expect(message.Version).To(Equal(uint8(3)))
		Expect(message.Nonce).To(Equal(uint32(3)))
		Expect(message.Origin).To(Equal(uint32(1)))
		Expect(message.Sender.String()).To(Equal("0xb32677d8121a50c7b960b8561ead86278a7d75ec786807983e1eebfcbc2d9cfc"))
		Expect(message.Destination).To(Equal(uint32(1287)))
		Expect(message.Recipient.String()).To(Equal("0x000000000000000000000000f254e1ce6b468e5c118214d13faa263011046715"))

		Expect(len(message.Body)).To(Equal(64))

		// Keccak-256 calculated with different tool
		Expect(message.Id().String()).To(Equal("0xe9c31f5e798ea079b6c7193020fe8a665445bbfec0e514354494651574ac8c1f"))

		Expect(message.String()).To(Equal(rawMessage))
	})

	It("Decode (invalid) Hyperlane Warp Message (too short)", func() {
		// Arrange
		tooShortMessage := validWarpMessage[:len(validWarpMessage)-2]

		// Act
		message, err := util.ParseHyperlaneMessage(hexutil.MustDecode(tooShortMessage))

		// Assert
		Expect(err.Error()).To(Equal("invalid hyperlane message"))
		Expect(message.Version).To(Equal(uint8(0)))
		Expect(message.Nonce).To(Equal(uint32(0)))
		Expect(message.Origin).To(Equal(uint32(0)))
		Expect(message.Sender.String()).To(Equal(util.NewZeroAddress().String()))
		Expect(message.Destination).To(Equal(uint32(0)))
		Expect(message.Recipient.String()).To(Equal(util.NewZeroAddress().String()))

		Expect(len(message.Body)).To(Equal(0))
	})
})
