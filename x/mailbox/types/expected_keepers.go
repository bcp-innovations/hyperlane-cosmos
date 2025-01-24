package types

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// Event Hooks
// These can be utilized to communicate between a warp keeper and another
// keeper which must take particular actions

type MailboxHooks interface {
	Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message HyperlaneMessage) error
}

type MailboxHooksWrapper struct{ MailboxHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (MailboxHooksWrapper) IsOnePerModuleType() {}
