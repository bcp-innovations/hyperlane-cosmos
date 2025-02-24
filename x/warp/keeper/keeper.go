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

	hexAddressFactory util.HexAddressFactory

	// state management

	Params collections.Item[types.Params]
	Schema collections.Schema
	// <tokenId> -> Token
	HypTokens      collections.Map[uint64, types.HypToken]
	HypTokensCount collections.Sequence
	// <tokenId> <domain> -> RemoteRouter
	EnrolledRouters collections.Map[collections.Pair[uint64, uint32], types.RemoteRouter]

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

	factory, err := util.NewHexAddressFactory(types.HEX_ADDRESS_CLASS_IDENTIFIER)
	if err != nil {
		panic(err)
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:               cdc,
		addressCodec:      addressCodec,
		authority:         authority,
		enabledTokens:     enabledTokens,
		hexAddressFactory: factory,
		HypTokens:         collections.NewMap(sb, types.HypTokenKey, "hyptokens", collections.Uint64Key, codec.CollValue[types.HypToken](cdc)),
		Params:            collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		HypTokensCount:    collections.NewSequence(sb, types.HypTokensCountKey, "hyptokens_count"),
		EnrolledRouters:   collections.NewMap(sb, types.EnrolledRoutersKey, "enrolled_routers", collections.PairKeyCodec(collections.Uint64Key, collections.Uint32Key), codec.CollValue[types.RemoteRouter](cdc)),
		bankKeeper:        bankKeeper,
		mailboxKeeper:     mailboxKeeper,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k *Keeper) GetAddressFromToken(token types.HypToken) util.HexAddress {
	return k.hexAddressFactory.GenerateId(uint32(token.TokenType), token.Id)
}

func (k *Keeper) ReceiverIsmId(ctx context.Context, recipient util.HexAddress) (util.HexAddress, error) {
	// Return nil when the module is not the recipient of the Handle call
	if !k.hexAddressFactory.IsClassMember(recipient) {
		return util.NewZeroAddress(), nil
	}

	token, err := k.HypTokens.Get(ctx, recipient.GetInternalId())
	if err != nil {
		return util.NewZeroAddress(), nil
	}

	hexAddress, err := util.DecodeHexAddress(token.IsmId)
	if err != nil {
		return util.NewZeroAddress(), err
	}

	return hexAddress, nil
}

func (k *Keeper) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
	// Return nil when the module is not the recipient of the Handle call
	if !k.hexAddressFactory.IsClassMember(message.Recipient) {
		return nil
	}

	goCtx := sdk.UnwrapSDKContext(ctx)

	// Return nil when the module is not the recipient of the Handle call
	token, err := k.HypTokens.Get(ctx, message.Recipient.GetInternalId())
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

	remoteRouter, err := k.EnrolledRouters.Get(ctx, collections.Join(message.Recipient.GetInternalId(), message.Origin))
	if err != nil {
		return fmt.Errorf("no enrolled router found for origin %d", message.Origin)
	}

	if message.Sender.String() != remoteRouter.ReceiverContract {
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
