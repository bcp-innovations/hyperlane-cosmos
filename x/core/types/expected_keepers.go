package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// Interchain Security Hooks

type InterchainSecurityHooks interface {
	Verify(ctx sdk.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error)
}

type InterchainSecurityHooksWrapper struct{ InterchainSecurityHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (InterchainSecurityHooksWrapper) IsOnePerModuleType() {}

// PostDispatchHooks

type PostDispatchHooks interface {
	PostDispatch(ctx sdk.Context, hookId util.HexAddress, metadata any, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
}

type PostDispatchHooksWrapper struct{ PostDispatchHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (PostDispatchHooksWrapper) IsOnePerModuleType() {}

// Event Hooks
// These can be utilized to communicate between a warp keeper and another
// keeper which must take particular actions
// TODO

type MailboxHooks interface {
	Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message util.HyperlaneMessage) error
}

type MailboxHooksWrapper struct{ MailboxHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (MailboxHooksWrapper) IsOnePerModuleType() {}

// External Keepers

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}
