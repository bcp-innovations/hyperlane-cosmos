package types_test

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - aggregation_ism.go

* ValidateAggregationISM (invalid) threshold is zero
* ValidateAggregationISM (invalid) empty modules list
* ValidateAggregationISM (invalid) threshold exceeds number of modules
* ValidateAggregationISM (invalid) duplicate modules
* ValidateAggregationISM (valid) single module with threshold 1
* ValidateAggregationISM (valid) multiple modules with threshold less than count
* ValidateAggregationISM (valid) multiple modules with threshold equal to count

*/

var _ = Describe("aggregation_ism.go", Ordered, func() {
	It("ValidateAggregationISM (invalid) threshold is zero", func() {
		// Arrange
		module1 := util.CreateMockHexAddress("module", 1)
		modules := []util.HexAddress{module1}

		// Act
		err := types.ValidateAggregationISM(modules, 0)

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("threshold must be greater than 0"))
	})

	It("ValidateAggregationISM (invalid) empty modules list", func() {
		// Arrange
		modules := []util.HexAddress{}

		// Act
		err := types.ValidateAggregationISM(modules, 1)

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("modules list cannot be empty"))
	})

	It("ValidateAggregationISM (invalid) threshold exceeds number of modules", func() {
		// Arrange
		module1 := util.CreateMockHexAddress("module", 1)
		module2 := util.CreateMockHexAddress("module", 2)
		modules := []util.HexAddress{module1, module2}

		// Act
		err := types.ValidateAggregationISM(modules, 3)

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("threshold (3) cannot exceed number of modules (2)"))
	})

	It("ValidateAggregationISM (invalid) duplicate modules", func() {
		// Arrange
		module1 := util.CreateMockHexAddress("module", 1)
		modules := []util.HexAddress{module1, module1} // duplicate

		// Act
		err := types.ValidateAggregationISM(modules, 1)

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("duplicate module found"))
	})

	It("ValidateAggregationISM (valid) single module with threshold 1", func() {
		// Arrange
		module1 := util.CreateMockHexAddress("module", 1)
		modules := []util.HexAddress{module1}

		// Act
		err := types.ValidateAggregationISM(modules, 1)

		// Assert
		Expect(err).To(BeNil())
	})

	It("ValidateAggregationISM (valid) multiple modules with threshold less than count", func() {
		// Arrange
		module1 := util.CreateMockHexAddress("module", 1)
		module2 := util.CreateMockHexAddress("module", 2)
		module3 := util.CreateMockHexAddress("module", 3)
		modules := []util.HexAddress{module1, module2, module3}

		// Act
		err := types.ValidateAggregationISM(modules, 2) // 2 out of 3

		// Assert
		Expect(err).To(BeNil())
	})

	It("ValidateAggregationISM (valid) multiple modules with threshold equal to count", func() {
		// Arrange
		module1 := util.CreateMockHexAddress("module", 1)
		module2 := util.CreateMockHexAddress("module", 2)
		module3 := util.CreateMockHexAddress("module", 3)
		modules := []util.HexAddress{module1, module2, module3}

		// Act
		err := types.ValidateAggregationISM(modules, 3) // 3 out of 3 (all must pass)

		// Assert
		Expect(err).To(BeNil())
	})

	It("ValidateAggregationISM (valid) large number of modules", func() {
		// Arrange
		var modules []util.HexAddress
		for i := 0; i < 100; i++ {
			modules = append(modules, util.CreateMockHexAddress("module", int64(i)))
		}

		// Act
		err := types.ValidateAggregationISM(modules, 50) // 50 out of 100

		// Assert
		Expect(err).To(BeNil())
	})
})
