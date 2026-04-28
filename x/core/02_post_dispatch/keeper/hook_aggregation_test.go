package keeper_test

import (
	"cosmossdk.io/math"
	i "github.com/bcp-innovations/hyperlane-cosmos/tests/integration"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("hook_aggregation_test.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	var mailboxId util.HexAddress
	var noopHookId util.HexAddress
	var merkleHookId util.HexAddress
	var aggregationHookId util.HexAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")

		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())

		mailboxId, err = createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())

		noopHookId, err = createDummyNoopHook(s, creator.Address)
		Expect(err).To(BeNil())

		merkleHookId, err = createDummyMerkleTreeHook(s, creator.Address, mailboxId)
		Expect(err).To(BeNil())

		aggregationHookId, err = createDummyAggregationHook(s, creator.Address, []util.HexAddress{noopHookId, merkleHookId})
		Expect(err).To(BeNil())
	})

	It("Create (valid) Aggregation Hook", func() {
		qs := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.PostDispatchKeeper)

		hook, err := qs.AggregationHook(s.Ctx(), &types.QueryAggregationHookRequest{Id: aggregationHookId.String()})
		Expect(err).To(BeNil())
		Expect(hook.AggregationHook.Owner).To(Equal(creator.Address))
		Expect(hook.AggregationHook.Hooks).To(Equal([]util.HexAddress{noopHookId, merkleHookId}))

		hooks, err := qs.AggregationHooks(s.Ctx(), &types.QueryAggregationHooksRequest{})
		Expect(err).To(BeNil())
		Expect(hooks.AggregationHooks).To(HaveLen(1))
		Expect(hooks.AggregationHooks[0].Id).To(Equal(aggregationHookId))
	})

	It("Create (invalid) Aggregation Hook (empty hooks)", func() {
		_, err := s.RunTx(&types.MsgCreateAggregationHook{
			Owner: creator.Address,
			Hooks: []util.HexAddress{},
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("hooks cannot be empty"))
	})

	It("Create (invalid) Aggregation Hook (zero hook)", func() {
		_, err := s.RunTx(&types.MsgCreateAggregationHook{
			Owner: creator.Address,
			Hooks: []util.HexAddress{util.NewZeroAddress()},
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("hook id cannot be zero"))
	})

	It("Create (invalid) Aggregation Hook (non-existing hook)", func() {
		nonExistingHook := util.CreateMockHexAddress("missing-hook", 1)

		_, err := s.RunTx(&types.MsgCreateAggregationHook{
			Owner: creator.Address,
			Hooks: []util.HexAddress{nonExistingHook},
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(nonExistingHook.String()))
	})

	It("Create (invalid) Aggregation Hook (duplicate hook)", func() {
		_, err := s.RunTx(&types.MsgCreateAggregationHook{
			Owner: creator.Address,
			Hooks: []util.HexAddress{noopHookId, noopHookId},
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("duplicate hook id"))
	})

	It("Create (invalid) Aggregation Hook (nested aggregation hook)", func() {
		_, err := s.RunTx(&types.MsgCreateAggregationHook{
			Owner: creator.Address,
			Hooks: []util.HexAddress{aggregationHookId},
		})

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("nested aggregation hook is not allowed"))
	})

	It("AggregationHook HookType", func() {
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(aggregationHookId)
		Expect(err).To(BeNil())

		Expect((*handler).HookType()).To(Equal(uint8(types.POST_DISPATCH_HOOK_TYPE_AGGREGATION)))
	})

	It("AggregationHook (valid) exists", func() {
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(aggregationHookId)
		Expect(err).To(BeNil())

		exists, err := (*handler).Exists(s.Ctx(), aggregationHookId)
		Expect(err).To(BeNil())
		Expect(exists).To(BeTrue())
	})

	It("AggregationHook (invalid) exists", func() {
		aggregationHookId[31] = byte(10)
		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(aggregationHookId)
		Expect(err).To(BeNil())

		exists, err := (*handler).Exists(s.Ctx(), aggregationHookId)
		Expect(err).To(BeNil())
		Expect(exists).To(BeFalse())
	})

	It("AggregationHook QuoteDispatch sums child quotes", func() {
		igp1, igp2 := createConfiguredIgpPair(s, creator)
		aggregationHookId, err := createDummyAggregationHook(s, creator.Address, []util.HexAddress{igp1, igp2})
		Expect(err).To(BeNil())

		handler, err := s.App().HyperlaneKeeper.PostDispatchRouter().GetModule(aggregationHookId)
		Expect(err).To(BeNil())

		quote, err := (*handler).QuoteDispatch(s.Ctx(), mailboxId, aggregationHookId, util.StandardHookMetadata{GasLimit: math.NewInt(10)}, util.HyperlaneMessage{Destination: 2})

		Expect(err).To(BeNil())
		Expect(quote).To(Equal(sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(20)))))
	})

	It("AggregationHook PostDispatch executes child hooks in order", func() {
		merkleHookId2, err := createDummyMerkleTreeHook(s, creator.Address, mailboxId)
		Expect(err).To(BeNil())

		aggregationHookId, err := createDummyAggregationHook(s, creator.Address, []util.HexAddress{merkleHookId, merkleHookId2})
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       0,
			Origin:      11,
			Sender:      util.CreateMockHexAddress("sender", 1),
			Destination: 2,
			Recipient:   util.CreateMockHexAddress("recipient", 1),
			Body:        []byte("test"),
		}

		charged, err := s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), mailboxId, aggregationHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(BeNil())
		Expect(charged).To(Equal(sdk.NewCoins()))

		var insertedHookIds []string
		for _, event := range s.Ctx().EventManager().Events() {
			if event.Type != "hyperlane.core.post_dispatch.v1.EventInsertedIntoTree" {
				continue
			}
			attr, found := event.GetAttribute("merkle_tree_hook_id")
			if found {
				insertedHookIds = append(insertedHookIds, attr.Value)
			}
		}

		Expect(insertedHookIds).To(Equal([]string{`"` + merkleHookId.String() + `"`, `"` + merkleHookId2.String() + `"`}))
	})

	It("AggregationHook PostDispatch passes remaining maxFee to child hooks", func() {
		igp1, igp2 := createConfiguredIgpPair(s, creator)
		aggregationHookId, err := createDummyAggregationHook(s, creator.Address, []util.HexAddress{igp1, igp2})
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       0,
			Origin:      11,
			Sender:      util.CreateMockHexAddress("sender", 1),
			Destination: 2,
			Recipient:   util.CreateMockHexAddress("recipient", 1),
			Body:        []byte("test"),
		}
		metadata := util.StandardHookMetadata{
			Address:  creator.AccAddress,
			GasLimit: math.NewInt(10),
		}

		charged, err := s.App().HyperlaneKeeper.PostDispatch(
			s.Ctx(),
			mailboxId,
			aggregationHookId,
			metadata,
			message,
			sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(20))),
		)
		Expect(err).To(BeNil())
		Expect(charged).To(Equal(sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(20)))))

		_, err = s.App().HyperlaneKeeper.PostDispatch(
			s.Ctx(),
			mailboxId,
			aggregationHookId,
			metadata,
			message,
			sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(15))),
		)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("required payment exceeds max hyperlane fee"))
	})

	It("AggregationHook PostDispatch returns child failure", func() {
		wrongMailboxId, err := createDummyMailbox(s, creator.Address)
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       0,
			Origin:      11,
			Sender:      util.CreateMockHexAddress("sender", 1),
			Destination: 2,
			Recipient:   util.CreateMockHexAddress("recipient", 1),
			Body:        []byte("test"),
		}

		_, err = s.App().HyperlaneKeeper.PostDispatch(s.Ctx(), wrongMailboxId, aggregationHookId, util.StandardHookMetadata{}, message, sdk.NewCoins())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(types.ErrSenderIsNotDesignatedMailbox.Error()))
	})

	It("AggregationHook genesis export preserves ordered child hooks", func() {
		genesis := keeper.ExportGenesis(s.Ctx(), s.App().HyperlaneKeeper.PostDispatchKeeper)

		Expect(genesis.AggregationHooks).To(HaveLen(1))
		Expect(genesis.AggregationHooks[0].Id).To(Equal(aggregationHookId))
		Expect(genesis.AggregationHooks[0].Hooks).To(Equal([]util.HexAddress{noopHookId, merkleHookId}))
	})
})

func createConfiguredIgpPair(s *i.KeeperTestSuite, creator i.TestValidatorAddress) (util.HexAddress, util.HexAddress) {
	igp1, err := createDummyIgp(s, creator.Address, "acoin")
	Expect(err).To(BeNil())
	igp2, err := createDummyIgp(s, creator.Address, "acoin")
	Expect(err).To(BeNil())

	for _, igp := range []util.HexAddress{igp1, igp2} {
		_, err = s.RunTx(&types.MsgSetDestinationGasConfig{
			Owner: creator.Address,
			IgpId: igp,
			DestinationGasConfig: &types.DestinationGasConfig{
				RemoteDomain: 2,
				GasOracle: &types.GasOracle{
					TokenExchangeRate: math.NewInt(1e10),
					GasPrice:          math.NewInt(1),
				},
				GasOverhead: math.ZeroInt(),
			},
		})
		Expect(err).To(BeNil())
	}

	return igp1, igp2
}
