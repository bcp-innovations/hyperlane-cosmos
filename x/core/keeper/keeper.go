package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	hooks types.MailboxHooks
	// state management
	Mailboxes collections.Map[[]byte, types.Mailbox]
	Messages  collections.KeySet[[]byte]
	// Key is the Receiver address (util.HexAddress) and value is the util.HexAddress of the ISM
	ReceiverIsmMapping collections.Map[[]byte, []byte]
	MailboxesSequence  collections.Sequence
	Validators         collections.Map[[]byte, types.Validator]
	ValidatorsSequence collections.Sequence
	// IGP
	Igp                        collections.Map[[]byte, types.Igp]
	IgpDestinationGasConfigMap collections.Map[collections.Pair[[]byte, uint32], types.DestinationGasConfig]
	IgpSequence                collections.Sequence
	// ISM
	Isms         collections.Map[[]byte, types.Ism]
	IsmsSequence collections.Sequence

	Params collections.Item[types.Params]
	Schema collections.Schema

	bankKeeper types.BankKeeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, bankKeeper types.BankKeeper) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:                        cdc,
		addressCodec:               addressCodec,
		authority:                  authority,
		Mailboxes:                  collections.NewMap(sb, types.MailboxesKey, "mailboxes", collections.BytesKey, codec.CollValue[types.Mailbox](cdc)),
		Messages:                   collections.NewKeySet(sb, types.MessagesKey, "messages", collections.BytesKey),
		Params:                     collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		hooks:                      nil,
		MailboxesSequence:          collections.NewSequence(sb, types.MailboxesSequenceKey, "mailboxes_sequence"),
		Validators:                 collections.NewMap(sb, types.ValidatorsKey, "validators", collections.BytesKey, codec.CollValue[types.Validator](cdc)),
		ValidatorsSequence:         collections.NewSequence(sb, types.ValidatorsSequencesKey, "validators_sequence"),
		Igp:                        collections.NewMap(sb, types.IgpKey, "igp", collections.BytesKey, codec.CollValue[types.Igp](cdc)),
		IgpDestinationGasConfigMap: collections.NewMap(sb, types.IgpDestinationGasConfigMapKey, "igp_destination_gas_config_map", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), codec.CollValue[types.DestinationGasConfig](cdc)),
		IgpSequence:                collections.NewSequence(sb, types.IgpSequenceKey, "igp_sequence"),
		Isms:                       collections.NewMap(sb, types.IsmsKey, "isms", collections.BytesKey, codec.CollValue[types.Ism](cdc)),
		IsmsSequence:               collections.NewSequence(sb, types.IsmsSequencesKey, "isms_sequence"),
		ReceiverIsmMapping:         collections.NewMap(sb, types.ReceiverIsmKey, "receiver_ism", collections.BytesKey, collections.BytesValue),
		bankKeeper:                 bankKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k *Keeper) RegisterReceiverIsm(ctx context.Context, receiver util.HexAddress, mailboxId util.HexAddress, ismId string) error {
	prefixedIsmId, err := util.DecodeHexAddress(ismId)
	if err != nil || ismId == "" {
		// Use DefaultISM if no ISM is specified
		mailbox, err := k.Mailboxes.Get(ctx, mailboxId.Bytes())
		if err != nil {
			return err
		}

		prefixedIsmId, err = util.DecodeHexAddress(mailbox.DefaultIsm)
		if err != nil {
			return err
		}
	}

	exists, err := k.IsmIdExists(ctx, prefixedIsmId)
	if err != nil || !exists {
		return err
	}

	has, err := k.ReceiverIsmMapping.Has(ctx, receiver.Bytes())
	if err != nil || has {
		return err
	}

	return k.ReceiverIsmMapping.Set(ctx, receiver.Bytes(), prefixedIsmId.Bytes())
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

func (k Keeper) IgpIdExists(ctx context.Context, igpId util.HexAddress) (bool, error) {
	igp, err := k.Igp.Has(ctx, igpId.Bytes())
	if err != nil {
		return false, err
	}
	return igp, nil
}
