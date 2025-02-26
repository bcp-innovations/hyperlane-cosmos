package types

import (
	"context"

	//"fmt"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/cosmos/gogoproto/proto"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ MailboxHooks = MultiMailboxHooks{}

type MultiMailboxHooks []MailboxHooks

func NewMultiMailboxHooks(hooks ...MailboxHooks) MultiMailboxHooks {
	return hooks
}

func (h MultiMailboxHooks) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
	for i := range h {
		if err := h[i].Handle(ctx, mailboxId, message); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiMailboxHooks) ReceiverIsmId(ctx context.Context, recipient util.HexAddress) (util.HexAddress, error) {
	var receiverIsm util.HexAddress
	for i := range h {
		ismId, err := h[i].ReceiverIsmId(ctx, recipient)
		if err != nil {
			return util.HexAddress{}, nil
		}
		if !ismId.IsZeroAddress() {
			if receiverIsm.IsZeroAddress() {
				receiverIsm = ismId
			} else {
				return util.HexAddress{}, ErrMultipleReceiverIsm
			}
		}
	}

	if !receiverIsm.IsZeroAddress() {
		return receiverIsm, nil
	} else {
		return util.HexAddress{}, ErrNoReceiverISM
	}
}
