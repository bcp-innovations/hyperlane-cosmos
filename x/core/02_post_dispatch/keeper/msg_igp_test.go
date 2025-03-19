package keeper_test

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_igp.go

* CreateIgp (invalid) with invalid denom
* CreateIgp (valid)
* SetDestinationGasConfig (invalid) for non-existing IGP
* SetDestinationGasConfig (invalid) with wrong owner
* SetDestinationGasConfig (invalid) without gas oracle
* SetDestinationGasConfig (invalid) for invalid IGP
* SetDestinationGasConfig (valid)
* MsgSetIgpOwner (invalid) for invalid IGP
* MsgSetIgpOwner (invalid) for non-existing IGP
* MsgSetIgpOwner (invalid) called by non-owner
* MsgSetIgpOwner (valid)
*/

var _ = Describe("msg_igp.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var gasPayer i.TestValidatorAddress

	denom := "acoin"

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		gasPayer = i.GenerateTestValidatorAddress("Payer")

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		err = s.MintBaseCoins(gasPayer.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	// SetDestinationGasConfig
	It("SetDestinationGasConfig (invalid) for non-existing IGP", func() {
		// Arrange
		nonExistingIgp := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		_, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: nonExistingIgp,
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 1,
				GasOracle: &types.GasOracle{
					TokenExchangeRate: math.NewInt(1e10),
					GasPrice:          math.NewInt(1),
				},
				GasOverhead: math.NewInt(200000),
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp does not exist: %s", nonExistingIgp)))
	})

	It("SetDestinationGasConfig (invalid) with wrong owner", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: gasPayer.Address,
			IgpId: response.Id,
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 1,
				GasOracle: &types.GasOracle{
					TokenExchangeRate: math.NewInt(1e10),
					GasPrice:          math.NewInt(1),
				},
				GasOverhead: math.NewInt(200000),
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to set DestinationGasConfigs: %s is not the owner of igp with id %s", gasPayer.Address, response.Id)))
	})

	It("SetDestinationGasConfig (invalid) without gas oracle", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: "denom",
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: response.Id,
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 1,
				GasOracle:    nil,
				GasOverhead:  math.NewInt(200000),
			},
		})

		// Assert
		Expect(err.Error()).To(Equal("failed to set DestinationGasConfigs: gas Oracle is required"))
	})

	// IGP creation
	It("CreateIgp (invalid) with invalid denom", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: "123HYPERLANE!",
		})

		// Assert
		Expect(err.Error()).To(Equal("denom 123HYPERLANE! is invalid"))
	})

	It("CreateIgp (valid)", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.Owner).To(Equal(creator.Address))
		Expect(igp.Denom).To(Equal(denom))
		Expect(igp.ClaimableFees.IsZero()).To(BeTrue())
	})

	// SetDestinationGasConfig
	It("SetDestinationGasConfig (invalid) for invalid IGP", func() {
		// Arrange
		invalidIgp := "0x12345"

		_, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: invalidIgp,
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 1,
				GasOracle: &types.GasOracle{
					TokenExchangeRate: math.NewInt(1e10),
					GasPrice:          math.NewInt(1),
				},
				GasOverhead: math.NewInt(200000),
			},
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp id %s is invalid: invalid hex address length", invalidIgp)))
	})

	It("SetDestinationGasConfig (valid)", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: igpId.String(),
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 1,
				GasOracle: &types.GasOracle{
					TokenExchangeRate: math.NewInt(1e10),
					GasPrice:          math.NewInt(1),
				},
				GasOverhead: math.NewInt(200000),
			},
		})

		// Assert
		Expect(err).To(BeNil())

		rng := collections.NewPrefixedPairRange[uint64, uint32](igpId.GetInternalId())

		iter, err := s.App().HyperlaneKeeper.PostDispatchKeeper.IgpDestinationGasConfigs.Iterate(s.Ctx(), rng)
		Expect(err).To(BeNil())

		destinationGasConfigs, err := iter.Values()
		Expect(err).To(BeNil())

		Expect(destinationGasConfigs).To(HaveLen(1))
		Expect(destinationGasConfigs[0].RemoteDomain).To(Equal(uint32(1)))
		Expect(destinationGasConfigs[0].GasOracle.TokenExchangeRate).To(Equal(math.NewInt(1e10)))
		Expect(destinationGasConfigs[0].GasOracle.GasPrice).To(Equal(math.NewInt(1)))
		Expect(destinationGasConfigs[0].GasOverhead).To(Equal(math.NewInt(200000)))
	})

	// MsgSetIgpOwner
	It("MsgSetIgpOwner (invalid) for invalid IGP", func() {
		// Arrange
		invalidIgp := "0x12345"
		_, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:    creator.Address,
			NewOwner: creator.Address,
			IgpId:    invalidIgp,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid hex address length"))
	})

	It("MsgSetIgpOwner (invalid) for non-existing IGP", func() {
		// Arrange
		nonExistingIgp := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"
		_, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:    creator.Address,
			NewOwner: creator.Address,
			IgpId:    nonExistingIgp,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp does not exist: %s", nonExistingIgp)))
	})

	It("MsgSetIgpOwner (invalid) called by non-owner", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:    gasPayer.Address,
			NewOwner: creator.Address,
			IgpId:    igpId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own igp with id %s", gasPayer.Address, igpId)))
	})

	It("MsgSetIgpOwner (valid)", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:    creator.Address,
			NewOwner: gasPayer.Address,
			IgpId:    igpId.String(),
		})

		// Assert
		Expect(err).To(BeNil())

		// Check if the owner has been updated
		igp, err := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(err).To(BeNil())
		Expect(igp.Owner).To(Equal(gasPayer.Address))
	})
})
