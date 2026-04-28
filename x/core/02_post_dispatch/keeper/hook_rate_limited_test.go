package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("hook_rate_limited_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	var mailboxId util.HexAddress
	var noopHookId util.HexAddress
	var rateLimitedHookId util.HexAddress
	var tokenId util.HexAddress
	var recipient util.HexAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		mailboxId, err = createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())

		noopHookId, err = createDummyNoopHook(s, creator.Address)
		Expect(err).To(BeNil())

		_, err = s.RunTx(&coreTypes.MsgSetMailbox{
			Owner:        creator.Address,
			MailboxId:    mailboxId,
			DefaultHook:  &noopHookId,
			RequiredHook: &noopHookId,
		})
		Expect(err).To(BeNil())

		rateLimitedHookId, err = createDummyRateLimitedHook(s, creator.Address, mailboxId)
		Expect(err).To(BeNil())

		tokenId = util.CreateMockHexAddress("warp-token", 1)
		recipient = util.CreateMockHexAddress("recipient", 1)
	})

	It("Create (valid) Rate Limited Hook", func() {
		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)

		hook, err := qs.RateLimitedHook(s.Ctx(), &types.QueryRateLimitedHookRequest{Id: rateLimitedHookId.String()})
		Expect(err).To(BeNil())
		Expect(hook.RateLimitedHook.Owner).To(Equal(creator.Address))
		Expect(hook.RateLimitedHook.MailboxId).To(Equal(mailboxId))

		hooks, err := qs.RateLimitedHooks(s.Ctx(), &types.QueryRateLimitedHooksRequest{})
		Expect(err).To(BeNil())
		Expect(hooks.RateLimitedHooks).To(HaveLen(1))
		Expect(hooks.RateLimitedHooks[0].Id).To(Equal(rateLimitedHookId))
	})

	It("Create (invalid) Rate Limited Hook (mailbox does not exist)", func() {
		_, err := s.RunTx(&types.MsgCreateRateLimitedHook{
			Owner:     creator.Address,
			MailboxId: util.CreateMockHexAddress("missing-mailbox", 1),
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrMailboxDoesNotExist.Error()))
	})

	It("SetRateLimit creates a full token rate limit and rejects invalid owners/capacities", func() {
		_, err := s.RunTx(&types.MsgSetRateLimit{
			Owner:       creator.Address,
			HookId:      rateLimitedHookId,
			TokenId:     tokenId,
			MaxCapacity: math.NewInt(86_400),
		})
		Expect(err).To(BeNil())

		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)
		tokenRateLimit, err := qs.TokenRateLimit(s.Ctx(), &types.QueryTokenRateLimitRequest{
			HookId:  rateLimitedHookId.String(),
			TokenId: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokenRateLimit.EffectiveTokenRateLimit.TokenRateLimit.FilledLevel).To(Equal(math.NewInt(86_400)))
		Expect(tokenRateLimit.EffectiveTokenRateLimit.CurrentLevel).To(Equal(math.NewInt(86_400)))
		Expect(tokenRateLimit.EffectiveTokenRateLimit.EffectiveCapacity).To(Equal(math.NewInt(86_400)))

		_, err = s.RunTx(&types.MsgSetRateLimit{
			Owner:       i.GenerateTestValidatorAddress("wrong-owner").Address,
			HookId:      rateLimitedHookId,
			TokenId:     tokenId,
			MaxCapacity: math.NewInt(86_400),
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrUnauthorized.Error()))

		_, err = s.RunTx(&types.MsgSetRateLimit{
			Owner:       creator.Address,
			HookId:      rateLimitedHookId,
			TokenId:     tokenId,
			MaxCapacity: math.NewInt(86_399),
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrRateLimitNotSet.Error()))
	})

	It("RemoveRateLimit removes a token rate limit and unconfigured tokens reject", func() {
		setRateLimit(s, creator.Address, rateLimitedHookId, tokenId, math.NewInt(86_400))

		_, err := s.RunTx(&types.MsgRemoveRateLimit{
			Owner:   creator.Address,
			HookId:  rateLimitedHookId,
			TokenId: tokenId,
		})
		Expect(err).To(BeNil())

		message := dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(1))
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrRateLimitNotConfigured.Error()))
	})

	It("RateLimitedHook HookType and Exists", func() {
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(rateLimitedHookId)
		Expect(err).To(BeNil())

		Expect((*handler).HookType()).To(Equal(uint8(types.POST_DISPATCH_HOOK_TYPE_RATE_LIMITED)))

		exists, err := (*handler).Exists(s.Ctx(), rateLimitedHookId)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())

		missingHookId := rateLimitedHookId
		missingHookId[31] = 42
		exists, err = (*handler).Exists(s.Ctx(), missingHookId)
		Expect(err).To(BeNil())
		Expect(exists).To(BeFalse())
	})

	It("QuoteDispatch returns zero coins", func() {
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(rateLimitedHookId)
		Expect(err).To(BeNil())

		quote, err := (*handler).QuoteDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, util.HyperlaneMessage{})
		Expect(err).To(BeNil())
		Expect(quote).To(Equal(sdk.NewCoins()))
	})

	It("PostDispatch consumes a configured token rate limit", func() {
		setRateLimit(s, creator.Address, rateLimitedHookId, tokenId, math.NewInt(86_400))

		message := dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(100))
		charged, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(BeNil())
		Expect(charged).To(Equal(sdk.NewCoins()))

		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)
		tokenRateLimit, err := qs.TokenRateLimit(s.Ctx(), &types.QueryTokenRateLimitRequest{
			HookId:  rateLimitedHookId.String(),
			TokenId: tokenId.String(),
		})
		Expect(err).To(BeNil())
		Expect(tokenRateLimit.EffectiveTokenRateLimit.CurrentLevel).To(Equal(math.NewInt(86_300)))
	})

	It("PostDispatch rejects wrong mailbox, malformed bodies, non-latest messages, and exceeded limits", func() {
		setRateLimit(s, creator.Address, rateLimitedHookId, tokenId, math.NewInt(86_400))

		message := dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(1))
		wrongMailboxId, err := createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), wrongMailboxId, rateLimitedHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrSenderIsNotDesignatedMailbox.Error()))

		malformed := dispatchRawTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, []byte("short"))
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, malformed, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("token message body must be at least 64 bytes"))

		stale := dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(1))
		_ = dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(1))
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, stale, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrInvalidDispatchedMessage.Error()))

		exceeded := dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(86_401))
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, exceeded, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrRateLimitExceeded.Error()))
	})

	It("PostDispatch keeps independent token rate limits per token and refills over time", func() {
		tokenId2 := util.CreateMockHexAddress("warp-token", 2)
		ctx := s.Ctx().WithBlockTime(time.Unix(1_000, 0))

		setRateLimitWithCtx(ctx, s, creator.Address, rateLimitedHookId, tokenId, math.NewInt(864_000))
		setRateLimitWithCtx(ctx, s, creator.Address, rateLimitedHookId, tokenId2, math.NewInt(864_000))

		message := dispatchTestMessage(s, ctx, mailboxId, tokenId, recipient, math.NewInt(1_000))
		_, err := s.App().HyperlaneKeeper.PostDispatch(ctx, mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(BeNil())

		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)
		tokenRateLimit1, err := qs.TokenRateLimit(ctx, &types.QueryTokenRateLimitRequest{HookId: rateLimitedHookId.String(), TokenId: tokenId.String()})
		Expect(err).To(BeNil())
		tokenRateLimit2, err := qs.TokenRateLimit(ctx, &types.QueryTokenRateLimitRequest{HookId: rateLimitedHookId.String(), TokenId: tokenId2.String()})
		Expect(err).To(BeNil())
		Expect(tokenRateLimit1.EffectiveTokenRateLimit.CurrentLevel).To(Equal(math.NewInt(863_000)))
		Expect(tokenRateLimit2.EffectiveTokenRateLimit.CurrentLevel).To(Equal(math.NewInt(864_000)))

		ctx = ctx.WithBlockTime(time.Unix(1_010, 0))
		tokenRateLimit1, err = qs.TokenRateLimit(ctx, &types.QueryTokenRateLimitRequest{HookId: rateLimitedHookId.String(), TokenId: tokenId.String()})
		Expect(err).To(BeNil())
		Expect(tokenRateLimit1.EffectiveTokenRateLimit.CurrentLevel).To(Equal(math.NewInt(863_100)))

		ctx = ctx.WithBlockTime(time.Unix(1_000+int64(types.RateLimitDurationSeconds)+1, 0))
		tokenRateLimit1, err = qs.TokenRateLimit(ctx, &types.QueryTokenRateLimitRequest{HookId: rateLimitedHookId.String(), TokenId: tokenId.String()})
		Expect(err).To(BeNil())
		Expect(tokenRateLimit1.EffectiveTokenRateLimit.CurrentLevel).To(Equal(math.NewInt(864_000)))
	})

	It("Genesis export preserves rate limited hooks and token rate limits", func() {
		setRateLimit(s, creator.Address, rateLimitedHookId, tokenId, math.NewInt(86_400))
		message := dispatchTestMessage(s, s.Ctx(), mailboxId, tokenId, recipient, math.NewInt(1))
		_, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, rateLimitedHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(BeNil())

		genesis := keeper.ExportGenesis(s.Ctx(), s.App().HyperlaneKeeper.PostDispatchKeeper)
		Expect(genesis.RateLimitedHooks).To(HaveLen(1))
		Expect(genesis.RateLimitedHooks[0].Id).To(Equal(rateLimitedHookId))
		Expect(genesis.TokenRateLimits).To(HaveLen(1))
		Expect(genesis.TokenRateLimits[0].HookId).To(Equal(rateLimitedHookId))
		Expect(genesis.TokenRateLimits[0].TokenId).To(Equal(tokenId))
	})
})

func setRateLimit(s *i.KeeperTestSuite, owner string, hookId, tokenId util.HexAddress, maxCapacity math.Int) {
	_, err := s.RunTx(&types.MsgSetRateLimit{
		Owner:       owner,
		HookId:      hookId,
		TokenId:     tokenId,
		MaxCapacity: maxCapacity,
	})
	Expect(err).To(BeNil())
}

func setRateLimitWithCtx(ctx sdk.Context, s *i.KeeperTestSuite, owner string, hookId, tokenId util.HexAddress, maxCapacity math.Int) {
	msgServer := keeper.NewMsgServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)
	_, err := msgServer.SetRateLimit(ctx, &types.MsgSetRateLimit{
		Owner:       owner,
		HookId:      hookId,
		TokenId:     tokenId,
		MaxCapacity: maxCapacity,
	})
	Expect(err).To(BeNil())
}

func dispatchTestMessage(s *i.KeeperTestSuite, ctx sdk.Context, mailboxId, sender, recipient util.HexAddress, amount math.Int) util.HyperlaneMessage {
	return dispatchRawTestMessage(s, ctx, mailboxId, sender, recipient, tokenMessageBody(amount))
}

func dispatchRawTestMessage(s *i.KeeperTestSuite, ctx sdk.Context, mailboxId, sender, recipient util.HexAddress, body []byte) util.HyperlaneMessage {
	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	Expect(err).To(BeNil())

	message := util.HyperlaneMessage{
		Version:     coreTypes.MESSAGE_VERSION,
		Nonce:       mailbox.MessageSent,
		Origin:      mailbox.LocalDomain,
		Sender:      sender,
		Destination: 2,
		Recipient:   recipient,
		Body:        body,
	}

	messageId, err := s.App().HyperlaneKeeper.DispatchMessage(ctx, mailboxId, sender, sdk.NewCoins(), message.Destination, recipient, body, util.StandardHookMetadata{}, nil)
	Expect(err).To(BeNil())
	Expect(messageId).To(Equal(message.Id()))

	return message
}

func tokenMessageBody(amount math.Int) []byte {
	amountBytes := amount.BigInt().Bytes()
	body := make([]byte, 64)
	copy(body[32+(32-len(amountBytes)):64], amountBytes)
	return body
}
