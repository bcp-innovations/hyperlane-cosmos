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
* Verify (invalid) no registered route
* Verify (invalid) registered itself
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

	It("Verify (invalid) overflow due to self reference", func() {
		// Arrange

		// registry mock ISM
		mockIsm := i.CreateMockIsm(s.App().HyperlaneKeeper.IsmRouter())

		routingIsm := createRoutingIsm()
		setRoute(routingIsm, routingIsm, 1)

		message := util.HyperlaneMessage{
			Origin: 1,
		}

		// set an explicity gas limit
		ctx := s.Ctx().WithGasMeter(storetypes.NewGasMeter(200000))

		panic := false

		act := func() {
			defer func() {
				if r := recover(); r != nil {
					panic = true
				}
			}()
			// satisfy linter
			result, err := s.App().HyperlaneKeeper.Verify(ctx, routingIsm, []byte{}, message)
			Expect(err).To(BeNil())
			Expect(result).To(BeTrue())
		}
		// Act
		act()

		Expect(panic).To(BeTrue())

		// verify mock ISM was called
		Expect(mockIsm.CallCount()).To(Equal(0))
	})
})
