package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	ismkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/_interchain_security/keeper"
	postdispatchkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/keeper"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/query"
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
	// IGP
	Igp                        collections.Map[[]byte, types.Igp]
	IgpDestinationGasConfigMap collections.Map[collections.Pair[[]byte, uint32], types.DestinationGasConfig]
	IgpSequence                collections.Sequence

	Params collections.Item[types.Params]
	Schema collections.Schema

	bankKeeper types.BankKeeper

	IsmKeeper ismkeeper.Keeper
	ismHooks  types.InterchainSecurityHooks

	PostDispatchKeeper postdispatchkeeper.Keeper
	postDispatchHooks  types.PostDispatchHooks
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
		Igp:                        collections.NewMap(sb, types.IgpKey, "igp", collections.BytesKey, codec.CollValue[types.Igp](cdc)),
		IgpDestinationGasConfigMap: collections.NewMap(sb, types.IgpDestinationGasConfigMapKey, "igp_destination_gas_config_map", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), codec.CollValue[types.DestinationGasConfig](cdc)),
		IgpSequence:                collections.NewSequence(sb, types.IgpSequenceKey, "igp_sequence"),
		ReceiverIsmMapping:         collections.NewMap(sb, types.ReceiverIsmKey, "receiver_ism", collections.BytesKey, collections.BytesValue),
		bankKeeper:                 bankKeeper,

		// REFACTORED
		IsmKeeper:          ismkeeper.NewKeeper(cdc, storeService),
		PostDispatchKeeper: postdispatchkeeper.NewKeeper(cdc, storeService),
	}

	k.IsmKeeper.SetCoreKeeper(k)

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k *Keeper) RegisterReceiverIsm(ctx context.Context, receiver util.HexAddress, mailboxId util.HexAddress, ismId string) error {
	var prefixedIsmId util.HexAddress
	var err error

	// Use DefaultISM if no ISM is specified
	if ismId == "" {
		mailbox, err := k.Mailboxes.Get(ctx, mailboxId.Bytes())
		if err != nil {
			return err
		}

		prefixedIsmId, err = util.DecodeHexAddress(mailbox.DefaultIsm)
		if err != nil {
			return err
		}
	} else {
		prefixedIsmId, err = util.DecodeHexAddress(ismId)
		if err != nil {
			return fmt.Errorf("invalid ism id: %s", err)
		}
	}

	exists, err := k.IsmKeeper.IsmIdExists(ctx, prefixedIsmId)
	if err != nil || !exists {
		return fmt.Errorf("ism with id %s does not exist", prefixedIsmId)
	}

	return k.ReceiverIsmMapping.Set(ctx, receiver.Bytes(), prefixedIsmId.Bytes())
}

func (k *Keeper) PostDispatchHooks() types.PostDispatchHooks {
	if k.postDispatchHooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiPostDispatchHooks{}
	}

	return k.postDispatchHooks
}

func (k *Keeper) SetPostDispatchHooks(sh types.PostDispatchHooks) {
	if k.postDispatchHooks != nil {
		panic("cannot set mailbox hooks twice")
	}

	k.postDispatchHooks = sh
}

func (k *Keeper) IsmHooks() types.InterchainSecurityHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiInterchainSecurityHooks{}
	}

	return k.ismHooks
}

func (k *Keeper) SetIsmHooks(sh types.InterchainSecurityHooks) {
	if k.ismHooks != nil {
		panic("cannot set mailbox hooks twice")
	}

	k.ismHooks = sh
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

func (k Keeper) LocalDomain(ctx context.Context) (uint32, error) {
	params, err := k.Params.Get(ctx)
	return params.Domain, err
}

func (k Keeper) MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error) {
	mailbox, err := k.Mailboxes.Has(ctx, mailboxId.Bytes())
	if err != nil {
		return false, err
	}
	return mailbox, nil
}

// TODO out-source to utils
func GetPaginatedFromMap[T any, K any](ctx context.Context, collection collections.Map[K, T], pagination *query.PageRequest) ([]T, *query.PageResponse, error) {
	// Parse basic pagination
	if pagination == nil {
		pagination = &query.PageRequest{CountTotal: true}
	}

	offset := pagination.Offset
	key := pagination.Key
	limit := pagination.Limit
	reverse := pagination.Reverse

	if limit == 0 {
		limit = query.DefaultLimit
	}

	pageResponse := query.PageResponse{}

	// user has to use either offset or key, not both
	if offset > 0 && key != nil {
		return nil, nil, fmt.Errorf("invalid request, either offset or key is expected, got both")
	}

	ordering := collections.OrderDescending
	if reverse {
		ordering = collections.OrderAscending
	}

	// TODO: subject to change -> use it as key so we can jump to the offset directly
	it, err := collection.IterateRaw(ctx, key, nil, ordering)
	if err != nil {
		return nil, nil, err
	}

	defer it.Close()

	data := make([]T, 0, limit)
	keyValues, err := it.KeyValues()
	if err != nil {
		return nil, nil, err
	}
	length := uint64(len(keyValues))

	i := uint64(offset)
	for ; i < limit+offset && i < length; i++ {
		data = append(data, keyValues[i].Value)
	}

	if i < length {
		encodedKey := keyValues[i].Key
		codec := collection.KeyCodec()
		buffer := make([]byte, codec.Size(encodedKey))
		_, err := codec.Encode(buffer, encodedKey)
		if err != nil {
			return nil, nil, err
		}
		pageResponse.NextKey = buffer
	}

	return data, &pageResponse, nil
}
