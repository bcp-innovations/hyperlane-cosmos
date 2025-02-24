package types

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

type CoreKeeper interface {
	LocalDomain(ctx context.Context) (uint32, error)
	MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error)
}
