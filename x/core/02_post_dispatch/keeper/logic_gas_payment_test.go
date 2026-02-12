package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_gas_payment.go

* PayForGas (invalid) for non-existing IGP
* PayForGas (invalid) with zero amount
* PayForGas (invalid) with an invalid sender
* PayForGas (invalid) with a non-funded sender
* PayForGas (invalid) insufficient amount
* PayForGas (valid)
* Claim (invalid) for non-existing IGP
* Claim (invalid) from non-owner address
* Claim (invalid) with invalid address
* Claim (invalid) when claimable fees are zero
* Claim (valid)

*/

var _ = Describe("logic_gas_payment.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var gasPayer i.TestValidatorAddress

	denom := "acoin"

	var messageIdTest util.HexAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		gasPayer = i.GenerateTestValidatorAddress("Payer")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		messageIdTest, err = util.DecodeHexAddress("0x000000000000000000000000000000000000006D657373616765496454657374")
		Expect(err).To(BeNil())
	})

	// PayForGas
	It("PayForGas (invalid) for non-existing IGP", func() {
		// Arrange
		nonExistingIgp, err := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")
		Expect(err).To(BeNil())

		res, err := s.RunTx(&types.MsgCreateIgp{
			Owner: creator.Address,
			Denom: denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateIgpResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		igpId := response.Id

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
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             nonExistingIgp,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, math.NewInt(250000)),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp does not exist: %s", nonExistingIgp)))
	})

	It("PayForGas (invalid) with zero amount", func() {
		// NOTE: Negative amount panics at sdk.NewCoins()
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
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, math.ZeroInt()),
		})

		// Assert
		Expect(err.Error()).To(Equal("amount must be greater than zero"))
	})

	It("PayForGas (invalid) insufficient amount", func() {
		// NOTE: Negative amount panics at sdk.NewCoins()
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
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, math.NewInt(5)),
		})

		// Assert
		Expect(err.Error()).To(Equal("required payment exceeds max hyperlane fee: 250000acoin"))
	})

	It("PayForGas (invalid) with an invalid sender", func() {
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
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address + "test",
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, math.NewInt(250000)),
		})

		// Assert
		Expect(err.Error()).To(Equal("decoding bech32 failed: invalid checksum (expected n7qqqp got nltest)"))
	})

	It("PayForGas (invalid) with a non-funded sender", func() {
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
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, math.NewInt(250000)),
		})

		// Assert
		Expect(err.Error()).To(Equal("spendable balance 0acoin is smaller than 250000acoin: insufficient funds"))
	})

	It("PayForGas (valid)", func() {
		// Arrange
		gasAmount := math.NewInt(250000)

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
		igpId := response.Id

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
		Expect(err).To(BeNil())

		gasPayerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom)

		// Act
		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, math.NewInt(250000)),
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom).Amount).To(Equal(gasPayerBalance.Amount.Sub(gasAmount)))

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.AmountOf(denom)).To(Equal(gasAmount))
	})

	// Claim
	It("Claim (invalid) for non-existing IGP", func() {
		// Arrange
		nonExistingIgp, err := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: gasPayer.Address,
			IgpId:  nonExistingIgp,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find igp with id: %s", nonExistingIgp)))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom).Amount).To(Equal(math.ZeroInt()))
	})

	It("Claim (invalid) from non-owner address", func() {
		// Arrange
		gasAmount := math.NewInt(200000)

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
		igpId := response.Id

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
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.ZeroInt(),
			Amount:            sdk.NewCoin(denom, gasAmount),
		})
		Expect(err).To(BeNil())

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.AmountOf(denom)).To(Equal(gasAmount))

		claimableFees := igp.ClaimableFees.AmountOf(denom)
		nonOwnerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom)

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: gasPayer.Address,
			IgpId:  igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to claim: %s is not permitted to claim", gasPayer.Address)))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom).Amount).To(Equal(nonOwnerBalance.Amount))

		igp, _ = s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.AmountOf(denom)).To(Equal(claimableFees))
	})

	It("Claim (invalid) with invalid address", func() {
		// Arrange
		gasAmount := math.NewInt(250000)

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
		igpId := response.Id

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
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, gasAmount),
		})
		Expect(err).To(BeNil())

		ownerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom)

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.AmountOf(denom)).To(Equal(gasAmount))

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: creator.Address + "test",
			IgpId:  igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to claim: %s is not permitted to claim", creator.Address+"test")))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount))
	})

	It("Claim (invalid) when claimable fees are zero", func() {
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

		ownerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom)

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.IsZero()).To(BeTrue())

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: creator.Address,
			IgpId:  igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal("no claimable fees left"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount))

		igp, _ = s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.IsZero()).To(BeTrue())
	})

	It("Claim (valid)", func() {
		// Arrange
		gasAmount := math.NewInt(250000)

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
		igpId := response.Id

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
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgPayForGas{
			Sender:            gasPayer.Address,
			IgpId:             igpId,
			MessageId:         messageIdTest,
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            sdk.NewCoin(denom, gasAmount),
		})
		Expect(err).To(BeNil())

		igp, _ := s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.AmountOf(denom)).To(Equal(gasAmount))

		ownerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom)

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: creator.Address,
			IgpId:  igpId,
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount.Add(igp.ClaimableFees.AmountOf(denom))))

		igp, _ = s.App().HyperlaneKeeper.PostDispatchKeeper.Igps.Get(s.Ctx(), igpId.GetInternalId())
		Expect(igp.ClaimableFees.IsZero()).To(BeTrue())
	})
})
