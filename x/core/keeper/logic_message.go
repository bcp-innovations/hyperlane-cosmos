package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ProcessMessage verifies and processes an incoming message.
// It checks mailbox existence, prevents replay attacks, verifies through the specified ISM,
// and forwards the message to the recipient if valid.
func (k Keeper) ProcessMessage(
	ctx sdk.Context,
	mailboxId util.HexAddress,
	rawMessage []byte,
	metadata []byte,
) error {
	message, err := util.ParseHyperlaneMessage(rawMessage)
	if err != nil {
		return err
	}

	// Check for valid message version
	if message.Version != types.MESSAGE_VERSION {
		return fmt.Errorf("unsupported message version %d", message.Version)
	}

	// Check if mailbox exists and increment counter
	mailbox, err := k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return fmt.Errorf("failed to find mailbox with id: %s", mailboxId.String())
	}
	mailbox.MessageReceived++

	if message.Destination != mailbox.LocalDomain {
		return fmt.Errorf("message destination %v does not match local domain %v", message.Destination, mailbox.LocalDomain)
	}

	err = k.Mailboxes.Set(ctx, mailboxId.GetInternalId(), mailbox)
	if err != nil {
		return err
	}

	// Check replay protection
	key := collections.Join(mailboxId.GetInternalId(), message.Id().Bytes())
	received, err := k.Messages.Has(ctx, key)
	if err != nil {
		return err
	}
	if received {
		return fmt.Errorf("already received messsage with id %s", message.Id().String())
	}
	err = k.Messages.Set(ctx, key)
	if err != nil {
		return err
	}

	// Verify message
	ismId, err := k.ReceiverIsmId(ctx, message.Recipient)
	if err != nil {
		if errors.IsOf(err, types.ErrNoReceiverISM) {
			ismId = mailbox.DefaultIsm
		} else {
			return err
		}
	}

	verified, err := k.Verify(ctx, ismId, metadata, message)
	if err != nil {
		return err
	}
	if !verified {
		return fmt.Errorf("ism verification failed")
	}

	err = k.Handle(ctx, mailboxId, message)
	if err != nil {
		return err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.Process{
		OriginMailboxId: mailboxId.String(),
		Origin:          message.Origin,
		Sender:          message.Sender.String(),
		Recipient:       message.Recipient.String(),
		MessageId:       message.Id().String(),
		Message:         message.String(),
	})

	return nil
}

// DispatchMessage sends a Hyperlane message to a destination chain.
// It verifies the mailbox, constructs and emits the message,
// and calls the required and optional post-dispatch hooks while enforcing max fee limits.
func (k Keeper) DispatchMessage(
	ctx sdk.Context,
	originMailboxId util.HexAddress,
	// sender address on the origin chain (e.g. token id)
	sender util.HexAddress,
	// the maximum amount of tokens the dispatch is allowed to cost
	maxFee sdk.Coins,

	destinationDomain uint32,
	// Recipient address on the destination chain (e.g. smart contract)
	recipient util.HexAddress,
	body []byte,
	// Custom metadata for postDispatch Hook
	metadata util.StandardHookMetadata,
	postDispatchHookId *util.HexAddress,
) (messageId util.HexAddress, error error) {
	mailbox, err := k.Mailboxes.Get(ctx, originMailboxId.GetInternalId())
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to find mailbox with id: %v", originMailboxId.String())
	}

	// check for valid mailbox state
	if mailbox.RequiredHook == nil {
		return util.HexAddress{}, types.ErrRequiredHookNotSet
	}
	if mailbox.DefaultHook == nil {
		return util.HexAddress{}, types.ErrDefaultHookNotSet
	}

	hypMsg := util.HyperlaneMessage{
		Version:     types.MESSAGE_VERSION,
		Nonce:       mailbox.MessageSent,
		Origin:      mailbox.LocalDomain,
		Sender:      sender,
		Destination: destinationDomain,
		Recipient:   recipient,
		Body:        body,
	}
	mailbox.MessageSent++

	err = k.Messages.Set(ctx, collections.Join(originMailboxId.GetInternalId(), hypMsg.Id().Bytes()))
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.Mailboxes.Set(ctx, originMailboxId.GetInternalId(), mailbox)
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.Dispatch{
		OriginMailboxId: originMailboxId.String(),
		Sender:          sender.String(),
		Destination:     destinationDomain,
		Recipient:       recipient.String(),
		Message:         hypMsg.String(),
	})

	chargedCoinsRequired, err := k.PostDispatch(ctx, originMailboxId, *mailbox.RequiredHook, metadata, hypMsg, maxFee)
	if err != nil {
		return util.HexAddress{}, err
	}

	remainingCoins, neg := maxFee.SafeSub(chargedCoinsRequired...)
	if neg {
		return util.HexAddress{}, fmt.Errorf("remaining coins cannot be negative")
	}

	if postDispatchHookId == nil {
		postDispatchHookId = mailbox.DefaultHook
	}
	if postDispatchHookId == nil {
		return util.HexAddress{}, types.ErrDefaultHookNotSet
	}
	chargedCoinsDefault, err := k.PostDispatch(ctx, originMailboxId, *postDispatchHookId, metadata, hypMsg, remainingCoins)
	if err != nil {
		return util.HexAddress{}, err
	}

	chargedCoins := chargedCoinsRequired.Add(chargedCoinsDefault...)

	if chargedCoins.IsAnyGT(maxFee) {
		return util.HexAddress{}, fmt.Errorf("maxFee exceeded %s > %s", chargedCoins.String(), maxFee.String())
	}

	return hypMsg.Id(), nil
}
