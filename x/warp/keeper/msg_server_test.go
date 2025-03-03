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

* MsgCreateSyntheticToken (invalid) invalid Mailbox ID
* MsgCreateSyntheticToken (invalid) non-existing Mailbox ID
* MsgCreateSyntheticToken (invalid) non-existing ISM ID
* MsgCreateSyntheticToken (valid) with default ISM ID
* MsgCreateSyntheticToken (valid)
* MsgCreateCollateralToken (invalid) invalid denom
* MsgCreateCollateralToken (invalid) invalid Mailbox ID
* MsgCreateCollateralToken (invalid) non-existing Mailbox ID
* MsgCreateCollateralToken (invalid) non-existing ISM ID
* MsgCreateCollateralToken (valid) with default ISM ID
* MsgCreateCollateralToken (valid)
* MsgEnrollRemoteRouter (invalid) invalid Token ID
* MsgEnrollRemoteRouter (invalid) non-existing Token ID
* MsgEnrollRemoteRouter (invalid) non-owner address
* MsgEnrollRemoteRouter (invalid) invalid remote router
* MsgEnrollRemoteRouter (valid)
* MsgUnrollRemoteRouter (invalid) invalid Token ID
* MsgUnrollRemoteRouter (invalid) non-existing Token ID
* MsgUnrollRemoteRouter (invalid) non-owner address
* MsgUnrollRemoteRouter (invalid) non-existing remote domain
* MsgUnrollRemoteRouter (valid)
* MsgSetInterchainSecurityModule (invalid) empty ISM ID
* MsgSetInterchainSecurityModule (invalid) invalid Token ID
* MsgSetInterchainSecurityModule (invalid) non-owner address
* MsgSetInterchainSecurityModule (invalid) invalid ISM ID
* MsgSetInterchainSecurityModule (valid)
* MsgRemoteTransfer (invalid) invalid Token ID
* MsgRemoteTransfer (invalid) non-existing Token ID

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

	It("MsgCreateSyntheticToken (invalid) invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		invalidMailboxId := mailboxId.String() + "test"

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: invalidMailboxId,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: invalid hex address length"))
	})

	It("MsgCreateSyntheticToken (invalid) non-existing Mailbox ID", func() {
		// Arrange

		nonExistingMailboxId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: nonExistingMailboxId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	// TODO should it be allowed to set invalid ISM ids?
	PIt("MsgCreateSyntheticToken (invalid) non-existing ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		nonExistingIsmId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
		})
		Expect(err).To(BeNil())

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
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateSyntheticToken (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   invalidDenom,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("origin denom %s is invalid", invalidDenom)))
	})

	It("MsgCreateCollateralToken (invalid) invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		invalidMailboxId := mailboxId.String() + "test"

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: invalidMailboxId,
			OriginDenom:   denom,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: invalid hex address length"))
	})

	It("MsgCreateCollateralToken (invalid) non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: nonExistingMailboxId,
			OriginDenom:   denom,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	// TODO
	PIt("MsgCreateCollateralToken (invalid) non-existing ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		nonExistingIsmId := "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0"

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgCreateCollateralToken (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		// Act
		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgEnrollRemoteRouter (invalid) invalid Token ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.RemoteRouters).To(HaveLen(0))
	})

	It("MsgEnrollRemoteRouter (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		_, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
			RemoteRouter: &remoteRouter,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        sender.Address,
			TokenId:      tokenId.String(),
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		tokens, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.RemoteRouters).To(HaveLen(1))
		Expect(tokens.RemoteRouters[0]).To(Equal(&remoteRouter))
	})

	It("MsgUnrollRemoteRouter (invalid) invalid Token ID", func() {
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokens.Tokens[0].Owner).To(Equal(owner.Address))
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
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
			TokenId:        tokenId.String() + "test",
			ReceiverDomain: secondRemoteRouter.ReceiverDomain,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("invalid token id %s", tokenId.String()+"test")))

		routers, err = keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(2))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))
		Expect(routers.RemoteRouters[1]).To(Equal(&secondRemoteRouter))
	})

	It("MsgUnrollRemoteRouter (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId := "0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647"

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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())

		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
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
			TokenId:        tokenId.String(),
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
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
			TokenId:        tokenId.String(),
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
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
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

		routers, err := keeper.NewQueryServerImpl(s.App().WarpKeeper).RemoteRouters(s.Ctx(), &types.QueryRemoteRoutersRequest{
			Id: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(routers.RemoteRouters).To(HaveLen(1))
		Expect(routers.RemoteRouters[0]).To(Equal(&remoteRouter))

		_, err = s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner.Address,
			TokenId:      tokenId.String(),
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
			TokenId:        tokenId.String(),
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

	It("MsgSetInterchainSecurityModule (invalid) empty ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId.String(),
			IsmId:    "",
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal("new owner or ism id required"))
	})

	It("MsgSetInterchainSecurityModule (invalid) invalid Token ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId.String() + "test",
			IsmId:    secondIsmId.String(),
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("invalid token id %s", tokenId.String()+"test")))
	})

	It("MsgSetInterchainSecurityModule (invalid) non-owner address", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    sender.Address,
			TokenId:  tokenId.String(),
			IsmId:    secondIsmId.String(),
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own token with id %s", sender.Address, tokenId.String())))
	})

	It("MsgSetInterchainSecurityModule (invalid) invalid ISM ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId.String(),
			IsmId:    secondIsmId.String() + "test",
			NewOwner: "",
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid ism id: invalid hex address length"))
	})

	It("MsgSetInterchainSecurityModule (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, owner.Address, "noop", false, 1)

		secondIsmId := createNoopIsm(s, owner.Address)

		res, err := s.RunTx(&types.MsgCreateCollateralToken{
			Owner:         owner.Address,
			OriginMailbox: mailboxId.String(),
			OriginDenom:   denom,
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err := util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgSetToken{
			Owner:    owner.Address,
			TokenId:  tokenId.String(),
			IsmId:    secondIsmId.String(),
			NewOwner: "",
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("MsgRemoteTransfer (invalid) non-existing Token ID", func() {
		// Arrange
		nonExistingTokenId, _ := util.DecodeHexAddress("0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0")

		// Act
		_, err := s.RunTx(&types.MsgRemoteTransfer{
			Sender:            sender.Address,
			TokenId:           nonExistingTokenId,
			DestinationDomain: 0,
			Recipient:         nonExistingTokenId,
			Amount:            math.ZeroInt(),
			IgpId:             &nonExistingTokenId,
			GasLimit:          math.ZeroInt(),
			MaxFee:            sdk.NewCoin(denom, math.ZeroInt()),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find token with id: %s", nonExistingTokenId)))
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

	igpId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return igpId
}

func createMerkleHook(s *i.KeeperTestSuite, creator string, mailboxId string) util.HexAddress {
	res, err := s.RunTx(&pdTypes.MsgCreateMerkleTreeHook{
		Owner:     creator,
		MailboxId: mailboxId,
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateMerkleTreeHookResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	hookId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return hookId
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
		Owner:      creator,
		DefaultIsm: ismId,
	})
	Expect(err).To(BeNil())

	var response coreTypes.MsgCreateMailboxResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	merkleHook := createMerkleHook(s, creator, mailboxId.String())

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

	ismId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return ismId
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismTypes.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response ismTypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	ismId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return ismId
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId string, domain uint32) error {
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
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

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
			OriginMailbox: mailboxId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateCollateralTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err = util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())

	case 2:
		res, err := s.RunTx(&types.MsgCreateSyntheticToken{
			Owner:         owner,
			OriginMailbox: mailboxId.String(),
		})
		Expect(err).To(BeNil())

		var response types.MsgCreateSyntheticTokenResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())
		tokenId, err = util.DecodeHexAddress(response.Id)
		Expect(err).To(BeNil())
	}

	if remoteRouter != nil {
		_, err := s.RunTx(&types.MsgEnrollRemoteRouter{
			Owner:        owner,
			TokenId:      tokenId.String(),
			RemoteRouter: remoteRouter,
		})
		Expect(err).To(BeNil())
	}

	_, err := s.RunTx(&types.MsgSetToken{
		Owner:    owner,
		TokenId:  tokenId.String(),
		IsmId:    ismId.String(),
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
