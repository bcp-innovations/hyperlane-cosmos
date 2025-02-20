package keeper_test

import (
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

		message := util.HyperlaneMessage{
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
