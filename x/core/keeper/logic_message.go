package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/errors"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ProcessMessage(ctx sdk.Context, mailboxId util.HexAddress, rawMessage []byte, metadata []byte) error {
	message, err := util.ParseHyperlaneMessage(rawMessage)
	if err != nil {
		return err
	}

	// Check if mailbox exists and increment counter
	mailbox, err := k.Mailboxes.Get(ctx, mailboxId.Bytes())
	if err != nil {
		return fmt.Errorf("failed to find mailbox with id: %s", mailboxId.String())
	}
	mailbox.MessageReceived++

	err = k.Mailboxes.Set(ctx, mailboxId.Bytes(), mailbox)
	if err != nil {
		return err
	}

	// Check replay protection
	key := collections.Join(mailboxId.Bytes(), message.Id().Bytes())
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

	ismId, err := k.ReceiverIsmId(ctx, message.Recipient)
	if err != nil {
		if errors.IsOf(err, types.ErrNoReceiverISM) {
			ismId, _ = util.DecodeHexAddress(mailbox.DefaultIsm)
		} else {
			return err
		}
	}

	// New logic
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
	metadata []byte,
	postDispatchHookId util.HexAddress,
) (messageId util.HexAddress, error error) {
	mailbox, err := k.Mailboxes.Get(ctx, originMailboxId.Bytes())
	if err != nil {
		return util.HexAddress{}, fmt.Errorf("failed to find mailbox with id: %v", originMailboxId.String())
	}

	localDomain, err := k.LocalDomain(ctx)
	if err != nil {
		return util.HexAddress{}, err
	}

	hypMsg := util.HyperlaneMessage{
		Version:     3,
		Nonce:       mailbox.MessageSent,
		Origin:      localDomain,
		Sender:      sender,
		Destination: destinationDomain,
		Recipient:   recipient,
		Body:        body,
	}
	mailbox.MessageSent++

	err = k.Messages.Set(ctx, collections.Join(originMailboxId.Bytes(), hypMsg.Id().Bytes()))
	if err != nil {
		return util.HexAddress{}, err
	}

	err = k.Mailboxes.Set(ctx, originMailboxId.Bytes(), mailbox)
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

	// TODO: check if we want to allow nil-able hooks
	requiredHookAddress, err := util.DecodeHexAddress(mailbox.RequiredHook)
	if err != nil {
		return util.HexAddress{}, err
	}
	chargedCoinsRequired, err := k.PostDispatch(ctx, originMailboxId, requiredHookAddress, metadata, hypMsg, maxFee)
	if err != nil {
		return util.HexAddress{}, err
	}

	if postDispatchHookId.IsZeroAddress() {
		defaultHookAddress, err := util.DecodeHexAddress(mailbox.DefaultHook)
		if err != nil {
			return util.HexAddress{}, err
		}
		postDispatchHookId = defaultHookAddress
	}

	remainingCoins, neg := maxFee.SafeSub(chargedCoinsRequired...)
	if neg {
		return util.HexAddress{}, fmt.Errorf("remaining coins cannot be negative")
	}

	chargedCoinsDefault, err := k.PostDispatch(ctx, originMailboxId, postDispatchHookId, metadata, hypMsg, remainingCoins)
	if err != nil {
		return util.HexAddress{}, err
	}

	chargedCoins := chargedCoinsRequired.Add(chargedCoinsDefault...)

	if chargedCoins.IsAnyGT(maxFee) {
		return util.HexAddress{}, fmt.Errorf("maxFee exceeded %s > %s", chargedCoins.String(), maxFee.String())
	}

	return hypMsg.Id(), nil
}
