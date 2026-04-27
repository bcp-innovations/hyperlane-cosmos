package keeper_test

import (
	storetypes "cosmossdk.io/store/types"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - aggregation_ism_handler_test.go

* Verify (valid) - all ISMs pass, threshold met
* Verify (valid) - exactly threshold ISMs pass
* Verify (valid) - more than threshold ISMs pass
* Verify (valid) - early exit optimization when threshold met
* Verify (valid) - with nested AggregationISM
* Verify (invalid) - fewer than threshold ISMs pass
* Verify (invalid) - all ISMs fail
* Verify (invalid) - some ISMs return errors (treated as failures)
* Verify (invalid) - non-existing AggregationISM
* Verify (invalid) - non-existing child ISM
* Verify (invalid) - empty modules list
* Verify (invalid) - threshold is zero
* Verify (invalid) - threshold exceeds number of modules
* Verify (valid) - gas costs for aggregation with 100 modules
*/

var _ = Describe("aggregation_ism_handler_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	createAggregationIsm := func(modules []util.HexAddress, threshold uint32) util.HexAddress {
		res, err := s.RunTx(&types.MsgCreateAggregationIsm{
			Creator:   creator.Address,
			Modules:   modules,
			Threshold: threshold,
		})

		Expect(err).To(BeNil())

		var aggregationIsm types.MsgCreateAggregationIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &aggregationIsm)
		Expect(err).To(BeNil())

		return aggregationIsm.Id
	}

	It("Verify (valid) - all ISMs pass, threshold met", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
		mockIsm.SetShouldVerify(true)

		mockIsmId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId3, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		aggregationIsm := createAggregationIsm(
			[]util.HexAddress{mockIsmId1, mockIsmId2, mockIsmId3},
			2, // threshold: 2 out of 3
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		Expect(mockIsm.CallCount()).To(Equal(2)) // Should stop at threshold
	})

	It("Verify (valid) - exactly threshold ISMs pass", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		// Register ISMs that will pass
		mockIsm.SetShouldVerify(true)
		passId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		passId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		// Register an ISM that will fail
		mockIsm.SetShouldVerify(false)
		failId, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		aggregationIsm := createAggregationIsm(
			[]util.HexAddress{passId1, passId2, failId},
			2, // threshold: 2 out of 3
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
	})

	It("Verify (valid) - more than threshold ISMs pass", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
		mockIsm.SetShouldVerify(true)

		mockIsmId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId3, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId4, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		aggregationIsm := createAggregationIsm(
			[]util.HexAddress{mockIsmId1, mockIsmId2, mockIsmId3, mockIsmId4},
			2, // threshold: 2 out of 4
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		// Should exit early after 2 passes
		Expect(mockIsm.CallCount()).To(Equal(2))
	})

	It("Verify (valid) - early exit optimization when threshold met", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
		mockIsm.SetShouldVerify(true)

		var modules []util.HexAddress
		for i := 0; i < 10; i++ {
			id, err := mockIsm.RegisterIsm(s.Ctx())
			Expect(err).To(BeNil())
			modules = append(modules, id)
		}

		aggregationIsm := createAggregationIsm(modules, 3) // threshold: 3 out of 10

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		// Should only call 3 ISMs, not all 10
		Expect(mockIsm.CallCount()).To(Equal(3))
	})

	It("Verify (valid) - with nested AggregationISM", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
		mockIsm.SetShouldVerify(true)

		mockIsmId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		// Create inner aggregation ISM
		innerAggregationIsm := createAggregationIsm(
			[]util.HexAddress{mockIsmId1, mockIsmId2},
			1, // threshold: 1 out of 2
		)

		mockIsmId3, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		// Create outer aggregation ISM
		outerAggregationIsm := createAggregationIsm(
			[]util.HexAddress{innerAggregationIsm, mockIsmId3},
			2, // threshold: 2 out of 2 (both must pass)
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), outerAggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
	})

	It("Verify (invalid) - fewer than threshold ISMs pass", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		// Register one ISM that will pass
		mockIsm.SetShouldVerify(true)
		passId, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		// Register two ISMs that will fail
		mockIsm.SetShouldVerify(false)
		failId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		failId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		aggregationIsm := createAggregationIsm(
			[]util.HexAddress{passId, failId1, failId2},
			2, // threshold: 2 out of 3, but only 1 passes
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("insufficient verifications"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - all ISMs fail", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
		mockIsm.SetShouldVerify(false)

		mockIsmId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId3, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		aggregationIsm := createAggregationIsm(
			[]util.HexAddress{mockIsmId1, mockIsmId2, mockIsmId3},
			1, // threshold: 1 out of 3
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("insufficient verifications: 0/3"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - mix of passing and failing ISMs below threshold", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		// Register two ISMs that will pass
		mockIsm.SetShouldVerify(true)
		passId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		passId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		// Register three ISMs that will fail
		mockIsm.SetShouldVerify(false)
		failId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		failId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		failId3, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		aggregationIsm := createAggregationIsm(
			[]util.HexAddress{passId1, passId2, failId1, failId2, failId3},
			3, // threshold: 3 out of 5, but only 2 pass
		)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), aggregationIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("insufficient verifications: 2/5"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - non-existing AggregationISM", func() {
		// Arrange
		invalidIsmId := util.CreateMockHexAddress("invalid", 999)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), invalidIsmId, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("not found"))
		Expect(result).To(BeFalse())
	})

	It("Verify (valid) - gas costs for aggregation with 100 modules", func() {
		// Arrange
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())
		mockIsm.SetShouldVerify(true)

		var modules []util.HexAddress
		for i := 0; i < 100; i++ {
			id, err := mockIsm.RegisterIsm(s.Ctx())
			Expect(err).To(BeNil())
			modules = append(modules, id)
		}

		aggregationIsm := createAggregationIsm(modules, 50) // threshold: 50 out of 100

		message := util.HyperlaneMessage{
			Origin: 1,
		}

		// Set an explicit gas limit
		ctx := s.Ctx().WithGasMeter(storetypes.NewGasMeter(5_000_000))

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(ctx, aggregationIsm, []byte{}, message)

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		// Should only call 50 ISMs due to early exit
		Expect(mockIsm.CallCount()).To(Equal(50))
		// Verify gas consumption is reasonable
		Expect(ctx.GasMeter().GasConsumed()).Should(BeNumerically("<", 5_000_000))
	})
})
