package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

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

func (k *Keeper) CreateMerkleTreeHook(ctx context.Context, msg *types.MsgCreateMerkleTreeHook) (util.HexAddress, error) {
	if exists, err := k.coreKeeper.MailboxIdExists(ctx, msg.MailboxId); !exists || err != nil {
		return util.HexAddress{}, errors.Wrapf(types.ErrMailboxDoesNotExist, "%s", msg.MailboxId)
	}

	nextId, err := k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE)
	if err != nil {
		return util.HexAddress{}, err
	}
	merkleTreeHook := types.MerkleTreeHook{
		Id:        nextId,
		MailboxId: msg.MailboxId,
		Owner:     msg.Owner,
		Tree:      types.ProtoFromTree(util.NewTree(util.ZeroHashes, 0)),
	}

	err = k.merkleTreeHooks.Set(ctx, merkleTreeHook.Id.GetInternalId(), merkleTreeHook)
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateMerkleTreeHook{
		MerkleTreeHookId: merkleTreeHook.Id,
		MailboxId:        merkleTreeHook.MailboxId,
		Owner:            merkleTreeHook.Owner,
	})

	return nextId, nil
}

func (ms msgServer) CreateMerkleTreeHook(ctx context.Context, msg *types.MsgCreateMerkleTreeHook) (*types.MsgCreateMerkleTreeHookResponse, error) {
	nextId, err := ms.k.CreateMerkleTreeHook(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateMerkleTreeHookResponse{
		Id: nextId,
	}, nil
}

func (k *Keeper) CreateNoopHook(ctx context.Context, msg *types.MsgCreateNoopHook) (util.HexAddress, error) {
	nextId, err := k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_UNUSED)
	if err != nil {
		return util.HexAddress{}, err
	}
	noopHook := types.NoopHook{
		Id:    nextId,
		Owner: msg.Owner,
	}

	err = k.noopHooks.Set(ctx, nextId.GetInternalId(), noopHook)
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateNoopHook{
		NoopHookId: noopHook.Id,
		Owner:      noopHook.Owner,
	})

	return nextId, nil
}

func (ms msgServer) CreateNoopHook(ctx context.Context, msg *types.MsgCreateNoopHook) (*types.MsgCreateNoopHookResponse, error) {
	nextId, err := ms.k.CreateNoopHook(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateNoopHookResponse{
		Id: nextId,
	}, nil
}
