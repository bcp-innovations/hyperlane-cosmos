package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

//
// Interchain Security Hooks

type InterchainSecurityHooks interface {
	Verify(ctx sdk.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error)
}

type InterchainSecurityHooksWrapper struct{ InterchainSecurityHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (InterchainSecurityHooksWrapper) IsOnePerModuleType() {}

//
// Post Dispatch Hooks

type PostDispatchHooks interface {
	PostDispatch(ctx sdk.Context, hookId util.HexAddress, metadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
}

type PostDispatchHooksWrapper struct{ PostDispatchHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (PostDispatchHooksWrapper) IsOnePerModuleType() {}

//
// Mailbox Hooks

type MailboxHooks interface {
	Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error
	ReceiverIsmId(ctx context.Context, recipient util.HexAddress) (util.HexAddress, error)
}

type MailboxHooksWrapper struct{ MailboxHooks }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (MailboxHooksWrapper) IsOnePerModuleType() {}

//
// External Keepers

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
}
