package keeper

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MerkleTreeHookHandler struct {
	k Keeper
}

var _ types.PostDispatchHookHandler = MerkleTreeHookHandler{}

func (i MerkleTreeHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE
}

func (i MerkleTreeHookHandler) SupportsMetadata(_ []byte) bool {
	return false
}

// TODO add mailbox id IMPORTANT: Double check if caller = mailboxId
func (i MerkleTreeHookHandler) PostDispatch(ctx sdk.Context, hookId uint64, rawMetadata []byte, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	merkleTreeHook, err := i.k.merkleTreeHooks.Get(ctx, hookId)
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
		// MailboxId: mailbox.Id, TODO mailbox Id?? What if multiple mailboxes insert into the same tree-hook
	})

	merkleTreeHook.Tree = types.ProtoFromTree(tree)

	if err := i.k.merkleTreeHooks.Set(ctx, hookId, merkleTreeHook); err != nil {
		return nil, err
	}

	return nil, nil
}
