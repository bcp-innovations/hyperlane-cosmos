package keeper

import (
	"fmt"

	"cosmossdk.io/math"

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
	received, err := k.Messages.Has(ctx, message.Id().Bytes())
	if err != nil {
		return err
	}
	if received {
		return fmt.Errorf("already received messsage with id %s", message.Id().String())
	}
	err = k.Messages.Set(ctx, message.Id().Bytes())
	if err != nil {
		return err
	}

	// TODO convert to hook for mailbox Client
	rawIsmAddress, err := k.ReceiverIsmMapping.Get(ctx, message.Recipient.Bytes())
	if err != nil {
		return fmt.Errorf("failed to get receiver ism address for recipient: %s", message.Recipient.String())
	}

	ismId := util.HexAddress(rawIsmAddress)

	// New logic
	verified, err := k.ismHooks.Verify(ctx, ismId, metadata, message)
	if err != nil {
		return err
	}
	if !verified {
		return fmt.Errorf("ism verification failed")
	}

	err = k.Hooks().Handle(ctx, mailboxId, message.Origin, message.Sender, message)
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

func (k Keeper) DispatchMessage2(
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
	metadata any,
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

	err = k.Messages.Set(ctx, hypMsg.Id().Bytes())
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

	requiredHookAddress, err := util.DecodeHexAddress(mailbox.RequiredHook)
	if err != nil {
		return util.HexAddress{}, err
	}
	remainingCoins, err := k.postDispatchHooks.PostDispatch(ctx, requiredHookAddress, metadata, hypMsg, maxFee)
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

	finalCoins, err := k.postDispatchHooks.PostDispatch(ctx, postDispatchHookId, metadata, hypMsg, remainingCoins)
	if err != nil {
		return util.HexAddress{}, err
	}

	chargedCoins := finalCoins.Add(remainingCoins...)

	if chargedCoins.IsAnyGT(maxFee) {
		return util.HexAddress{}, fmt.Errorf("maxFee exceeded %s > %s", chargedCoins.String(), maxFee.String())
	}

	return hypMsg.Id(), nil
}
