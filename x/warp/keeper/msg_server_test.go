package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"

	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	ismTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	pdTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
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

* MsgCreateSyntheticToken (invalid) non-existing Mailbox ID
* MsgCreateSyntheticToken (invalid) when module disabled synthetic tokens
* MsgCreateSyntheticToken (valid)
* MsgCreateCollateralToken (invalid) invalid denom
* MsgCreateCollateralToken (invalid) non-existing Mailbox ID
* MsgCreateCollateralToken (invalid) when module disabled collateral tokens
* MsgCreateCollateralToken (valid)
* MsgEnrollRemoteRouter (invalid) non-existing Token ID
* MsgEnrollRemoteRouter (invalid) non-owner address
* MsgEnrollRemoteRouter (invalid) update with non-owner address
* MsgEnrollRemoteRouter (invalid) invalid remote router
* MsgEnrollRemoteRouter (valid)
* MsgUnrollRemoteRouter (invalid) non-existing Token ID
* MsgUnrollRemoteRouter (invalid) non-owner address
* MsgUnrollRemoteRouter (invalid) non-existing remote domain
* MsgUnrollRemoteRouter (valid)
* MsgSetToken (invalid) non-existing Token ID
* MsgSetToken (invalid) empty new-owner and ISM ID
* MsgSetToken (invalid) non-existing ISM ID
* MsgSetToken (invalid) non-owner address
* MsgSetToken (valid)
* MsgRemoteTransfer (invalid) non-existing Token ID
* MsgRemoteTransfer (invalid) invalid CustomHookMetadata

*/

var denom = "acoin"

var _ = Describe("msg_server.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var owner i.TestValidatorAddress
	var sender i.TestValidatorAddress
	var noopPostDispatchHandler *i.NoopPostDispatchHookHandler

	BeforeEach(func() {
		s = i.NewCleanChain()
		owner = i.GenerateTestValidatorAddress("Owner")
		sender = i.GenerateTestValidatorAddress("Sender")
		err := s.MintBaseCoins(owner.Address, 1_000_000)
		Expect(err).To(BeNil())

		noopPostDispatchHandler = i.CreateNoopDispatchHookHandler(s.App().HyperlaneKeeper.PostDispatchRouter())
		_, err = noopPostDispatchHandler.CreateHook(s.Ctx())
		Expect(err).To(BeNil())
	})

	It("MsgCreateSyntheticToken (invalid) non-existing Mailbox ID", func() {
		// Arrange

		nonExistingMailboxId, _ := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: nonExistingMailboxId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("MsgCreateSyntheticToken (invalid) when module disabled synthetic tokens", func() {
		// Arrange

		// Create new chain with only collateral tokens enabled
		s = i.NewCleanChainWithEnabledTokens([]int32{1})
		err := s.MintBaseCoins(owner.Address, 1_000_000)
		Expect(err).To(BeNil())

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err = s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
		})

		// Assert
		Expect(err.Error()).To(Equal("module disabled synthetic tokens"))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(0))
	})

	It("MsgCreateSyntheticToken (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateCollateralToken (invalid) invalid denom", func() {
		// Arrange
		invalidDenom := "123HYPERLANE!"

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   invalidDenom,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("origin denom %s is invalid", invalidDenom)))
	})

	It("MsgCreateCollateralToken (invalid) non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, _ := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: nonExistingMailboxId,
			OriginDenom:   denom,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("MsgCreateCollateralToken (invalid) when module disabled collateral tokens", func() {
		// Arrange

		// Create new chain with only synthetic tokens enabled
		s = i.NewCleanChainWithEnabledTokens([]int32{2})
		err := s.MintBaseCoins(owner.Address, 1_000_000)
		Expect(err).To(BeNil())

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err = s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})

		// Assert
		Expect(err.Error()).To(Equal("module disabled collateral tokens"))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(0))
	})

	It("MsgCreateCollateralToken (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgEnrollRemoteRouter (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
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
	})

	It("MsgEnrollRemoteRouter (invalid) non-owner address", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        sender.Address,
			TokenId:      tokenId,
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own token with id %s", sender.Address, tokenId.String())))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (invalid) update with non-owner address", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        sender.Address,
			TokenId:      tokenId,
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own token with id %s", sender.Address, tokenId.String())))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.RemoteRouters).To(HaveLen(1))
		Expect(tokens.RemoteRouters[0]).To(Equal(&remoteRouter))
	})

	It("MsgEnrollRemoteRouter (invalid) invalid remote router", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: nil,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid remote router"))

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (valid)", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &remoteRouter,
		})

		// Assert
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.RemoteRouters).To(HaveLen(1))
		Expect(tokens.RemoteRouters[0]).To(Equal(&remoteRouter))
	})

	It("MsgUnrollRemoteRouter (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")

		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		secondRemoteRouter := types.RemoteRouter{
			ReceiverDomain:   2,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def1",
			Gas:              math.NewInt(50000),
		}

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &secondRemoteRouter,
		})
		Expect(err).To(BeNil())

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))

		// Act
		_, err = s.RunTx(&types.MsgUnrollRemoteRouter{
			Owner:          owner.Address,
			TokenId:        nonExistingTokenId,
			ReceiverDomain: secondRemoteRouter.ReceiverDomain,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("token with id %s not found", nonExistingTokenId)))

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))
	})

	It("MsgUnrollRemoteRouter (invalid) non-owner address", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		secondRemoteRouter := types.RemoteRouter{
			ReceiverDomain:   2,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def1",
			Gas:              math.NewInt(50000),
		}

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &secondRemoteRouter,
		})
		Expect(err).To(BeNil())

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))

		// Act
		_, err = s.RunTx(&types.MsgUnrollRemoteRouter{
			Owner:          sender.Address,
			TokenId:        tokenId,
			ReceiverDomain: secondRemoteRouter.ReceiverDomain,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own token with id %s", sender.Address, tokenId.String())))

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))
	})

	It("MsgUnrollRemoteRouter (invalid) non-existing remote domain", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		secondRemoteRouter := types.RemoteRouter{
			ReceiverDomain:   2,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def1",
			Gas:              math.NewInt(50000),
		}

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &secondRemoteRouter,
		})
		Expect(err).To(BeNil())

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))

		// Act
		_, err = s.RunTx(&types.MsgUnrollRemoteRouter{
			Owner:          owner.Address,
			TokenId:        tokenId,
			ReceiverDomain: 3,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find remote router for domain %v", 3)))

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))
	})

	It("MsgUnrollRemoteRouter (valid)", func() {
		// Arrange
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		secondRemoteRouter := types.RemoteRouter{
			ReceiverDomain:   2,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def1",
			Gas:              math.NewInt(50000),
		}

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		err = s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId,
			RemoteRouter: &secondRemoteRouter,
		})
		Expect(err).To(BeNil())

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))

		// Act
		_, err = s.RunTx(&types.MsgUnrollRemoteRouter{
			Owner:          owner.Address,
			TokenId:        tokenId,
			ReceiverDomain: secondRemoteRouter.ReceiverDomain,
		})

		// Assert
		Expect(err).To(BeNil())

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))
	})

	It("MsgSetToken (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId, _ := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")

		// Act
		_, err := s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  nonExistingTokenId,
			IsmId:    nil,
			NewOwner: "new_owner",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find token with id: %s", nonExistingTokenId.String())))
	})

	It("MsgSetToken (invalid) empty new-owner and ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId,
			IsmId:    nil,
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal("new owner or ism id required"))
	})

	It("MsgSetToken (invalid) non-existing ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)
		nonExistingIsm, _ := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId,
			IsmId:    &nonExistingIsm,
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", nonExistingIsm.String())))
	})

	It("MsgSetToken (invalid) non-owner address", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    sender.Address,
			TokenId:  tokenId,
			IsmId:    &secondIsmId,
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own token with id %s", sender.Address, tokenId.String())))
	})

	It("MsgSetToken (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)
		newOwner := "new_owner"

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId,
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId := response.Id

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId,
			IsmId:    &secondIsmId,
			NewOwner: newOwner,
		})

		// Assert
		Expect(err).To(BeNil())

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens).To(HaveLen(1))
		Expect(tokens.Tokens[0].Owner).To(Equal(newOwner))
		Expect(tokens.Tokens[0].IsmId.String()).To(Equal(secondIsmId.String()))
	})

	It("MsgRemoteTransfer (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId, _ := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")

		// Act
		_, err := s.RunTx(&types.MsgRemoteTransfer{
			Sender:             sender.Address,
			TokenId:            nonExistingTokenId,
			DestinationDomain:  0,
			Recipient:          nonExistingTokenId,
			Amount:             math.ZeroInt(),
			CustomHookId:       &nonExistingTokenId,
			GasLimit:           math.ZeroInt(),
			MaxFee:             sdk.NewCoin(denom, math.ZeroInt()),
			CustomHookMetadata: "",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find token with id: %s", nonExistingTokenId)))
	})

	It("MsgRemoteTransfer (invalid) invalid CustomHookMetadata", func() {
		// Arrange
		tokenId, _, _, _ := createToken(s, nil, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_SYNTHETIC)
		invalidCustomHookMetadata := "invalid_custom_hook_metadata"

		// Act
		_, err := s.RunTx(&types.MsgRemoteTransfer{
			Sender:             sender.Address,
			TokenId:            tokenId,
			DestinationDomain:  0,
			Recipient:          tokenId,
			Amount:             math.ZeroInt(),
			CustomHookId:       &tokenId,
			GasLimit:           math.ZeroInt(),
			MaxFee:             sdk.NewCoin(denom, math.ZeroInt()),
			CustomHookMetadata: invalidCustomHookMetadata,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid custom hook metadata"))
	})
})

// Utils
func createIgp(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&pdTypes.MsgCreateIgp{
		Owner: creator,
		Denom: denom,
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createMerkleHook(s *i.KeeperTestSuite, creator string, mailboxId util.HexAddress) util.HexAddress {
	res, err := s.RunTx(&pdTypes.MsgCreateMerkleTreeHook{
		Owner:     creator,
		MailboxId: mailboxId,
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateMerkleTreeHookResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
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

	err := setDestinationGasConfig(s, creator, igpId, destinationDomain)
	Expect(err).To(BeNil())

	res, err := s.RunTx(&coreTypes.MsgCreateMailbox{
		Owner:      creator,
		DefaultIsm: ismId,
	})
	Expect(err).To(BeNil())

	var response coreTypes.MsgCreateMailboxResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId := response.Id

	merkleHook := createMerkleHook(s, creator, mailboxId)

	_, err = s.RunTx(&coreTypes.MsgSetMailbox{
		Owner:        creator,
		MailboxId:    mailboxId,
		DefaultIsm:   &ismId,
		DefaultHook:  &igpId,
		RequiredHook: &merkleHook,
		NewOwner:     creator,
	})
	Expect(err).To(BeNil())

	if err != nil {
		return [32]byte{}, [32]byte{}, [32]byte{}
	}

	return verifyNewMailbox(s, res, creator, igpId.String(), ismId.String(), igpRequired), igpId, ismId
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismTypes.MsgCreateMerkleRootMultisigIsm{
		Creator: creator,
		Validators: []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
		},
		Threshold: 2,
	})
	Expect(err).To(BeNil())

	var response ismTypes.MsgCreateMerkleRootMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismTypes.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response ismTypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId util.HexAddress, domain uint32) error {
	_, err := s.RunTx(&pdTypes.MsgSetDestinationGasConfig{
		Owner: creator,
		IgpId: igpId,
		DestinationGasConfig: &pdTypes.DestinationGasConfig{
			RemoteDomain: domain,
			GasOracle: &pdTypes.GasOracle{
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
	mailboxId := response.Id

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
	Expect(err).To(BeNil())
	Expect(mailbox.Owner).To(Equal(creator))
	Expect(mailbox.DefaultIsm.String()).To(Equal(ismId))
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))
	if igpId != "" {
		Expect(mailbox.DefaultHook.String()).To(Equal(igpId))
	} else {
		Expect(mailbox.DefaultHook).To(BeNil())
	}

	//if igpRequired {
	//	Expect(mailbox.Igp.Required).To(BeTrue()) TODO
	//} else {
	//	Expect(mailbox.Igp.Required).To(BeFalse())
	//}

	mailboxes, err := coreKeeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &coreTypes.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Owner).To(Equal(creator))

	return mailboxId
}

func createToken(s *i.KeeperTestSuite, remoteRouter *types.RemoteRouter, owner, sender string, tokenType types.HypTokenType) (util.HexAddress, util.HexAddress, util.HexAddress, util.HexAddress) {
	mailboxId, igpId, ismId := createValidMailbox(s, owner, "noop", false, 1)

	var tokenId util.HexAddress
	switch tokenType {
	case 1:
		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner,
			OriginDenom:   denom,
			OriginMailbox: mailboxId,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId = response.Id

	case 2:
		res, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner,
			OriginMailbox: mailboxId,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateSyntheticTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId = response.Id
	}

	if remoteRouter != nil {
		_, err := s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner,
			TokenId:      tokenId,
			RemoteRouter: remoteRouter,
		})
		Expect(err).To(BeNil())
	}

	_, err := s.RunTx(&types.MsgSetToken{
		Owner:    owner,
		TokenId:  tokenId,
		IsmId:    &ismId,
		NewOwner: "",
	})
	Expect(err).To(BeNil())

	tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).Tokens(s.Ctx(), &types.QueryTokensRequest{})
	Expect(err).To(BeNil())
	Expect(tokens.Tokens).To(HaveLen(1))
	Expect(tokens.Tokens[0].Owner).To(Equal(owner))

	routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
		Id: tokenId.String(),
	})
	Expect(err).To(BeNil())

	if remoteRouter != nil {
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(remoteRouter))
	} else {
		Expect(routers.RemoteRouters).To(HaveLen(0))
	}
	return tokenId, mailboxId, igpId, ismId
}
