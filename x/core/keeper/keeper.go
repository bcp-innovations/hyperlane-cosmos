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
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	// first key is the mailbox ID, second key is the message ID
	Messages collections.KeySet[collections.Pair[[]byte, []byte]]
	// Key is the Receiver address (util.HexAddress) and value is the util.HexAddress of the ISM
	MailboxesSequence collections.Sequence

	Params collections.Item[types.Params]
	Schema collections.Schema

	bankKeeper types.BankKeeper

	IsmKeeper          ismkeeper.Keeper
	PostDispatchKeeper postdispatchkeeper.Keeper

	ismRouter          *util.Router[util.InterchainSecurityModule]
	postDispatchRouter *util.Router[util.PostDispatchModule]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, bankKeeper types.BankKeeper) Keeper {
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
		bankKeeper:        bankKeeper,

		// REFACTORED
		IsmKeeper:          ismkeeper.NewKeeper(cdc, storeService),
		PostDispatchKeeper: postdispatchkeeper.NewKeeper(cdc, storeService, bankKeeper),

		ismRouter:          util.NewRouter[util.InterchainSecurityModule](types.IsmRouterKey, "router_ism", sb),
		postDispatchRouter: util.NewRouter[util.PostDispatchModule](types.PostDispatchRouterKey, "router_post_dispatch", sb),
	}

	k.IsmKeeper.SetCoreKeeper(k)
	k.PostDispatchKeeper.SetCoreKeeper(k)

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) IsmRouter() *util.Router[util.InterchainSecurityModule] {
	return k.ismRouter
}

func (k *Keeper) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	handler, err := k.ismRouter.GetModule(ctx, ismId)
	if err != nil {
		return false, err
	}
	return (*handler).Verify(ctx, ismId, metadata, message)
}

func (k *Keeper) IsmExists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	handler, err := k.ismRouter.GetModule(ctx, ismId)
	if err != nil {
		return false, err
	}
	return (*handler).Exists(ctx, ismId)
}

func (k *Keeper) AssertIsmExists(ctx context.Context, id string) error {
	ismId, err := util.DecodeHexAddress(id)
	if err != nil {
		return fmt.Errorf("ism id %s is invalid: %s", id, err.Error())
	}
	ismExists, err := k.IsmExists(ctx, ismId)
	if err != nil || !ismExists {
		return fmt.Errorf("ism with id %s does not exist", ismId.String())
	}

	return nil
}

func (k Keeper) PostDispatchRouter() *util.Router[util.PostDispatchModule] {
	return k.postDispatchRouter
}

func (k *Keeper) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, metadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	handler, err := k.postDispatchRouter.GetModule(ctx, hookId)
	if err != nil {
		return sdk.NewCoins(), err
	}
	return (*handler).PostDispatch(ctx, mailboxId, hookId, metadata, message, maxFee)
}

func (k *Keeper) PostDispatchHookExists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	handler, err := k.postDispatchRouter.GetModule(ctx, hookId)
	if err != nil {
		return false, err
	}
	return (*handler).Exists(ctx, hookId)
}

func (k *Keeper) AssertPostDispatchHookExists(ctx context.Context, id string) error {
	hookId, err := util.DecodeHexAddress(id)
	if err != nil {
		return fmt.Errorf("hook id %s is invalid: %s", id, err.Error())
	}
	hookExists, err := k.PostDispatchHookExists(ctx, hookId)
	if err != nil || !hookExists {
		return fmt.Errorf("hook with id %s does not exist", hookId.String())
	}
	return nil
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
