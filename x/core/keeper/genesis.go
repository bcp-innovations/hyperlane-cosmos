package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"cosmossdk.io/collections"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	if err := k.ismRouter.SetInternalSequence(ctx, data.IsmSequence); err != nil {
		return err
	}

	if err := k.postDispatchRouter.SetInternalSequence(ctx, data.PostDispatchSequence); err != nil {
		return err
	}

	if err := k.appRouter.SetInternalSequence(ctx, data.AppSequence); err != nil {
		return err
	}

	for _, m := range data.Mailboxes {
		if err := k.Mailboxes.Set(ctx, m.Id.GetInternalId(), m); err != nil {
			return err
		}
	}

	if err := k.MailboxesSequence.Set(ctx, uint64(len(data.Mailboxes))); err != nil {
		return err
	}

	for _, message := range data.Messages {
		if err := k.MessagesMap.Set(ctx, collections.Join(message.MailboxId, message.MessageId.Bytes()), true); err != nil {
			return err
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	ismSequence, err := k.ismRouter.GetInternalSequence(ctx)
	if err != nil {
		return nil, err
	}
	postDispatchSequence, err := k.postDispatchRouter.GetInternalSequence(ctx)
	if err != nil {
		return nil, err
	}
	appSequence, err := k.appRouter.GetInternalSequence(ctx)
	if err != nil {
		return nil, err
	}

	messages := make([]types.GenesisMailboxMessageWrapper, 0)
	err = k.MessagesMap.Walk(ctx, nil, func(key collections.Pair[uint64, []byte], _ bool) (stop bool, err error) {
		messages = append(messages, types.GenesisMailboxMessageWrapper{
			MailboxId: key.K1(),
			MessageId: util.HexAddress(key.K2()),
		})
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	mailboxes := make([]types.Mailbox, 0)
	err = k.Mailboxes.Walk(ctx, nil, func(key uint64, value types.Mailbox) (stop bool, err error) {
		mailboxes = append(mailboxes, value)
		return false, nil
	})
	if err != nil {
		return nil, err
	}

	return &types.GenesisState{
		Mailboxes: mailboxes,
		Messages:  messages,

		IsmSequence:          ismSequence,
		PostDispatchSequence: postDispatchSequence,
		AppSequence:          appSequence,
	}, nil
}
