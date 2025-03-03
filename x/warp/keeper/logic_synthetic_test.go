package keeper_test

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_synthetic.go

* MsgRemoteTransfer (invalid) non-enrolled router (Synthetic)
* MsgRemoteTransfer (invalid) empty cosmos sender (Synthetic)
* MsgRemoteTransfer (invalid) invalid cosmos sender (Synthetic)
* MsgRemoteTransfer (invalid) empty recipient (Synthetic)
* MsgRemoteTransfer (invalid) invalid recipient (Synthetic)
* MsgRemoteTransfer (invalid) no enrolled router (Synthetic)
* MsgRemoteTransfer (invalid) receiver contract (Synthetic)
* MsgRemoteTransfer (invalid) insufficient funds (Synthetic)
* MsgRemoteTransfer (valid) (Synthetic)
* MsgRemoteTransfer && MsgRemoteReceiveSynthetic (valid) (Synthetic)

*/

var _ = Describe("logic_synthetic.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var owner i.TestValidatorAddress
	var sender i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		owner = i.GenerateTestValidatorAddress("Owner")
		sender = i.GenerateTestValidatorAddress("Sender")
		err := s.MintBaseCoins(owner.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("MsgRemoteTransfer (invalid) non-enrolled router (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, nil, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId,
			DestinationDomain: 1,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal("no enrolled router found for destination domain 1"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) empty cosmos sender (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            "",
			TokenId:           tokenId,
			DestinationDomain: remoteRouter.ReceiverDomain,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal("empty address string is not allowed"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) invalid cosmos sender (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            "Test123!",
			TokenId:           tokenId,
			DestinationDomain: remoteRouter.ReceiverDomain,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal("decoding bech32 failed: string not all lowercase or all uppercase"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) no enrolled router (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId,
			DestinationDomain: 2,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("no enrolled router found for destination domain %d", 2)))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) receiver contract (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865de",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId,
			DestinationDomain: remoteRouter.ReceiverDomain,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to decode receiver contract address %s", remoteRouter.ReceiverContract)))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) insufficient funds (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		// Act
		_, err := s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId,
			DestinationDomain: 1,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("spendable balance 0%s is smaller than 100%s: insufficient funds", syntheticDenom, syntheticDenom)))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer && MsgRemoteReceiveSynthetic (valid) (Synthetic)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, mailboxId, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)

		syntheticDenom := "hyperlane/" + tokenId.String()

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId,
			DestinationDomain: remoteRouter.ReceiverDomain,
			Recipient:         receiverAddress,
			Amount:            amount,
			CustomHookId:      &igpId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            maxFee,
		})
		Expect(err).To(BeNil())

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount.Int64()).To(Equal(senderBalance.Amount.Sub(amount).Int64()))

		receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
		Expect(err).To(BeNil())

		warpRecipient, err := sdk.GetFromBech32(sender.Address, "hyp")
		Expect(err).To(BeNil())

		warpPayload, err := types.NewWarpPayload(warpRecipient, *big.NewInt(amount.Int64()))
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       1,
			Origin:      remoteRouter.ReceiverDomain,
			Sender:      receiverContract,
			Destination: 0,
			Recipient:   tokenId,
			Body:        warpPayload.Bytes(),
		}

		senderBalance = s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		_, err = s.RunTx(&coreTypes.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount.Add(amount)))
	})
})
