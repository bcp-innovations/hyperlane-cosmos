package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ProcessMessage(ctx sdk.Context, mailboxIdString string, rawMessage []byte, metadata []byte) error {
	message, err := types.ParseHyperlaneMessage(rawMessage)
	if err != nil {
		return err
	}

	mailboxId, err := util.DecodeHexAddress(mailboxIdString)
	if err != nil {
		return err
	}

	// Check if mailbox exists and increment counter
	mailbox, err := k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return err
	}
	mailbox.MessageReceived++

	err = k.Mailboxes.Set(ctx, mailboxId.Bytes(), mailbox)
	if err != nil {
		return err
	}

	// Check replay protection
	received, err := k.Messages.Has(ctx, message.Id().Bytes())
	if err != nil {
		return err
	}
	if received {
		return fmt.Errorf("already received messsage")
	}
	err = k.Messages.Set(ctx, message.Id().Bytes())
	if err != nil {
		return err
	}

	rawIsmAddress, err := k.ReceiverIsmMapping.Get(ctx, message.Recipient.Bytes())
	if err != nil {
		return err
	}

	ismId := util.HexAddress(rawIsmAddress)

	verified, err := k.Verify(ctx, ismId, metadata, message)
	if err != nil {
		return err
	}
	if !verified {
		return fmt.Errorf("threshold not reached")
	}

	err = k.Hooks().Handle(ctx, mailboxId, message.Origin, message.Sender, message)
	if err != nil {
		return err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.Process{
		OriginMailboxId: mailboxIdString,
		Origin:          message.Origin,
		Sender:          message.Sender.String(),
		Recipient:       message.Recipient.String(),
		MessageId:       message.Id().String(),
		Message:         message.String(),
	})

	return nil
}

func (k Keeper) DispatchMessage(
	ctx sdk.Context,
	originMailboxId util.HexAddress,
	destinationDomain uint32,
	// Recipient address on the destination chain (e.g. smart contract)
	recipient util.HexAddress,
	// sender address on the origin chain (e.g. token id)
	sender util.HexAddress,
	body []byte,
	// Custom IGP settings
	cosmosSender string,
	customIgpId string,
	gasLimit math.Int,
	maxFee math.Int,
) (messageId util.HexAddress, error error) {
	mailbox, err := k.Mailboxes.Get(ctx, originMailboxId.Bytes())
	if err != nil {
		return util.HexAddress{}, err
	}

	hypMsg := types.HyperlaneMessage{
		Version:     3,
		Nonce:       mailbox.MessageSent,
		Origin:      k.LocalDomain(ctx),
		Sender:      sender,
		Destination: destinationDomain,
		Recipient:   recipient,
		Body:        body,
	}
	mailbox.MessageSent++

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return util.HexAddress{}, err
	}

	count := tree.GetCount()

	if err = tree.Insert(hypMsg.Id()); err != nil {
		return util.HexAddress{}, err
	}

	err = k.Messages.Set(ctx, hypMsg.Id().Bytes())
	if err != nil {
		return util.HexAddress{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.InsertedIntoTree{
		MessageId: hypMsg.Id().String(),
		Index:     count,
		MailboxId: mailbox.Id,
	})

	mailbox.Tree = types.ProtoFromTree(tree)

	err = k.Mailboxes.Set(ctx, originMailboxId.Bytes(), mailbox)
	if err != nil {
		return util.HexAddress{}, err
	}

	// Interchain Gas Payment
	igpId, err := util.DecodeHexAddress(mailbox.Igp.Id)
	if err != nil {
		return util.HexAddress{}, err
	}

	if !mailbox.Igp.Required && customIgpId != "" {
		igpId, err = util.DecodeHexAddress(customIgpId)
		if err != nil {
			return util.HexAddress{}, nil
		}
	}

	err = k.PayForGas(ctx, cosmosSender, igpId, hypMsg.Id().String(), destinationDomain, gasLimit, maxFee)
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.Dispatch{
		OriginMailboxId: originMailboxId.String(),
		Sender:          sender.String(),
		Destination:     destinationDomain,
		Recipient:       recipient.String(),
		Message:         hypMsg.String(),
	})

	return hypMsg.Id(), nil
}

func (k Keeper) LocalDomain(ctx context.Context) uint32 {
	params, err := k.Params.Get(ctx)
	if err != nil {
		// fallback to default if params can't be read
		return 75898669
	}
	return params.Domain
}
