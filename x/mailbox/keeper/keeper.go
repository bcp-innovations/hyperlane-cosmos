package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
	"github.com/KYVENetwork/hyperlane-cosmos/util"
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
	Mailboxes collections.Map[[]byte, types.Mailbox]
	Messages  collections.KeySet[collections.Pair[[]byte, []byte]]
	// Key is the Receiver address (util.HexAddress) and value is the util.HexAddress of the ISM
	ReceiverIsmMapping collections.Map[[]byte, []byte]
	MailboxesSequence  collections.Sequence
	Validators         collections.Map[[]byte, types.Validator]
	ValidatorsSequence collections.Sequence
	Params             collections.Item[types.Params]
	Schema             collections.Schema
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, ismKeeper types.IsmKeeper) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:                cdc,
		addressCodec:       addressCodec,
		authority:          authority,
		Mailboxes:          collections.NewMap(sb, types.MailboxesKey, "mailboxes", collections.BytesKey, codec.CollValue[types.Mailbox](cdc)),
		Messages:           collections.NewKeySet(sb, types.MessagesKey, "messages", collections.PairKeyCodec(collections.BytesKey, collections.BytesKey)),
		Params:             collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		hooks:              nil,
		MailboxesSequence:  collections.NewSequence(sb, types.MailboxesSequenceKey, "mailboxes_sequence"),
		Validators:         collections.NewMap(sb, types.ValidatorsKey, "validators", collections.BytesKey, codec.CollValue[types.Validator](cdc)),
		ValidatorsSequence: collections.NewSequence(sb, types.ValidatorsSequencesKey, "validators_sequence"),
		ismKeeper:          ismKeeper,
		ReceiverIsmMapping: collections.NewMap(sb, types.ReceiverIsmKey, "receiver_ism", collections.BytesKey, collections.BytesValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k *Keeper) RegisterReceiverIsm(ctx context.Context, receiver util.HexAddress, ismId util.HexAddress) error {
	exists, err := k.ismKeeper.IsmIdExists(ctx, ismId.String())
	if err != nil || !exists {
		return err
	}

	has, err := k.ReceiverIsmMapping.Has(ctx, receiver.Bytes())
	if err != nil || has {
		return err
	}

	return k.ReceiverIsmMapping.Set(ctx, receiver.Bytes(), ismId.Bytes())
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
