package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/ism/types"
	mailboxTypes "github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
)

func (ms msgServer) VerifyDryRun(ctx context.Context, req *types.MsgVerifyDryRun) (*types.MsgVerifyDryRunResponse, error) {
	rawMessage, err := util.DecodeEthHex(req.Message)
	if err != nil {
		return nil, err
	}

	message, err := mailboxTypes.ParseHyperlaneMessage(rawMessage)
	if err != nil {
		return nil, err
	}

	metadata, err := util.DecodeEthHex(req.Metadata)
	if err != nil {
		return nil, err
	}

	ismId, err := util.DecodeHexAddress(req.IsmId)
	if err != nil {
		return nil, err
	}

	verified, err := ms.k.Verify(ctx, ismId, metadata, message)
	if err != nil {
		return nil, err
	}

	return &types.MsgVerifyDryRunResponse{
		Verified: verified,
	}, err
}
