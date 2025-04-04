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


* SetDestinationGasConfig (invalid) for non-existing IGP
* SetDestinationGasConfig (invalid) with wrong owner
* SetDestinationGasConfig (invalid) without gas oracle
* CreateIgp (invalid) with invalid denom
* CreateIgp (valid)
* SetDestinationGasConfig (invalid) for non-existing IGP
* SetDestinationGasConfig (valid)
* MsgSetIgpOwner (invalid) for non-existing IGP
* MsgSetIgpOwner (invalid) called by non-owner
* MsgSetToken (invalid) invalid new owner
* MsgSetToken (invalid) renounce ownership with new owner set
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
		nonExistingIgp, err := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgCreateIgp{
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
		igpId := response.Id

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.Owner).To(Equal(creator.Address))
		Expect(igp.Denom).To(Equal(denom))
		Expect(igp.ClaimableFees.IsZero()).To(BeTrue())
	})

	// SetDestinationGasConfig
	It("SetDestinationGasConfig (invalid) for non-existing IGP", func() {
		// Arrange
		invalidIgp, err := util.DecodeHexAddress("0x10df2f89cb24ed6078fc3949b4870e94a7e32e40e8d8c6b7bd74ccc2c933d760")
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgCreateIgp{
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
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp does not exist: %s", invalidIgp.String())))
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
		igpId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: igpId,
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

	It("MsgSetIgpOwner (invalid) for non-existing IGP", func() {
		// Arrange
		nonExistingIgp, err := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")
		Expect(err).To(BeNil())
		_, err = s.RunTx(&types.MsgCreateIgp{
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
		igpId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:    gasPayer.Address,
			NewOwner: creator.Address,
			IgpId:    igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own igp with id %s", gasPayer.Address, igpId)))
	})

	It("MsgSetToken (invalid) invalid new owner", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:             creator.Address,
			NewOwner:          "new_owner",
			IgpId:             igpId,
			RenounceOwnership: false,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid new owner"))
	})

	It("MsgSetToken (invalid) renounce ownership with new owner set", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:             creator.Address,
			NewOwner:          gasPayer.Address,
			IgpId:             igpId,
			RenounceOwnership: true,
		})

		// Assert
		Expect(err.Error()).To(Equal("cannot set new owner and renounce ownership at the same time"))
	})

	It("MsgSetIgpOwner (valid) - renounce ownership", func() {
		// Arrange
		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:             creator.Address,
			NewOwner:          "",
			IgpId:             igpId,
			RenounceOwnership: true,
		})

		// Assert
		Expect(err).To(BeNil())

		// Check if the owner has been updated
		igp, err := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(err).To(BeNil())
		Expect(igp.Owner).To(Equal(""))
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
		igpId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetIgpOwner{
			Owner:    creator.Address,
			NewOwner: gasPayer.Address,
			IgpId:    igpId,
		})

		// Assert
		Expect(err).To(BeNil())

		// Check if the owner has been updated
		igp, err := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(err).To(BeNil())
		Expect(igp.Owner).To(Equal(gasPayer.Address))
	})
})
