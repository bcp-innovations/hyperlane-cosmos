package keeper

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/igp/types"
)

func (ms msgServer) PayForGas(ctx context.Context, req *types.MsgPayForGas) (*types.MsgPayForGasResponse, error) {
	igpId, err := util.DecodeHexAddress(req.IgpId)
	if err != nil {
		return nil, err
	}

	return &types.MsgPayForGasResponse{}, ms.k.PayForGas(ctx, req.Sender, igpId, req.MessageId, req.DestinationDomain, req.GasLimit, req.MaxFee)
}
