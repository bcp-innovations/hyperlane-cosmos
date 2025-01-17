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

type IsmKeeper interface {
	IsmIdExists(ctx context.Context, ismId util.HexAddress) (bool, error)
	Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message HyperlaneMessage) (valid bool, err error)
}

type IgpKeeper interface {
	IgpIdExists(ctx context.Context, igpId util.HexAddress) (bool, error)
	PayForGas(ctx context.Context, sender string, igpId util.HexAddress, messageId string, destinationDomain uint32, gasLimit uint64, maxFee uint64) error
}
