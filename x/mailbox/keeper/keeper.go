package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	hooks     types.MailboxHooks
	ismKeeper types.IsmKeeper
	// state management
	Mailboxes         collections.Map[[]byte, types.Mailbox]
	Messages          collections.KeySet[collections.Pair[[]byte, []byte]]
	MailboxesSequence collections.Sequence
	Params            collections.Item[types.Params]
	Schema            collections.Schema
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, ismKeeper types.IsmKeeper) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:               cdc,
		addressCodec:      addressCodec,
		authority:         authority,
		Mailboxes:         collections.NewMap(sb, types.MailboxesKey, "mailboxes", collections.BytesKey, codec.CollValue[types.Mailbox](cdc)),
		Messages:          collections.NewKeySet(sb, types.MessagesKey, "messages", collections.PairKeyCodec(collections.BytesKey, collections.BytesKey)),
		Params:            collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		hooks:             nil,
		MailboxesSequence: collections.NewSequence(sb, types.MailboxesSequenceKey, "mailboxes_sequence"),
		ismKeeper:         ismKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k *Keeper) PostDispatchMerkleTree(ctx context.Context, messageId string, index uint32) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.InsertedIntoTree{
		MessageId: messageId,
		Index:     index,
	})
}

// Hooks gets the hooks for staking *Keeper {
func (k *Keeper) Hooks() types.MailboxHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiMailboxHooks{}
	}

	return k.hooks
}

// SetHooks sets the validator hooks.  In contrast to other receivers, this method must take a pointer due to nature
// of the hooks interface and SDK start up sequence.
func (k *Keeper) SetHooks(sh types.MailboxHooks) {
	if k.hooks != nil {
		panic("cannot set mailbox hooks twice")
	}

	k.hooks = sh
}
