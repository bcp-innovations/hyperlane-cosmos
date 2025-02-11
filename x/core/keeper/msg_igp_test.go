package keeper_test

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_igp.go

* CreateIgp (invalid) with invalid denom
* CreateIgp (valid)
* SetDestinationGasConfig (invalid) for invalid IGP
* SetDestinationGasConfig (valid)
* PayForGas (invalid) for invalid IGP
* PayForGas (valid)
* Claim (invalid) for invalid ISM
* Claim (valid)

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

		igp, _ := s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.Owner).To(Equal(creator.Address))
		Expect(igp.Denom).To(Equal(denom))
		Expect(igp.ClaimableFees).To(Equal(math.NewInt(0)))
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
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: invalid hex address length", invalidIgp)))
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

		rng := collections.NewPrefixedPairRange[[]byte, uint32](igpId.Bytes())

		iter, err := s.App().HyperlaneKeeper.IgpDestinationGasConfigMap.Iterate(s.Ctx(), rng)
		Expect(err).To(BeNil())

		destinationGasConfigs, err := iter.Values()
		Expect(err).To(BeNil())

		Expect(destinationGasConfigs).To(HaveLen(1))
		Expect(destinationGasConfigs[0].RemoteDomain).To(Equal(uint32(1)))
		Expect(destinationGasConfigs[0].GasOracle.TokenExchangeRate).To(Equal(math.NewInt(1e10)))
		Expect(destinationGasConfigs[0].GasOracle.GasPrice).To(Equal(math.NewInt(1)))
		Expect(destinationGasConfigs[0].GasOverhead).To(Equal(math.NewInt(200000)))
	})

	// PayForGas
	It("PayForGas (invalid) for invalid IGP", func() {
		// Arrange
		invalidIgp := "0x12345"

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
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             invalidIgp,
			MessageId:         "testMessageId",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(1),
			Amount:            math.NewInt(10),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: invalid hex address length", invalidIgp)))
	})

	It("PayForGas (valid)", func() {
		// Arrange
		gasAmount := math.NewInt(10)

		err := s.MintBaseCoins(gasPayer.Address, 1_000_000)
		Expect(err).To(BeNil())

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
		Expect(err).To(BeNil())

		gasPayerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom)

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId.String(),
			MessageId:         "messageIdTest",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            gasAmount,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom).Amount).To(Equal(gasPayerBalance.Amount.Sub(gasAmount)))

		igp, _ := s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(gasAmount))
	})

	// Claim
	It("Claim (invalid) for invalid ISM", func() {
		// Arrange
		gasAmount := math.NewInt(10)

		err := s.MintBaseCoins(gasPayer.Address, 1_000_000)
		Expect(err).To(BeNil())

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
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId.String(),
			MessageId:         "messageIdTest",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            gasAmount,
		})
		Expect(err).To(BeNil())

		ownerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom)

		igp, _ := s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(gasAmount))

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: creator.Address,
			IgpId:  igpId.String() + "test",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism id %s is invalid: %s", igpId.String()+"test", "invalid hex address length")))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount))
	})

	It("Claim (valid)", func() {
		// Arrange
		gasAmount := math.NewInt(10)

		err := s.MintBaseCoins(gasPayer.Address, 1_000_000)
		Expect(err).To(BeNil())

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
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId.String(),
			MessageId:         "messageIdTest",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            gasAmount,
		})
		Expect(err).To(BeNil())

		igp, _ := s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(gasAmount))

		ownerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom)

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: creator.Address,
			IgpId:  igpId.String(),
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount.Add(igp.ClaimableFees)))

		igp, _ = s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(math.ZeroInt()))
	})
})
