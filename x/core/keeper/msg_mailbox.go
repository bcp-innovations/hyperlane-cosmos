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
	if req.DefaultHook != nil {
		if err := ms.k.AssertPostDispatchHookExists(ctx, *req.DefaultHook); err != nil {
			return nil, err
		}
	}

	// Check required hook is valid if set.
	// The "required" means that this hook can not be overridden by the message dispatcher
	if req.RequiredHook != nil {
		if err := ms.k.AssertPostDispatchHookExists(ctx, *req.RequiredHook); err != nil {
			return nil, err
		}
	}

	mailboxCount, err := ms.k.MailboxesSequence.Next(ctx)
	if err != nil {
		return nil, err
	}

	identifier := [20]byte{}
	copy(identifier[:], types.ModuleName)

	// generate a new unique id that is compliant with the way the router generates ids
	prefixedId := util.GenerateHexAddress(identifier, uint32(types.ModuleId), mailboxCount)

	newMailbox := types.Mailbox{
		Id:              prefixedId,
		Owner:           req.Owner,
		MessageSent:     0,
		MessageReceived: 0,
		DefaultIsm:      req.DefaultIsm,
		DefaultHook:     req.DefaultHook,
		RequiredHook:    req.RequiredHook,
	}

	if err = ms.k.Mailboxes.Set(ctx, prefixedId.GetInternalId(), newMailbox); err != nil {
		return nil, err
	}

	return &types.MsgCreateMailboxResponse{Id: prefixedId.String()}, nil
}

// DispatchMessage assumes an Interchain GasPaymaster as a hook, as there are currently no other hooks available
// TODO: Remove
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

	_, err = util.DecodeHexAddress(req.CustomIgp)
	if req.CustomIgp != "" && err != nil {
		return nil, fmt.Errorf("invalid customIgp: %s", err)
	}

	msgId, err := ms.k.DispatchMessage(goCtx, mailBoxId, sender, sdk.NewCoins(req.MaxFee), req.Destination, recipient, bodyBytes, util.StandardHookMetadata{
		GasLimit:           req.GasLimit,
		Address:            accSender,
		CustomHookMetadata: []byte{},
	}, nil)
	if err != nil {
		return nil, err
	}

	return &types.MsgDispatchMessageResponse{
		MessageId: msgId.String(),
	}, nil
}

func (ms msgServer) ProcessMessage(ctx context.Context, req *types.MsgProcessMessage) (*types.MsgProcessMessageResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)

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

	if err = ms.k.ProcessMessage(goCtx, req.MailboxId, messageBytes, metadataBytes); err != nil {
		return nil, err
	}

	return &types.MsgProcessMessageResponse{}, nil
}

func (ms msgServer) SetMailbox(ctx context.Context, req *types.MsgSetMailbox) (*types.MsgSetMailboxResponse, error) {
	mailboxId := req.MailboxId
	mailbox, err := ms.k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find mailbox with id: %v", mailboxId.String())
	}

	if mailbox.Owner != req.Owner {
		return nil, fmt.Errorf("%s does not own mailbox with id %s", req.Owner, mailboxId.String())
	}

	if req.DefaultIsm != nil {
		if err = ms.k.AssertIsmExists(ctx, *req.DefaultIsm); err != nil {
			return nil, fmt.Errorf("ism with id %s does not exist", req.DefaultIsm)
		}

		mailbox.DefaultIsm = *req.DefaultIsm
	}

	if req.DefaultHook != nil {
		if err := ms.k.AssertPostDispatchHookExists(ctx, *req.DefaultHook); err != nil {
			return nil, err
		}
		mailbox.DefaultHook = req.DefaultHook
	}

	if req.RequiredHook != nil {
		if err := ms.k.AssertPostDispatchHookExists(ctx, *req.RequiredHook); err != nil {
			return nil, err
		}
		mailbox.RequiredHook = req.RequiredHook
	}

	if req.NewOwner != "" {
		mailbox.Owner = req.NewOwner
	}

	if err = ms.k.Mailboxes.Set(ctx, mailboxId.GetInternalId(), mailbox); err != nil {
		return nil, err
	}

	return &types.MsgSetMailboxResponse{}, nil
}
