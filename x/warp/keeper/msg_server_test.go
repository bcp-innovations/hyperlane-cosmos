package keeper_test

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	coreKeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - msg_server.go

*

*/

var denom = "acoin"

var _ = Describe("msg_server.go", Ordered, func() {
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

	// DispatchMessage
	It("MsgRemoteTransfer (invalid) non-enrolled router", func() {
		// Arrange
		receiverAddress := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
		}

		amount := math.NewInt(100)
		maxFee := math.NewInt(250000)

		mailboxId, igpId, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		err = s.App().HyperlaneKeeper.RegisterReceiverIsm(s.Ctx(), tokenId, mailboxId, ismId.String())
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(1))
		Expect(tokens.Tokens[0].RemoteRouters[0]).To(Equal(&remoteRouter))

		// Act
		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId.String(),
			DestinationDomain: 2,
			Recipient:         receiverAddress,
			Amount:            amount,
			IgpId:             igpId.String(),
			GasLimit:          math.NewInt(50000),
			MaxFee:            maxFee,
		})

		// Assert
		Expect(err.Error()).To(Equal("no enrolled router found for destination domain 2"))
	})

	It("MsgRemoteTransfer (valid) (Synthetic)", func() {
		// Arrange
		receiverAddress := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
		}

		amount := math.NewInt(100)
		maxFee := math.NewInt(250000)

		mailboxId, igpId, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateSyntheticTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		syntheticDenom := "hyperlane/" + tokenId.String()

		err = s.MintBaseCoins(sender.Address, math.NewInt(maxFee.Int64()).Uint64())
		Expect(err).To(BeNil())

		err = s.MintCoins(sender.Address, sdk.NewCoins(sdk.NewInt64Coin(syntheticDenom, amount.Int64())))
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		err = s.App().HyperlaneKeeper.RegisterReceiverIsm(s.Ctx(), tokenId, mailboxId, ismId.String())
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(1))
		Expect(tokens.Tokens[0].RemoteRouters[0]).To(Equal(&remoteRouter))

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom)

		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId.String(),
			DestinationDomain: remoteRouter.ReceiverDomain,
			Recipient:         receiverAddress,
			Amount:            amount,
			IgpId:             igpId.String(),
			GasLimit:          math.NewInt(50000),
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

		message := coreTypes.HyperlaneMessage{
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
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})
		Expect(err).To(BeNil())

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, syntheticDenom).Amount).To(Equal(senderBalance.Amount.Add(amount)))

		// Act

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateSyntheticToken (invalid) invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		invalidMailboxId := mailboxId.String() + "test"

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: invalidMailboxId,
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: invalid hex address length"))
	})

	It("MsgCreateSyntheticToken (invalid) non-existing Mailbox ID", func() {
		// Arrange
		_, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		nonExistingMailboxId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: nonExistingMailboxId,
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("MsgCreateSyntheticToken (invalid) non-existing ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		nonExistingIsmId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			IsmId:         nonExistingIsmId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", nonExistingIsmId)))
	})

	It("MsgCreateSyntheticToken (valid) with default ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			IsmId:         "",
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateSyntheticToken (valid)", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateCollateralToken (invalid) invalid denom", func() {
		// Arrange
		invalidDenom := "123HYPERLANE!"

		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   invalidDenom,
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("origin denom %s is invalid", invalidDenom)))
	})

	It("MsgCreateCollateralToken (invalid) invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		invalidMailboxId := mailboxId.String() + "test"

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: invalidMailboxId,
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: invalid hex address length"))
	})

	It("MsgCreateCollateralToken (invalid) non-existing Mailbox ID", func() {
		// Arrange
		_, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		nonExistingMailboxId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: nonExistingMailboxId,
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("MsgCreateCollateralToken (invalid) non-existing ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		nonExistingIsmId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         nonExistingIsmId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", nonExistingIsmId)))
	})

	It("MsgCreateCollateralToken (valid) with default ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         "",
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateCollateralToken (valid)", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgRemoteTransfer (valid) (Collateral)", func() {
		// Arrange
		receiverAddress := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
		}

		amount := math.NewInt(100)
		maxFee := math.NewInt(250000)

		mailboxId, igpId, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		err = s.App().HyperlaneKeeper.RegisterReceiverIsm(s.Ctx(), tokenId, mailboxId, ismId.String())
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(1))
		Expect(tokens.Tokens[0].RemoteRouters[0]).To(Equal(&remoteRouter))

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		_, err = s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           tokenId.String(),
			DestinationDomain: remoteRouter.ReceiverDomain,
			Recipient:         receiverAddress,
			Amount:            amount,
			IgpId:             igpId.String(),
			GasLimit:          math.NewInt(50000),
			MaxFee:            maxFee,
		})
		Expect(err).To(BeNil())

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Sub(amount.Add(maxFee))))

		receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
		Expect(err).To(BeNil())

		warpRecipient, err := sdk.GetFromBech32(sender.Address, "hyp")
		Expect(err).To(BeNil())

		warpPayload, err := types.NewWarpPayload(warpRecipient, *big.NewInt(amount.Int64()))
		Expect(err).To(BeNil())

		message := coreTypes.HyperlaneMessage{
			Version:     1,
			Nonce:       1,
			Origin:      remoteRouter.ReceiverDomain,
			Sender:      receiverContract,
			Destination: 0,
			Recipient:   tokenId,
			Body:        warpPayload.Bytes(),
		}

		senderBalance = s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		_, err = s.RunTx(&coreTypes.MsgProcessMessage{
			MailboxId: mailboxId.String(),
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})
		Expect(err).To(BeNil())

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Add(amount)))

		// Act

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgEnrollRemoteRouter (invalid) invalid Token ID (Collateral)", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String() + "test",
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("invalid token id %s", tokenId.String()+"test")))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (invalid) non-existing Token ID (Collateral)", func() {
		// Arrange
		nonExistingTokenId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      nonExistingTokenId,
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("token with id %s not found", nonExistingTokenId)))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (invalid) non-owner address (Collateral)", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        sender.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own token with id %s", sender.Address, tokenId.String())))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (invalid) invalid remote router (Collateral)", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid remote router"))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (valid) (Collateral)", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
		}

		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &remoteRouter,
		})

		// Assert
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(1))
		Expect(tokens.Tokens[0].RemoteRouters[0]).To(Equal(&remoteRouter))
	})

	It("MsgUnrollRemoteRouter (valid)", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
		}

		secondRemoteRouter := types.RemoteRouter{
			ReceiverDomain:   2,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def1",
		}

		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(1))
		Expect(tokens.Tokens[0].RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &secondRemoteRouter,
		})
		Expect(err).To(BeNil())

		tokens, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(2))
		Expect(tokens.Tokens[0].RemoteRouters[1]).To(Equal(&secondRemoteRouter))

		// Act
		_, err = s.RunTx(&types.MsgUnrollRemoteRouter{
			Owner:          owner.Address,
			TokenId:        tokenId.String(),
			ReceiverDomain: secondRemoteRouter.ReceiverDomain,
		})

		// Assert
		Expect(err).To(BeNil())
		tokens, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].RemoteRouters).To(HaveLen(1))
		Expect(tokens.Tokens[0].RemoteRouters[0]).To(Equal(&remoteRouter))
	})

	It("MsgSetInterchainSecurityModule (valid)", func() {
		// Arrange
		mailboxId, _, ismId := createValidMailbox(s, owner.Address, "noop", false, 1)

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
			IsmId:         ismId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetInterchainSecurityModule{
			Owner:   owner.Address,
			TokenId: tokenId.String(),
			IsmId:   secondIsmId.String(),
		})

		// Assert
		Expect(err).To(BeNil())
	})
})

// Utils
func createIgp(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&coreTypes.MsgCreateIgp{
		Owner: creator,
		Denom: denom,
	})
	Expect(err).To(BeNil())

	var response coreTypes.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	igpId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return igpId
}

func createValidMailbox(s *i.KeeperTestSuite, creator string, ism string, igpRequired bool, destinationDomain uint32) (util.HexAddress, util.HexAddress, util.HexAddress) {
	var ismId util.HexAddress
	switch ism {
	case "noop":
		ismId = createNoopIsm(s, creator)
	case "multisig":
		ismId = createMultisigIsm(s, creator)
	}

	igpId := createIgp(s, creator)

	err := setDestinationGasConfig(s, creator, igpId.String(), destinationDomain)
	Expect(err).To(BeNil())

	res, err := s.RunTx(&coreTypes.MsgCreateMailbox{
		Creator:    creator,
		DefaultIsm: ismId.String(),
		Igp: &coreTypes.InterchainGasPaymaster{
			Id:       igpId.String(),
			Required: igpRequired,
		},
	})
	Expect(err).To(BeNil())

	return verifyNewMailbox(s, res, creator, igpId.String(), ismId.String(), igpRequired), igpId, ismId
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&coreTypes.MsgCreateMultisigIsm{
		Creator: creator,
		MultiSig: &coreTypes.MultiSigIsm{
			Validators: []string{
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		},
	})
	Expect(err).To(BeNil())

	var response coreTypes.MsgCreateMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	ismId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return ismId
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&coreTypes.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response coreTypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	ismId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return ismId
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId string, domain uint32) error {
	_, err := s.RunTx(&coreTypes.MsgSetDestinationGasConfig{
		Owner: creator,
		IgpId: igpId,
		DestinationGasConfig: &coreTypes.DestinationGasConfig{
			RemoteDomain: domain,
			GasOracle: &coreTypes.GasOracle{
				TokenExchangeRate: math.NewInt(1e10),
				GasPrice:          math.NewInt(1),
			},
			GasOverhead: math.NewInt(200000),
		},
	})

	return err
}

func verifyNewMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, igpId, ismId string, igpRequired bool) util.HexAddress {
	var response coreTypes.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.Bytes())
	Expect(err).To(BeNil())
	Expect(mailbox.Creator).To(Equal(creator))
	Expect(mailbox.Igp.Id).To(Equal(igpId))
	Expect(mailbox.DefaultIsm).To(Equal(ismId))
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

	if igpRequired {
		Expect(mailbox.Igp.Required).To(BeTrue())
	} else {
		Expect(mailbox.Igp.Required).To(BeFalse())
	}

	mailboxes, err := coreKeeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &coreTypes.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Creator).To(Equal(creator))

	return mailboxId
}
