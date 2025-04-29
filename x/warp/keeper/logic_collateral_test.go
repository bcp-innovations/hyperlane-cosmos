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

TEST CASES - logic_collateral.go

* MsgRemoteTransfer (invalid) non-enrolled router (Collateral)
* MsgRemoteTransfer (invalid) empty cosmos sender (Collateral)
* MsgRemoteTransfer (invalid) invalid cosmos sender (Collateral)
* MsgRemoteTransfer (invalid) empty recipient (Collateral)
* MsgRemoteTransfer (invalid) invalid recipient (Collateral)
* MsgRemoteTransfer (invalid) no enrolled router for destination (Collateral)
* MsgRemoteTransfer (invalid) receiver contract (Collateral)
* MsgRemoteTransfer (invalid) insufficient funds (Collateral)
* MsgRemoteTransfer & MsgRemoteReceiveCollateral (invalid) not enough collateral (Collateral)
* MsgRemoteTransfer && MsgRemoteReceiveCollateral (valid) (Collateral)

*/

var _ = Describe("logic_collateral.go", Ordered, func() {
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

	It("MsgRemoteTransfer (invalid) non-enrolled router (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, nil, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

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
		Expect(err.Error()).To(Equal("no enrolled router found for destination domain 2"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) empty cosmos sender (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, nil, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

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
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) invalid cosmos sender (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, nil, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

		err := s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Amount.Int64()).Uint64())
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

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
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) no enrolled router for destination (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

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
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) receiver contract (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865de",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

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
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer (invalid) insufficient funds (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, _, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

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
		Expect(err.Error()).To(Equal("spendable balance 0acoin is smaller than 100acoin: insufficient funds"))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer & MsgRemoteReceiveCollateral (invalid) not enough collateral (Collateral)", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)

		tokenId, mailboxId, _, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		// Act
		receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
		Expect(err).To(BeNil())

		warpRecipient, err := sdk.GetFromBech32(sender.Address, "hyp")
		Expect(err).To(BeNil())

		warpPayload, err := types.NewWarpPayload(warpRecipient, *big.NewInt(amount.Int64()), nil)
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     3,
			Nonce:       1,
			Origin:      remoteRouter.ReceiverDomain,
			Sender:      receiverContract,
			Destination: 0,
			Recipient:   tokenId,
			Body:        warpPayload.Bytes(),
		}

		_, err = s.RunTx(&coreTypes.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(types.ErrNotEnoughCollateral.Error()))
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount))
	})

	It("MsgRemoteTransfer && MsgRemoteReceiveCollateral (valid) (Collateral)", func() {
		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, mailboxId, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL)
		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)
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

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Sub(amount.Add(maxFee.Amount))))

		receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
		Expect(err).To(BeNil())

		warpRecipient, err := sdk.GetFromBech32(sender.Address, "hyp")
		Expect(err).To(BeNil())

		warpPayload, err := types.NewWarpPayload(warpRecipient, *big.NewInt(amount.Int64()), nil)
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     3,
			Nonce:       1,
			Origin:      remoteRouter.ReceiverDomain,
			Sender:      receiverContract,
			Destination: 0,
			Recipient:   tokenId,
			Body:        warpPayload.Bytes(),
		}

		senderBalance = s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		_, err = s.RunTx(&coreTypes.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Add(amount)))
	})
})
