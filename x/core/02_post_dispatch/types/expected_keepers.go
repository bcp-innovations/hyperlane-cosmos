package types

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CoreKeeper interface {
	LocalDomain(ctx context.Context, mailboxId util.HexAddress) (uint32, error)
	MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error)
	PostDispatchRouter() *util.Router[util.PostDispatchModule]
}

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}
