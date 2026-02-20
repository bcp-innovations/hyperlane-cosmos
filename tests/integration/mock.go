package integration

import (
	"context"

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
	n.hooks[sequence] = struct{}{}
	return sequence, nil
}

func (n NoopPostDispatchHookHandler) Exists(_ context.Context, hookId util.HexAddress) (bool, error) {
	_, ok := n.hooks[hookId]
	return ok, nil
}

func (n NoopPostDispatchHookHandler) PostDispatch(ctx context.Context, _, hookId util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage, _ sdk.Coins) (sdk.Coins, error) {
	if has, err := n.Exists(ctx, hookId); err != nil || !has {
		return sdk.Coins{}, err
	}

	return sdk.NewCoins(), nil
}

func (n NoopPostDispatchHookHandler) QuoteDispatch(_ context.Context, _, _ util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage) (sdk.Coins, error) {
	return sdk.NewCoins(), nil
}

func (n NoopPostDispatchHookHandler) HookType() uint8 {
	return MOCK_TYPE_NOOP_POST_DISPATCH
}

const MOCK_TYPE_APP uint8 = 201

type CallInfo struct {
	count     int
	message   util.HyperlaneMessage
	mailboxId util.HexAddress
}
type MockApp struct {
	// Map from recipient to ISM ID
	apps     map[util.HexAddress]util.HexAddress
	router   *util.Router[util.HyperlaneApp]
	callinfo *CallInfo
	moduleId uint8
}

func CreateMockApp(router *util.Router[util.HyperlaneApp]) *MockApp {
	handler := MockApp{
		apps:     make(map[util.HexAddress]util.HexAddress),
		router:   router,
		callinfo: new(CallInfo),
		moduleId: MOCK_TYPE_APP,
	}

	router.RegisterModule(handler.moduleId, handler)

	return &handler
}

func (m MockApp) RegisterApp(ctx context.Context, ismId util.HexAddress) (util.HexAddress, error) {
	sequence, err := m.router.GetNextSequence(ctx, m.moduleId)
	if err != nil {
		return util.HexAddress{}, err
	}
	m.apps[sequence] = ismId
	return sequence, nil
}

func (m MockApp) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
	*m.callinfo = CallInfo{
		count:     m.callinfo.count + 1,
		message:   message,
		mailboxId: mailboxId,
	}
	return nil
}

func (m MockApp) CallInfo() (count int, message util.HyperlaneMessage, mailboxId util.HexAddress) {
	if m.callinfo == nil {
		return 0, util.HyperlaneMessage{}, util.HexAddress{}
	}
	return m.callinfo.count, m.callinfo.message, m.callinfo.mailboxId
}

func (m MockApp) Exists(_ context.Context, recipient util.HexAddress) (bool, error) {
	_, ok := m.apps[recipient]
	return ok, nil
}

func (m MockApp) ReceiverIsmId(_ context.Context, recipient util.HexAddress) (*util.HexAddress, error) {
	ismId, ok := m.apps[recipient]
	if !ok {
		return nil, nil
	}
	return &ismId, nil
}

const MOCK_ISM uint8 = 202

type MockIsm struct {
	router       *util.Router[util.InterchainSecurityModule]
	isms         map[util.HexAddress]bool // maps ISM ID to verification result
	calls        *int
	moduleId     uint8
	shouldVerify bool
}

func CreateMockIsm(router *util.Router[util.InterchainSecurityModule]) *MockIsm {
	handler := MockIsm{
		isms:         make(map[util.HexAddress]bool),
		router:       router,
		calls:        new(int),
		moduleId:     MOCK_ISM,
		shouldVerify: true, // default to passing verification
	}

	router.RegisterModule(handler.moduleId, handler)

	return &handler
}

func (m *MockIsm) SetShouldVerify(shouldVerify bool) {
	m.shouldVerify = shouldVerify
}

func (m MockIsm) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	_, ok := m.isms[ismId]
	if !ok {
		return false, nil
	}
	return true, nil
}

func (m MockIsm) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	*m.calls++
	result, ok := m.isms[ismId]
	if !ok {
		return false, nil
	}
	return result, nil
}

func (m *MockIsm) RegisterIsm(ctx context.Context) (util.HexAddress, error) {
	sequence, err := m.router.GetNextSequence(ctx, m.moduleId)
	if err != nil {
		return util.HexAddress{}, err
	}
	m.isms[sequence] = m.shouldVerify
	return sequence, nil
}

func (m MockIsm) CallCount() int {
	if m.calls == nil {
		return 0
	}
	return *m.calls
}
