package types

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// combine multiple mailbox hooks, all hook functions are run in array sequence
var _ MailboxHooks = MultiMailboxHooks{}

type MultiMailboxHooks []MailboxHooks

func NewMultiMailboxHooks(hooks ...MailboxHooks) MultiMailboxHooks {
	return hooks
}

func (h MultiMailboxHooks) Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message HyperlaneMessage) error {
	for i := range h {
		if err := h[i].Handle(ctx, mailboxId, origin, sender, message); err != nil {
			return err
		}
	}

	return nil
}
