package types_test

import (
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("aggregation_ism.go", Ordered, func() {
	It("ParseAggregationMetadata (valid) relayer encoding", func() {
		// Encoding as produced by the relayer's `format_metadata` for 3 modules where
		// modules 0 and 2 are provided: an 8 byte range per module followed by the metadata.
		metadata := bytesFromHexString(
			"00000018" + "0000001a" + // module 0: [24:26]
				"00000000" + "00000000" + // module 1: not provided
				"0000001a" + "0000001d" + // module 2: [26:29]
				"aaaa" + "bbbbbb",
		)

		parsed, err := types.ParseAggregationMetadata(metadata, 3)

		Expect(err).To(BeNil())
		Expect(parsed).To(Equal([][]byte{{0xaa, 0xaa}, nil, {0xbb, 0xbb, 0xbb}}))
		Expect(types.FormatAggregationMetadata(parsed)).To(Equal(metadata))
	})

	It("ParseAggregationMetadata (valid) empty sub-metadata", func() {
		metadata := types.FormatAggregationMetadata([][]byte{{}, nil})

		parsed, err := types.ParseAggregationMetadata(metadata, 2)

		Expect(err).To(BeNil())
		Expect(parsed[0]).To(Equal([]byte{}))
		Expect(parsed[1]).To(BeNil())
	})

	It("ParseAggregationMetadata (invalid) range exceeding uint32 metadata length", func() {
		metadata := bytesFromHexString("00000008" + "ffffffff")

		_, err := types.ParseAggregationMetadata(metadata, 1)

		Expect(err.Error()).To(Equal("invalid metadata range for module 0: [8:4294967295] with metadata length 8"))
	})
})
