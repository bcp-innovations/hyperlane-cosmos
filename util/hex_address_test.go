package util_test

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - hex_address.go

*

*/

var _ = Describe("hex_address.go", Ordered, func() {

	BeforeEach(func() {
	})

	zeroHex := "0x0000000000000000000000000000000000000000000000000000000000000000"

	It("Decode (valid) Zero Hex Address", func() {
		// Arrange

		// Act
		address, err := util.DecodeHexAddress(zeroHex)

		// Assert
		Expect(err).To(BeNil())
		Expect(address.IsZeroAddress()).To(BeTrue())
		Expect(address.GetInternalId()).To(Equal(uint64(0)))
		Expect(address.String()).To(Equal(zeroHex))
	})

	It("Decode (invalid) Hex Address (to short)", func() {
		// Arrange
		invalidZeroHex := "0x000000000000000000000000000000000000000000000000000000000000000"

		// Act
		address, err := util.DecodeHexAddress(invalidZeroHex)

		// Assert (Address equals Zero address)
		Expect(err.Error()).To(Equal("invalid hex address length"))
		Expect(address.IsZeroAddress()).To(BeTrue())
		Expect(address.GetInternalId()).To(Equal(uint64(0)))
		Expect(address.String()).To(Equal(zeroHex))
	})
})
