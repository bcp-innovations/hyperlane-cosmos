package keeper

import (
	"context"
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"fmt"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	mailboxkeeper "github.com/bcp-innovations/hyperlane-cosmos/x/core/keeper"
	coreTypes "github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated message.
	// typically, this should be the x/gov module account.
	authority string

	// state management

	Params         collections.Item[types.Params]
	Schema         collections.Schema
	HypTokens      collections.Map[[]byte, types.HypToken]
	HypTokensCount collections.Sequence

	bankKeeper    bankkeeper.Keeper
	mailboxKeeper *mailboxkeeper.Keeper
}

// NewKeeper creates a new Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	addressCodec address.Codec,
	storeService storetypes.KVStoreService,
	authority string,
	bankKeeper bankkeeper.Keeper,
	mailboxKeeper *mailboxkeeper.Keeper,
) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}
	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:            cdc,
		addressCodec:   addressCodec,
		authority:      authority,
		HypTokens:      collections.NewMap(sb, types.HypTokenKey, "hyptokens", collections.BytesKey, codec.CollValue[types.HypToken](cdc)),
		Params:         collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		HypTokensCount: collections.NewSequence(sb, types.HypTokensCountKey, "hyptokens_count"),
		bankKeeper:     bankKeeper,
		mailboxKeeper:  mailboxKeeper,
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

	token, err := k.HypTokens.Get(ctx, message.Recipient.Bytes())
	if err != nil {
		return err
	}

	payload, err := types.ParseWarpPayload(message.Body)
	if err != nil {
		return err
	}

	if util.HexAddress(token.OriginMailbox) != mailboxId {
		return fmt.Errorf("invalid origin mailbox address")
	}

	if origin != token.ReceiverDomain {
		return fmt.Errorf("invalid origin denom")
	}

	if sender != util.HexAddress(token.ReceiverContract) {
		return fmt.Errorf("invalid receiver contract")
	}

	// Check token type
	err = nil
	if token.TokenType == types.HYP_TOKEN_COLLATERAL {
		err = k.RemoteReceiveCollateral(goCtx, token, payload)
	} else if token.TokenType == types.HYP_TOKEN_SYNTHETIC {
		err = k.RemoteReceiveSynthetic(goCtx, token, payload)
	} else {
		panic("inconsistent store")
	}

	return err
}
