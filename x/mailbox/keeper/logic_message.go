package keeper

import (
	"cosmossdk.io/collections"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/mailbox/types"
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
	received, err := k.Messages.Has(ctx, collections.Join(message.Recipient.Bytes(), message.Id().Bytes()))
	if err != nil {
		return err
	}
	if received {
		return fmt.Errorf("already received messsage")
	}
	err = k.Messages.Set(ctx, collections.Join(message.Recipient.Bytes(), message.Id().Bytes()))
	if err != nil {
		return err
	}

	rawIsmAddress, err := k.ReceiverIsmMapping.Get(ctx, message.Recipient.Bytes())
	if err != nil {
		return err
	}

	ismId := util.HexAddress(rawIsmAddress)

	verified, err := k.ismKeeper.Verify(ctx, ismId, metadata, message)
	if err != nil {
		return err
	}
	if !verified {
		return fmt.Errorf("threshold not reached")
	}

	_ = k.Hooks().Handle(ctx, mailboxId, message.Origin, message.Sender, message)

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
) (messageId util.HexAddress, error error) {
	mailbox, err := k.Mailboxes.Get(ctx, originMailboxId.Bytes())
	if err != nil {
		return util.HexAddress{}, err
	}

	mailbox.MessageSent++

	hypMsg := types.HyperlaneMessage{
		Version:     1,
		Nonce:       mailbox.MessageSent,
		Origin:      k.LocalDomain(),
		Sender:      sender,
		Destination: destinationDomain,
		Recipient:   recipient,
		Body:        body,
	}

	tree, err := types.TreeFromProto(mailbox.Tree)
	if err != nil {
		return util.HexAddress{}, err
	}

	count := tree.GetCount()

	if err = tree.Insert(hypMsg.Id()); err != nil {
		return util.HexAddress{}, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.InsertedIntoTree{
		MessageId: hypMsg.String(),
		Index:     count,
	})

	mailbox.Tree = types.ProtoFromTree(tree)

	err = k.Mailboxes.Set(ctx, originMailboxId.Bytes(), mailbox)
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

func (k Keeper) LocalDomain() uint32 {
	// TODO use global param
	return 100
}
