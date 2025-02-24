package keeper

import (
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (ms msgServer) CreateMailbox(ctx context.Context, req *types.MsgCreateMailbox) (*types.MsgCreateMailboxResponse, error) {
	ismId, err := util.DecodeHexAddress(req.DefaultIsm)
	if err != nil {
		return nil, fmt.Errorf("ism id %s is invalid: %s", req.DefaultIsm, err.Error())
	}

	exists, err := ms.k.IsmKeeper.IsmIdExists(ctx, ismId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("ism with id %s does not exist", ismId.String())
	}

	igpId, err := util.DecodeHexAddress(req.Igp.Id)
	if err != nil {
		return nil, fmt.Errorf("igp id %s is invalid: %s", req.Igp.Id, err.Error())
	}

	exists, err = ms.k.IgpIdExists(ctx, igpId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("igp with id %s does not exist", igpId.String())
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
		DefaultIsm:      ismId.String(),
		Igp: &types.InterchainGasPaymaster{
			Id:       req.Igp.Id,
			Required: req.Igp.Required,
		},
		Tree: types.ProtoFromTree(tree),
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId.Bytes(), newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{Id: prefixedId.String()}, nil
}

func (ms msgServer) DispatchMessage(ctx context.Context, req *types.MsgDispatchMessage) (*types.MsgDispatchMessageResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	bodyBytes, err := hexutil.Decode(req.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid body: %s", err)
	}

	mailBoxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, fmt.Errorf("invalid mailbox id: %s", err)
	}

	sender, err := util.ParseFromCosmosAcc(req.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender: %s", err)
	}

	recipient, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient: %s", err)
	}

	msgId, err := ms.k.DispatchMessage(goCtx, mailBoxId, req.Destination, recipient, sender, bodyBytes, req.Sender, req.IgpId, req.GasLimit, req.MaxFee)
	if err != nil {
		return nil, err
	}

	return &types.MsgDispatchMessageResponse{
		MessageId: msgId.String(),
	}, nil
}

func (ms msgServer) ProcessMessage(ctx context.Context, req *types.MsgProcessMessage) (*types.MsgProcessMessageResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

	mailboxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, fmt.Errorf("invalid mailbox id: %s", err)
	}

	// Decode and parse message
	messageBytes, err := util.DecodeEthHex(req.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message")
	}

	if len(messageBytes) == 0 {
		return nil, fmt.Errorf("invalid message")
	}

	// Decode and parse metadata
	metadataBytes, err := util.DecodeEthHex(req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode metadata")
	}

	if err = ms.k.ProcessMessage(goCtx, mailboxId, messageBytes, metadataBytes); err != nil {
		return nil, err
	}

	return &types.MsgProcessMessageResponse{}, nil
}
