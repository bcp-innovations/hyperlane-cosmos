package keeper_test

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"fmt"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_igp.go

* Create (valid) IGP
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

	It("PayForGas (invalid) for non-existing IGP", func() {
		// Arrange
		nonExistingIgp := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

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
			IgpId:             nonExistingIgp,
			MessageId:         "testMessageId",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(1),
			Amount:            math.NewInt(10),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("igp does not exist: %s", nonExistingIgp)))
	})

	It("PayForGas (invalid) with zero amount", func() {
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
			IgpId:             igpId.String(),
			MessageId:         "testMessageId",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(1),
			Amount:            math.NewInt(0),
		})

		// Assert
		Expect(err.Error()).To(Equal("amount must be greater than zero"))
	})

	It("PayForGas (invalid) without message id", func() {
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
			IgpId:             igpId.String(),
			MessageId:         "",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(1),
			Amount:            math.NewInt(10),
		})

		// Assert
		Expect(err.Error()).To(Equal("message id cannot be empty"))
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
			Sender:            gasPayer.Address + "test",
			IgpId:             igpId.String(),
			MessageId:         "messageIdTest",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            math.NewInt(10),
		})

		// Assert
		Expect(err.Error()).To(Equal("decoding bech32 failed: invalid checksum (expected n7qqqp got nltest)"))
	})

	It("PayForGas (invalid) with an non-funded sender", func() {
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
			IgpId:             igpId.String(),
			MessageId:         "messageIdTest",
			DestinationDomain: 1,
			GasLimit:          math.NewInt(50000),
			Amount:            math.NewInt(10),
		})

		// Assert
		Expect(err.Error()).To(Equal("spendable balance 0acoin is smaller than 10acoin: insufficient funds"))
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

	It("Claim (invalid) from non-owner address", func() {
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

		claimableFees := igp.ClaimableFees
		nonOwnerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom)

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: gasPayer.Address,
			IgpId:  igpId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to claim: %s is not permitted to claim", gasPayer.Address)))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), gasPayer.AccAddress, denom).Amount).To(Equal(nonOwnerBalance.Amount))

		igp, _ = s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(claimableFees))
	})

	It("Claim (invalid) with invalid address", func() {
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
			Sender: creator.Address + "test",
			IgpId:  igpId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to claim: %s is not permitted to claim", creator.Address+"test")))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount))
	})

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
		igpId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		ownerBalance := s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom)

		igp, _ := s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(math.ZeroInt()))

		// Act
		_, err = s.RunTx(&types.MsgClaim{
			Sender: creator.Address,
			IgpId:  igpId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal("no claimable fees left"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), creator.AccAddress, denom).Amount).To(Equal(ownerBalance.Amount))

		igp, _ = s.App().HyperlaneKeeper.Igp.Get(s.Ctx(), igpId.Bytes())
		Expect(igp.ClaimableFees).To(Equal(math.ZeroInt()))
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
