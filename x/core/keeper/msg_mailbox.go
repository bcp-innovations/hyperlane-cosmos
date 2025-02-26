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

	// Check ism existence
	if err := ms.k.AssertIsmExists(ctx, req.DefaultIsm); err != nil {
		return nil, err
	}

	// Check default is valid if set
	if req.DefaultHook != "" {
		if err := ms.k.AssertPostDispatchHookExists(ctx, req.DefaultHook); err != nil {
			return nil, err
		}
	}

	// Check required hook is valid if set.
	// The "required" means that this hook can not be overridden by the message dispatcher
	if req.RequiredHook != "" {
		if err := ms.k.AssertPostDispatchHookExists(ctx, req.RequiredHook); err != nil {
			return nil, err
		}
	}

	mailboxCount, err := ms.k.MailboxesSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	prefixedId := util.CreateHexAddress(types.ModuleName, int64(mailboxCount))

	newMailbox := types.Mailbox{
		Id:              prefixedId.String(),
		Creator:         req.Creator,
		MessageSent:     0,
		MessageReceived: 0,
		DefaultIsm:      req.DefaultIsm,
		DefaultHook:     req.DefaultHook,
		RequiredHook:    req.RequiredHook,
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId.Bytes(), newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{Id: prefixedId.String()}, nil
}

// DispatchMessage assumes an Interchain GasPaymaster as a hook, as there are currently no other hooks available
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

	accSender, err := sdk.AccAddressFromBech32(req.Sender)
	if err != nil {
		return nil, fmt.Errorf("invalid sender: %s", err)
	}

	recipient, err := util.DecodeHexAddress(req.Recipient)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient: %s", err)
	}

	customIgpId, err := util.DecodeHexAddress(req.CustomIgp)
	if err != nil {
		return nil, fmt.Errorf("invalid customIgp: %s", err)
	}

	msgId, err := ms.k.DispatchMessage(goCtx, mailBoxId, sender, sdk.NewCoins(req.MaxFee), req.Destination, recipient, bodyBytes, util.StandardHookMetadata{
		Variant:  1,
		Value:    req.MaxFee.Amount,
		GasLimit: req.GasLimit,
		Address:  accSender,
	}.Bytes(), customIgpId)
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
