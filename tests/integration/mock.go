package integration

import (
	"context"

	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ util.PostDispatchModule = NoopPostDispatchHookHandler{}

type NoopPostDispatchHookHandler struct {
	hooks  map[util.HexAddress]struct{}
	router *util.Router[util.PostDispatchModule]
}

const MOCK_TYPE_NOOP_POST_DISPATCH uint8 = 200

func CreateNoopDispatchHookHandler(router *util.Router[util.PostDispatchModule]) *NoopPostDispatchHookHandler {
	handler := NoopPostDispatchHookHandler{
		hooks:  make(map[util.HexAddress]struct{}),
		router: router,
	}

	router.RegisterModule(MOCK_TYPE_NOOP_POST_DISPATCH, handler)

	return &handler
}

func (n NoopPostDispatchHookHandler) CreateHook(ctx context.Context) (util.HexAddress, error) {
	sequence, err := n.router.GetNextSequence(ctx, MOCK_TYPE_NOOP_POST_DISPATCH)
	if err != nil {
		return util.HexAddress{}, err
	}
	return sequence, nil
}

func (n NoopPostDispatchHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	_, ok := n.hooks[hookId]
	return ok, nil
}

func (n NoopPostDispatchHookHandler) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, metadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	if has, err := n.Exists(ctx, hookId); err != nil || !has {
		return sdk.Coins{}, err
	}

	return sdk.NewCoins(), nil
}

func (n NoopPostDispatchHookHandler) QuoteDispatch(ctx context.Context, hookId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (math.Int, error) {
	return math.ZeroInt(), nil
}

func (n NoopPostDispatchHookHandler) HookType() uint8 {
	return MOCK_TYPE_NOOP_POST_DISPATCH
}

func (n NoopPostDispatchHookHandler) SupportsMetadata(metadata []byte) bool {
	return false
}
