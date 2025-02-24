package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	//"fmt"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ MailboxHooks = MultiMailboxHooks{}

type MultiMailboxHooks []MailboxHooks

func NewMultiMailboxHooks(hooks ...MailboxHooks) MultiMailboxHooks {
	return hooks
}

func (h MultiMailboxHooks) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
	for i := range h {
		if err := h[i].Handle(ctx, mailboxId, message); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiMailboxHooks) ReceiverIsmId(ctx context.Context, recipient util.HexAddress) (util.HexAddress, error) {
	var receiverIsm util.HexAddress
	for i := range h {
		ismId, err := h[i].ReceiverIsmId(ctx, recipient)
		if err != nil {
			return util.HexAddress{}, nil
		}
		if !ismId.IsZeroAddress() {
			if receiverIsm.IsZeroAddress() {
				receiverIsm = ismId
			} else {
				return util.HexAddress{}, ErrMultipleReceiverIsm
			}
		}
	}

	return util.HexAddress{}, ErrNoReceiverISM
}

// Interchain Security Module Multi Wrapper

// combine multiple mailbox hooks, all hook functions are run in array sequence
var _ InterchainSecurityHooks = MultiInterchainSecurityHooks{}

type MultiInterchainSecurityHooks []InterchainSecurityHooks

func NewMultiInterchainSecurityHooks(hooks ...InterchainSecurityHooks) MultiInterchainSecurityHooks {
	return hooks
}

func (h MultiInterchainSecurityHooks) Verify(ctx sdk.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
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
