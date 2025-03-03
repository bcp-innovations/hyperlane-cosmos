package keeper

import (
	"context"

	"cosmossdk.io/collections"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

// InitGenesis initializes the module state from a genesis state.
func (k *Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	if err := k.ismRouter.SetInternalSequence(ctx, data.IsmSequence); err != nil {
		panic(err)
	}

	if err := k.postDispatchRouter.SetInternalSequence(ctx, data.PostDispatchSequence); err != nil {
		panic(err)
	}

	if err := k.appRouter.SetInternalSequence(ctx, data.AppSequence); err != nil {
		panic(err)
	}

	for _, m := range data.Mailboxes {
		if err := k.Mailboxes.Set(ctx, m.Id.GetInternalId(), *m); err != nil {
			panic(err)
		}
	}

	if err := k.MailboxesSequence.Set(ctx, uint64(len(data.Mailboxes))); err != nil {
		panic(err)
	}

	for _, message := range data.Messages {
		if err := k.Messages.Set(ctx, collections.Join(message.MailboxId, message.MessageId)); err != nil {
			panic(err)
		}
	}

	return nil
}

// ExportGenesis exports the module state to a genesis state.
func (k *Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	ismSequence, err := k.ismRouter.GetInternalSequence(ctx)
	if err != nil {
		panic(err)
	}
	postDispatchSequence, err := k.postDispatchRouter.GetInternalSequence(ctx)
	if err != nil {
		panic(err)
	}
	appSequence, err := k.appRouter.GetInternalSequence(ctx)
	if err != nil {
		panic(err)
	}

	messages := make([]*types.MailboxMessage, 0)
	err = k.Messages.Walk(ctx, nil, func(key collections.Pair[uint64, []byte]) (stop bool, err error) {
		messages = append(messages, &types.MailboxMessage{
			MailboxId: key.K1(),
			MessageId: key.K2(),
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	mailboxes := make([]*types.Mailbox, 0)
	err = k.Mailboxes.Walk(ctx, nil, func(key uint64, value types.Mailbox) (stop bool, err error) {
		mailboxes = append(mailboxes, &value)
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	return &types.GenesisState{
		Messages:             messages,
		IsmSequence:          ismSequence,
		PostDispatchSequence: postDispatchSequence,
		AppSequence:          appSequence,
	}, nil
}
