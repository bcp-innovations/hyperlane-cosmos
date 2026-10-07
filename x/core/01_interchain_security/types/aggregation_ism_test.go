package types_test

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - aggregation_ism.go

* Verify (invalid) - must be handled by the AggregationISMHandler
* Validate (valid)
* Validate (invalid) - too many modules
* ParseAggregationMetadata (valid) relayer encoding
* ParseAggregationMetadata (valid) empty sub-metadata
* ParseAggregationMetadata (invalid) range exceeding uint32 metadata length
*/

var _ = Describe("aggregation_ism.go", Ordered, func() {
	It("Verify (invalid) - must be handled by the AggregationISMHandler", func() {
		// Arrange
		ism := types.AggregationISM{}

		// Act
		result, err := ism.Verify(nil, []byte{}, util.HyperlaneMessage{})

		// Assert
		Expect(err.Error()).To(Equal("Verify should not be called on AggregationISM: unexpected error"))
		Expect(result).To(BeFalse())
	})

	It("Validate (valid)", func() {
		// Arrange
		ism := types.AggregationISM{
			Modules:   []util.HexAddress{util.CreateMockHexAddress("ism", 1), util.CreateMockHexAddress("ism", 2)},
			Threshold: 2,
		}

		// Act
		err := ism.Validate()

		// Assert
		Expect(err).To(BeNil())
	})

	It("Validate (invalid) - too many modules", func() {
		// Arrange
		ism := types.AggregationISM{Threshold: 1}
		for k := 0; k <= types.MaxAggregationModules; k++ {
			ism.Modules = append(ism.Modules, util.CreateMockHexAddress("ism", int64(k)))
		}

		// Act
		err := ism.Validate()

		// Assert
		Expect(err.Error()).To(Equal("too many modules: got 256, max 255"))
	})

	It("ParseAggregationMetadata (valid) relayer encoding", func() {
		// Arrange
		// Encoding as produced by the relayer's `format_metadata` for 3 modules where
		// modules 0 and 2 are provided: an 8 byte range per module followed by the metadata.
		metadata := bytesFromHexString(
			"00000018" + "0000001a" + // module 0: [24:26]
				"00000000" + "00000000" + // module 1: not provided
				"0000001a" + "0000001d" + // module 2: [26:29]
				"aaaa" + "bbbbbb",
		)

		// Act
		parsed, err := types.ParseAggregationMetadata(metadata, 3)

		// Assert
		Expect(err).To(BeNil())
		Expect(parsed).To(Equal([][]byte{{0xaa, 0xaa}, nil, {0xbb, 0xbb, 0xbb}}))
		Expect(types.FormatAggregationMetadata(parsed)).To(Equal(metadata))
	})

	It("ParseAggregationMetadata (valid) empty sub-metadata", func() {
		// Arrange
		metadata := types.FormatAggregationMetadata([][]byte{{}, nil})

		// Act
		parsed, err := types.ParseAggregationMetadata(metadata, 2)

		// Assert
		Expect(err).To(BeNil())
		Expect(parsed[0]).To(Equal([]byte{}))
		Expect(parsed[1]).To(BeNil())
	})

	It("ParseAggregationMetadata (invalid) range exceeding uint32 metadata length", func() {
		// Arrange
		metadata := bytesFromHexString("00000008" + "ffffffff")

		// Act
		_, err := types.ParseAggregationMetadata(metadata, 1)

		// Assert
		Expect(err.Error()).To(Equal("invalid metadata range for module 0: [8:4294967295] with metadata length 8"))
	})
})
