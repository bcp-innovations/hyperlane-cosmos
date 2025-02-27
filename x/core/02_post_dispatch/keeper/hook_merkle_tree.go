package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MerkleTreeHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = MerkleTreeHookHandler{}

func (i MerkleTreeHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := i.k.merkleTreeHooks.Has(ctx, hookId.GetInternalId())
	if err != nil {
		return false, err
	}
	return has, nil
}

func (i MerkleTreeHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE
}

func (i MerkleTreeHookHandler) SupportsMetadata(_ []byte) bool {
	return false
}

// TODO add mailbox id IMPORTANT: Double check if caller = mailboxId
func (i MerkleTreeHookHandler) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, rawMetadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	merkleTreeHook, err := i.k.merkleTreeHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	tree, err := types.TreeFromProto(merkleTreeHook.Tree)
	if err != nil {
		return nil, err
	}

	count := tree.GetCount()

	if err = tree.Insert(message.Id()); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.InsertedIntoTree{
		MessageId: message.Id().String(),
		Index:     count,
		MailboxId: mailboxId.String(),
	})

	merkleTreeHook.Tree = types.ProtoFromTree(tree)

	if err := i.k.merkleTreeHooks.Set(ctx, hookId.GetInternalId(), merkleTreeHook); err != nil {
		return nil, err
	}

	return nil, nil
}
