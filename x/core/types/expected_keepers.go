package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

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

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr []byte, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr []byte, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}
