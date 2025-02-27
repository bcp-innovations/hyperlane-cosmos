package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) CreateMerkleTreeHook(ctx context.Context, msg *types.MsgCreateMerkleTreeHook) (*types.MsgCreateMerkleTreeHookResponse, error) {
	tree := util.NewTree(util.ZeroHashes, 0)

	nextId, err := ms.k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE)
	if err != nil {
		return nil, err
	}
	merkleTreeHook := types.MerkleTreeHook{
		InternalId: nextId.GetInternalId(),
		Id:         nextId.String(),
		Owner:      msg.Owner,
		Tree:       types.ProtoFromTree(tree),
	}

	err = ms.k.merkleTreeHooks.Set(ctx, merkleTreeHook.InternalId, merkleTreeHook)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMerkleTreeHookResponse{
		Id: nextId.String(),
	}, nil
}
