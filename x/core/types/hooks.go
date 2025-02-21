package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	//"fmt"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// TODO refactor to Mailbox Client
var _ MailboxHooks = MultiMailboxHooks{}

type MultiMailboxHooks []MailboxHooks

func NewMultiMailboxHooks(hooks ...MailboxHooks) MultiMailboxHooks {
	return hooks
}

func (h MultiMailboxHooks) Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message util.HyperlaneMessage) error {
	for i := range h {
		if err := h[i].Handle(ctx, mailboxId, origin, sender, message); err != nil {
			return err
		}
	}

	return nil
}

// Interchain Security Module Multi Wrapper

// combine multiple mailbox hooks, all hook functions are run in array sequence
var _ InterchainSecurityHooks = MultiInterchainSecurityHooks{}

type MultiInterchainSecurityHooks []InterchainSecurityHooks

func NewMultiInterchainSecurityHooks(hooks ...InterchainSecurityHooks) MultiInterchainSecurityHooks {
	return hooks
}

func (h MultiInterchainSecurityHooks) Verify(ctx sdk.Context, ismId util.HexAddress, metadata any, message util.HyperlaneMessage) (bool, error) {
	for i := range h {
		verfied, err := h[i].Verify(ctx, ismId, metadata, message)
		if err != nil {
			return false, err
		}
		if verfied {
			return true, nil
		}
	}

	return false, nil
}

// Post Dispatch Hook Multi Wrapper

// combine multiple mailbox hooks, all hook functions are run in array sequence
var _ PostDispatchHooks = MultiPostDispatchHooks{}

type MultiPostDispatchHooks []PostDispatchHooks

func NewMultiPostDispatchHooks(hooks ...PostDispatchHooks) MultiPostDispatchHooks {
	return hooks
}

func (m MultiPostDispatchHooks) PostDispatch(ctx sdk.Context, hookId util.HexAddress, metadata any, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	for i := range m {
		coins, err := m[i].PostDispatch(ctx, hookId, metadata, message, maxFee)
		if err != nil {
			return nil, err
		}
		if coins != nil {
			return coins, nil
		}
	}
	return nil, nil
}
