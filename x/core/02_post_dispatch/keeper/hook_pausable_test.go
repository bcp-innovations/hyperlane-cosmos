package keeper_test

import (
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("hook_pausable_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	var mailboxId util.HexAddress
	var noopHookId util.HexAddress
	var pausableHookId util.HexAddress
	var sender util.HexAddress
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

		pausableHookId, err = createDummyPausableHook(s, creator.Address, mailboxId)
		Expect(err).To(BeNil())

		sender = util.CreateMockHexAddress("sender", 1)
		recipient = util.CreateMockHexAddress("recipient", 1)
	})

	It("Create (valid) Pausable Hook", func() {
		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)

		hook, err := qs.PausableHook(s.Ctx(), &types.QueryPausableHookRequest{Id: pausableHookId.String()})
		Expect(err).To(BeNil())
		Expect(hook.PausableHook.Owner).To(Equal(creator.Address))
		Expect(hook.PausableHook.MailboxId).To(Equal(mailboxId))
		Expect(hook.PausableHook.Paused).To(BeFalse())

		hooks, err := qs.PausableHooks(s.Ctx(), &types.QueryPausableHooksRequest{})
		Expect(err).To(BeNil())
		Expect(hooks.PausableHooks).To(HaveLen(1))
		Expect(hooks.PausableHooks[0].Id).To(Equal(pausableHookId))
	})

	It("Create (invalid) Pausable Hook (mailbox does not exist)", func() {
		_, err := s.RunTx(&types.MsgCreatePausableHook{
			Owner:     creator.Address,
			MailboxId: util.CreateMockHexAddress("missing-mailbox", 1),
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrMailboxDoesNotExist.Error()))
	})

	It("SetPausableHookPaused rejects invalid hook IDs and wrong owners", func() {
		wrongTypeHookId := util.CreateMockHexAddress("wrong-type", int64(types.POST_DISPATCH_HOOK_TYPE_RATE_LIMITED))
		_, err := s.RunTx(&types.MsgSetPausableHookPaused{
			Owner:  creator.Address,
			HookId: wrongTypeHookId,
			Paused: true,
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrInvalidPausableHook.Error()))

		_, err = s.RunTx(&types.MsgSetPausableHookPaused{
			Owner:  i.GenerateTestValidatorAddress("wrong-owner").Address,
			HookId: pausableHookId,
			Paused: true,
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrUnauthorized.Error()))

		missingHookId := pausableHookId
		missingHookId[31] = 42
		_, err = s.RunTx(&types.MsgSetPausableHookPaused{
			Owner:  creator.Address,
			HookId: missingHookId,
			Paused: true,
		})
		Expect(err).To(HaveOccurred())
	})

	It("PausableHook HookType and Exists", func() {
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(pausableHookId)
		Expect(err).To(BeNil())

		Expect((*handler).HookType()).To(Equal(uint8(types.POST_DISPATCH_HOOK_TYPE_PAUSABLE)))

		exists, err := (*handler).Exists(s.Ctx(), pausableHookId)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())

		missingHookId := pausableHookId
		missingHookId[31] = 42
		exists, err = (*handler).Exists(s.Ctx(), missingHookId)
		Expect(err).To(BeNil())
		Expect(exists).To(BeFalse())
	})

	It("QuoteDispatch returns zero coins", func() {
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(pausableHookId)
		Expect(err).To(BeNil())

		quote, err := (*handler).QuoteDispatch(s.Ctx(), mailboxId, pausableHookId, util.StandardHookMetadata{}, util.HyperlaneMessage{})
		Expect(err).To(BeNil())
		Expect(quote).To(Equal(sdk.NewCoins()))
	})

	It("PostDispatch allows unpaused hooks and rejects wrong mailboxes", func() {
		message := dispatchRawTestMessage(s, s.Ctx(), mailboxId, sender, recipient, []byte("message"))

		charged, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, pausableHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(BeNil())
		Expect(charged).To(Equal(sdk.NewCoins()))

		wrongMailboxId, err := createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), wrongMailboxId, pausableHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrSenderIsNotDesignatedMailbox.Error()))
	})

	It("PostDispatch rejects paused hooks and allows dispatch after unpause", func() {
		setPausableHookPaused(s, creator.Address, pausableHookId, true)

		message := dispatchRawTestMessage(s, s.Ctx(), mailboxId, sender, recipient, []byte("message"))
		_, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, pausableHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrPausableHookPaused.Error()))

		setPausableHookPaused(s, creator.Address, pausableHookId, false)

		charged, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, pausableHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(BeNil())
		Expect(charged).To(Equal(sdk.NewCoins()))
	})

	It("AggregationHook PostDispatch returns paused child failure", func() {
		setPausableHookPaused(s, creator.Address, pausableHookId, true)

		aggregationHookId, err := createDummyAggregationHook(s, creator.Address, []util.HexAddress{noopHookId, pausableHookId})
		Expect(err).To(BeNil())

		message := dispatchRawTestMessage(s, s.Ctx(), mailboxId, sender, recipient, []byte("message"))
		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, aggregationHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrPausableHookPaused.Error()))
	})

	It("Genesis export preserves pausable hooks", func() {
		setPausableHookPaused(s, creator.Address, pausableHookId, true)

		genesis := keeper.ExportGenesis(s.Ctx(), s.App().HyperlaneKeeper.PostDispatchKeeper)
		Expect(genesis.PausableHooks).To(HaveLen(1))
		Expect(genesis.PausableHooks[0].Id).To(Equal(pausableHookId))
		Expect(genesis.PausableHooks[0].Paused).To(BeTrue())
	})
})

func setPausableHookPaused(s *i.KeeperTestSuite, owner string, hookId util.HexAddress, paused bool) {
	_, err := s.RunTx(&types.MsgSetPausableHookPaused{
		Owner:  owner,
		HookId: hookId,
		Paused: paused,
	})
	Expect(err).To(BeNil())
}
