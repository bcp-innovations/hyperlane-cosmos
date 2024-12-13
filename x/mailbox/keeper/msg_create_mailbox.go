package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) CreateMailbox(ctx context.Context, req *types.MsgCreateMailbox) (*types.MsgCreateMailboxResponse, error) {
	mailboxCount, err := ms.k.MailboxesSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName, int64(mailboxCount))

	tree := types.NewTree(types.ZeroHashes, 0)

	newMailbox := types.Mailbox{
		Id:              prefixedId.String(),
		MessageSent:     0,
		MessageReceived: 0,
		Creator:         req.Creator,
		Tree:            types.ProtoFromTree(tree),
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId.Bytes(), newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{}, nil
}
