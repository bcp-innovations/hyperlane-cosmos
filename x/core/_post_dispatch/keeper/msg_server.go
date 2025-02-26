package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/_post_dispatch/types"
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

	merkleTreeHook := types.MerkleTreeHook{
		Id:    ms.k.idFactory.GenerateNewId(ctx),
		Owner: msg.Owner,
		Tree:  types.ProtoFromTree(tree),
	}

	err := ms.k.merkleTreeHooks.Set(ctx, merkleTreeHook.Id, merkleTreeHook)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMerkleTreeHookResponse{
		// TODO return id
	}, nil
}
