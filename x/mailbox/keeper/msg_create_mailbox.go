package keeper

import (
	"context"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) CreateMailbox(ctx context.Context, req *types.MsgCreateMailbox) (*types.MsgCreateMailboxResponse, error) {
	igpId, err := util.DecodeHexAddress(req.Igp.Id)
	if err != nil {
		return nil, err
	}

	exists, err := ms.k.igpKeeper.IgpIdExists(ctx, igpId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("InterchainGasPaymaster with id %s does not exist", igpId.String())
	}

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
		Igp: &types.InterchainGasPaymaster{
			Id:       req.Igp.Id,
			Required: req.Igp.Required,
		},
		Tree: types.ProtoFromTree(tree),
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId.Bytes(), newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{}, nil
}
