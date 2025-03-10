package util_test

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - hex_address.go

* Decode (valid) Zero Hex Address
* Decode (invalid) Hex Address (to short)
* Decode (invalid) Hex Address (invalid hex)
* Address generation
* Proto marshalling
* Proto marshalling 2
* Proto marshalling 2 (wrong length)
* Proto marshalling (JSON)
* Proto unmarshalling
* Proto unmarshalling (wrong length)
* Proto unmarshalling (invalid)
* Proto unmarshalling (JSON)
* Proto unmarshalling (JSON) (invalid json)
* Proto unmarshalling (JSON) (invalid hex address)

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

	It("Decode (invalid) Hex Address (invalid hex)", func() {
		// Arrange
		invalidZeroHex := "0x000000000000000000000000000000000000000000000000000000000000000g"

		// Act
		address, err := util.DecodeHexAddress(invalidZeroHex)

		// Assert (Address equals Zero address)
		Expect(err.Error()).To(Equal("encoding/hex: invalid byte: U+0067 'g'"))
		Expect(address.IsZeroAddress()).To(BeTrue())
		Expect(address.GetInternalId()).To(Equal(uint64(0)))
		Expect(address.String()).To(Equal(zeroHex))
	})

	It("Address generation", func() {
		// Arrange

		// Act
		identifier := make([]byte, 20)
		copy(identifier, "hyperlane")
		address := util.GenerateHexAddress([20]byte(identifier), 1, 1)

		// Assert
		Expect(address.String()).To(Equal("0x68797065726c616e650000000000000000000000000000010000000000000001"))
		Expect(address.GetInternalId()).To(Equal(uint64(1)))
		Expect(address.GetType()).To(Equal(uint32(1)))
	})

	It("Proto marshalling", func() {
		// Arrange
		identifier := make([]byte, 20)
		copy(identifier, "hyperlane")
		address := util.GenerateHexAddress([20]byte(identifier), 1, 1)

		// Act
		bytes, err := address.Marshal()
		Expect(err).To(BeNil())

		// Assert
		Expect(string(bytes)).To(Equal("0x68797065726c616e650000000000000000000000000000010000000000000001"))
		Expect(address.Size()).To(Equal(2 + 2*32))
	})

	It("Proto marshalling 2", func() {
		// Arrange
		identifier := make([]byte, 20)
		copy(identifier, "hyperlane")
		address := util.GenerateHexAddress([20]byte(identifier), 1, 1)

		bytes := make([]byte, 66)
		// Act
		n, err := address.MarshalTo(bytes)
		Expect(err).To(BeNil())

		// Assert
		Expect(n).To(Equal(66))
		Expect(string(bytes)).To(Equal("0x68797065726c616e650000000000000000000000000000010000000000000001"))
		Expect(address.Size()).To(Equal(2 + 2*32))
	})

	It("Proto marshalling 2 (wrong length)", func() {
		// Arrange
		identifier := make([]byte, 20)
		copy(identifier, "hyperlane")
		address := util.GenerateHexAddress([20]byte(identifier), 1, 1)

		bytes := make([]byte, 64)
		// Act
		n, err := address.MarshalTo(bytes)

		// Assert
		Expect(err.Error()).To(Equal("invalid hex address length: 64"))
		Expect(n).To(Equal(64))
		Expect(address.Size()).To(Equal(2 + 2*32))
	})

	It("Proto marshalling (JSON)", func() {
		// Arrange
		rawAddress := "0x68797065726c616e650000000000000000000000000000010000000000000001"
		address, _ := util.DecodeHexAddress(rawAddress)

		// Act
		bytes, err := address.MarshalJSON()
		Expect(err).To(BeNil())

		// Assert
		Expect(string(bytes)).To(Equal(fmt.Sprintf(`"%s"`, rawAddress)))
		Expect(address.Size()).To(Equal(2 + 2*32))
	})

	It("Proto unmarshalling", func() {
		// Arrange
		rawAddress := "0x68797065726c616e650000000000000000000000000000010000000000000001"
		var address util.HexAddress

		// Act
		err := address.Unmarshal([]byte(rawAddress))
		Expect(err).To(BeNil())

		// Assert
		Expect(address.String()).To(Equal("0x68797065726c616e650000000000000000000000000000010000000000000001"))
		Expect(address.Size()).To(Equal(2 + 2*32))
	})

	It("Proto unmarshalling (wrong length)", func() {
		// Arrange
		rawAddress := "0x68797065726c616e65000000000000000000000000000001000000000000000"
		var address util.HexAddress

		// Act
		err := address.Unmarshal([]byte(rawAddress))

		// Assert
		Expect(err.Error()).To(Equal("invalid hex address length"))
		Expect(address.String()).To(Equal(zeroHex))
	})

	It("Proto unmarshalling (invalid)", func() {
		// Arrange
		rawAddress := "0x68797065726c616e65000000000000000000000000000001000000000000000g"
		var address util.HexAddress

		// Act
		err := address.Unmarshal([]byte(rawAddress))

		// Assert
		Expect(err.Error()).To(Equal("encoding/hex: invalid byte: U+0067 'g'"))
		Expect(address.String()).To(Equal(zeroHex))
	})

	It("Proto unmarshalling (JSON)", func() {
		// Arrange
		rawAddress := "\"0x68797065726c616e650000000000000000000000000000010000000000000001\""
		var address util.HexAddress

		// Act
		err := address.UnmarshalJSON([]byte(rawAddress))
		Expect(err).To(BeNil())

		// Assert
		Expect(address.String()).To(Equal("0x68797065726c616e650000000000000000000000000000010000000000000001"))
		Expect(address.Size()).To(Equal(2 + 2*32))
	})

	It("Proto unmarshalling (JSON) (invalid json)", func() {
		// Arrange
		rawAddress := "\"0x68797065726c616e650000000000000000000000000000010000000000000001"
		var address util.HexAddress

		// Act
		err := address.UnmarshalJSON([]byte(rawAddress))

		// Assert
		Expect(err.Error()).To(Equal("unexpected end of JSON input"))
		Expect(address.String()).To(Equal(util.NewZeroAddress().String()))
	})

	It("Proto unmarshalling (JSON) (invalid hex address)", func() {
		// Arrange
		rawAddress := "\"0x68797065726c616e65000000000000000000000000000001000000000000000g\""
		var address util.HexAddress

		// Act
		err := address.UnmarshalJSON([]byte(rawAddress))

		// Assert
		Expect(err.Error()).To(Equal("encoding/hex: invalid byte: U+0067 'g'"))
		Expect(address.String()).To(Equal(util.NewZeroAddress().String()))
	})
})
