package types

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

type CoreKeeper interface {
	LocalDomain(ctx context.Context, mailboxId util.HexAddress) (uint32, error)
	MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error)
	IsmExists(ctx context.Context, ismId util.HexAddress) (bool, error)
	IsmRouter() *util.Router[util.InterchainSecurityModule]
	Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error)
}
