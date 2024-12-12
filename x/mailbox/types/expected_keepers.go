package types

import (
	"context"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
)

// Event Hooks
// These can be utilized to communicate between a warp keeper and another
// keeper which must take particular actions

type MailboxHooks interface {
	// TODO should we return an error?
	Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message HyperlaneMessage) error
}

type MailboxHooksWrapper struct{ MailboxHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (MailboxHooksWrapper) IsOnePerModuleType() {}

type IsmKeeper interface {
	IsmIdExists(ctx context.Context, ismId string) (bool, error)
	Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message HyperlaneMessage) (valid bool, err error)
}
