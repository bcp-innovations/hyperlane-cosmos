package keeper

import (
	"context"
	"fmt"
	"slices"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	mailboxkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	enabledTokens []int32

	// state management

	Params          collections.Item[types.Params]
	Schema          collections.Schema
	HypTokens       collections.Map[[]byte, types.HypToken]
	HypTokensCount  collections.Sequence
	EnrolledRouters collections.Map[collections.Pair[[]byte, uint32], types.RemoteRouter]

	bankKeeper    types.BankKeeper
	mailboxKeeper *mailboxkeeper.Keeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService storetypes.KVStoreService,
	authority string,
	bankKeeper types.BankKeeper,
	mailboxKeeper *mailboxkeeper.Keeper,
	enabledTokens []int32,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:             cdc,
		addressCodec:    addressCodec,
		authority:       authority,
		enabledTokens:   enabledTokens,
		HypTokens:       collections.NewMap(sb, types.HypTokenKey, "hyptokens", collections.BytesKey, codec.CollValue[types.HypToken](cdc)),
		Params:          collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		HypTokensCount:  collections.NewSequence(sb, types.HypTokensCountKey, "hyptokens_count"),
		EnrolledRouters: collections.NewMap(sb, types.EnrolledRoutersKey, "enrolled_routers", collections.PairKeyCodec(collections.BytesKey, collections.Uint32Key), codec.CollValue[types.RemoteRouter](cdc)),
		bankKeeper:      bankKeeper,
		mailboxKeeper:   mailboxKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k *Keeper) Handle(ctx context.Context, mailboxId util.HexAddress, origin uint32, sender util.HexAddress, message coreTypes.HyperlaneMessage) error {
	goCtx := sdk.UnwrapSDKContext(ctx)

	// Return nil when the module is not the recipient of the Handle call
	token, err := k.HypTokens.Get(ctx, message.Recipient.Bytes())
	if err != nil {
		return nil
	}

	payload, err := types.ParseWarpPayload(message.Body)
	if err != nil {
		return err
	}

	if util.HexAddress(token.OriginMailbox) != mailboxId {
		return fmt.Errorf("invalid origin mailbox address")
	}

	remoteRouter, err := k.EnrolledRouters.Get(ctx, collections.Join(message.Recipient.Bytes(), origin))
	if err != nil {
		return fmt.Errorf("no enrolled router found for origin %d", origin)
	}

	if sender.String() != remoteRouter.ReceiverContract {
		return fmt.Errorf("invalid receiver contract")
	}

	// Check token type
	err = nil
	if token.TokenType == types.HYP_TOKEN_TYPE_COLLATERAL {
		if !slices.Contains(k.enabledTokens, int32(types.HYP_TOKEN_TYPE_COLLATERAL)) {
			return fmt.Errorf("module disabled collateral tokens")
		}
		err = k.RemoteReceiveCollateral(goCtx, token, payload)
	} else if token.TokenType == types.HYP_TOKEN_TYPE_SYNTHETIC {
		if !slices.Contains(k.enabledTokens, int32(types.HYP_TOKEN_TYPE_SYNTHETIC)) {
			return fmt.Errorf("module disabled synthetic tokens")
		}
		err = k.RemoteReceiveSynthetic(goCtx, token, payload)
	} else {
		panic("inconsistent store")
	}

	return err
}
