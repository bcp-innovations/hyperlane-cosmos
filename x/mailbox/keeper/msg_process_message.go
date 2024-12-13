package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (ms msgServer) ProcessMessage(ctx context.Context, req *types.MsgProcessMessage) (*types.MsgProcessMessageResponse, error) {

	goCtx := sdk.UnwrapSDKContext(ctx)

	// Decode and parse message
	messageBytes, err := util.DecodeEthHex(req.Message)
	if err != nil {
		return nil, err
	}

	// Decode and parse metadata
	metadataBytes, err := util.DecodeEthHex(req.Metadata)
	if err != nil {
		return nil, err
	}

	if err = ms.k.ProcessMessage(goCtx, req.MailboxId, messageBytes, metadataBytes); err != nil {
		return nil, err
	}

	return &types.MsgProcessMessageResponse{}, nil
}
