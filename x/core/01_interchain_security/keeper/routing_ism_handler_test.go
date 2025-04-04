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

TEST CASES - routing_ism_handler_test.go

* Verify (valid)
* Verify (valid) with nested RoutingISM
* Verify (valid) - multiple registered routes
* Verify (invalid) - non-existing ISM
* Verify (invalid) - message for non-enrolled domain
* Verify (invalid) overflow due to self reference with large gas
* Verify (valid) gas costs for 1000 registered routes
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

	createRoutingIsm := func() util.HexAddress {
		res, err := s.RunTx(&types.MsgCreateRoutingIsm{
			Creator: creator.Address,
		})

		Expect(err).To(BeNil())

		var routingIsm types.MsgCreateRoutingIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &routingIsm)
		Expect(err).To(BeNil())

		return routingIsm.Id
	}

	setRoute := func(routingIsm, ism util.HexAddress, domain uint32) {
		// Act
		_, err := s.RunTx(&types.MsgSetRoutingIsmDomain{
			Owner: creator.Address,
			IsmId: routingIsm,
			Route: types.Route{
				Ism:    ism,
				Domain: domain,
			},
		})

		// Assert
		Expect(err).To(BeNil())
	}

	It("Verify (valid)", func() {
		// Arrange

		// registry mock ISM
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		mockIsmId, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		routingIsm := createRoutingIsm()
		setRoute(routingIsm, mockIsmId, 1)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), routingIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		// verify mock ISM was called
		Expect(mockIsm.CallCount()).To(Equal(1))
	})

	It("Verify (valid) with nested RoutingISM", func() {
		// Arrange

		// registry mock ISM
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		mockIsmId, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		routingIsmB := createRoutingIsm()
		setRoute(routingIsmB, mockIsmId, 1)

		routingIsmA := createRoutingIsm()
		setRoute(routingIsmA, routingIsmB, 1)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), routingIsmA, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		// verify mock ISM was called
		Expect(mockIsm.CallCount()).To(Equal(1))
	})

	It("Verify (valid) - multiple registered routes", func() {
		// Arrange

		// registry mock ISM
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		mockIsmId1, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())
		mockIsmId2, err := mockIsm.RegisterIsm(s.Ctx())
		Expect(err).To(BeNil())

		routingIsm := createRoutingIsm()
		setRoute(routingIsm, mockIsmId1, 1)
		setRoute(routingIsm, mockIsmId2, 2)

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), routingIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())

		// verify mock ISM was called
		// only one route should be called not both of them
		Expect(mockIsm.CallCount()).To(Equal(1))
	})

	It("Verify (invalid) - non-existing ISM", func() {
		// Arrange
		routingIsm := createRoutingIsm()

		// create invalid ism id
		routingIsm[31] = 10

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), routingIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err.Error()).To(Equal("collections: not found: key '10' of type <nil>"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) - message for non-enrolled domain", func() {
		// Arrange
		routingIsm := createRoutingIsm()

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(s.Ctx(), routingIsm, []byte{}, util.HyperlaneMessage{
			Origin: 1,
		})

		// Assert
		Expect(err.Error()).To(Equal("no route found for domain 1: no route found"))
		Expect(result).To(BeFalse())
	})

	It("Verify (invalid) overflow due to self reference with large gas", func() {
		// Arrange
		routingIsm := createRoutingIsm()
		setRoute(routingIsm, routingIsm, 1)

		message := util.HyperlaneMessage{
			Origin: 1,
		}

		// set an explict gas limit
		ctx := s.Ctx().WithGasMeter(storetypes.NewGasMeter(1_000_000_000)).WithBlockGasMeter(storetypes.NewGasMeter(1_000_000_000))

		hasPanicked := false

		act := func() {
			defer func() {
				if r := recover(); r != nil {
					hasPanicked = true
				}
			}()
			// satisfy linter
			result, err := s.App().HyperlaneKeeper.Verify(ctx, routingIsm, []byte{}, message)
			Expect(err).To(BeNil())
			Expect(result).To(BeTrue())
		}
		// Act
		act()

		Expect(hasPanicked).To(BeTrue())
	})

	It("Verify (valid) gas costs for 1000 registered routes", func() {
		// Arrange
		routingIsm := createRoutingIsm()
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		for k := uint32(1); k <= 1000; k++ {
			ism, err := mockIsm.RegisterIsm(s.Ctx())
			Expect(err).To(BeNil())
			setRoute(routingIsm, ism, k)
		}

		message := util.HyperlaneMessage{
			Origin: 1000,
		}

		// set an explict gas limit
		ctx := s.Ctx().WithGasMeter(storetypes.NewGasMeter(300_000))

		// Act
		result, err := s.App().HyperlaneKeeper.Verify(ctx, routingIsm, []byte{}, message)

		// Assert
		Expect(err).To(BeNil())
		Expect(result).To(BeTrue())
		Expect(mockIsm.CallCount()).To(Equal(1))
		Expect(ctx.GasMeter().GasConsumed()).Should(BeNumerically("<", 300_000))
	})
})
