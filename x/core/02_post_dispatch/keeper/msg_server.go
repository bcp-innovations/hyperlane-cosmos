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

func (k *Keeper) CreateAggregationHook(ctx context.Context, msg *types.MsgCreateAggregationHook) (util.HexAddress, error) {
	if err := k.validateAggregationHooks(ctx, msg.Hooks); err != nil {
		return util.HexAddress{}, err
	}

	nextId, err := k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_AGGREGATION)
	if err != nil {
		return util.HexAddress{}, err
	}

	aggregationHook := types.AggregationHook{
		Id:    nextId,
		Owner: msg.Owner,
		Hooks: msg.Hooks,
	}

	err = k.aggregationHooks.Set(ctx, aggregationHook.Id.GetInternalId(), aggregationHook)
	if err != nil {
		return util.HexAddress{}, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateAggregationHook{
		AggregationHookId: aggregationHook.Id,
		Owner:             aggregationHook.Owner,
		Hooks:             aggregationHook.Hooks,
	})

	return nextId, nil
}

func (ms msgServer) CreateAggregationHook(ctx context.Context, msg *types.MsgCreateAggregationHook) (*types.MsgCreateAggregationHookResponse, error) {
	nextId, err := ms.k.CreateAggregationHook(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateAggregationHookResponse{
		Id: nextId,
	}, nil
}

func (k *Keeper) validateAggregationHooks(ctx context.Context, hooks []util.HexAddress) error {
	if len(hooks) == 0 {
		return errors.Wrap(types.ErrInvalidAggregationHook, "hooks cannot be empty")
	}

	seen := make(map[util.HexAddress]struct{}, len(hooks))
	for _, hookId := range hooks {
		if hookId.IsZeroAddress() {
			return errors.Wrap(types.ErrInvalidAggregationHook, "hook id cannot be zero")
		}
		if hookId.GetType() == uint32(types.POST_DISPATCH_HOOK_TYPE_AGGREGATION) {
			return errors.Wrapf(types.ErrInvalidAggregationHook, "nested aggregation hook is not allowed: %s", hookId.String())
		}
		if _, ok := seen[hookId]; ok {
			return errors.Wrapf(types.ErrInvalidAggregationHook, "duplicate hook id: %s", hookId.String())
		}
		seen[hookId] = struct{}{}

		handler, err := k.coreKeeper.PostDispatchRouter().GetModule(hookId)
		if err != nil {
			return errors.Wrapf(types.ErrHookDoesNotExistOrIsNotRegistered, "%s", hookId.String())
		}
		exists, err := (*handler).Exists(ctx, hookId)
		if err != nil || !exists {
			return errors.Wrapf(types.ErrHookDoesNotExistOrIsNotRegistered, "%s", hookId.String())
		}
	}

	return nil
}
