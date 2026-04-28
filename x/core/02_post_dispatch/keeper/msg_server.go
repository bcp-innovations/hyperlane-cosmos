package keeper

import (
	"context"
	stderrors "errors"

	"cosmossdk.io/collections"
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

func (k *Keeper) CreatePausableHook(ctx context.Context, msg *types.MsgCreatePausableHook) (util.HexAddress, error) {
	if exists, err := k.coreKeeper.MailboxIdExists(ctx, msg.MailboxId); !exists || err != nil {
		return util.HexAddress{}, errors.Wrapf(types.ErrMailboxDoesNotExist, "%s", msg.MailboxId)
	}

	nextId, err := k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_PAUSABLE)
	if err != nil {
		return util.HexAddress{}, err
	}

	hook := types.PausableHook{
		Id:        nextId,
		Owner:     msg.Owner,
		MailboxId: msg.MailboxId,
		Paused:    false,
	}

	if err := k.pausableHooks.Set(ctx, hook.Id.GetInternalId(), hook); err != nil {
		return util.HexAddress{}, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreatePausableHook{
		PausableHookId: hook.Id,
		Owner:          hook.Owner,
		MailboxId:      hook.MailboxId,
		Paused:         hook.Paused,
	})

	return nextId, nil
}

func (ms msgServer) CreatePausableHook(ctx context.Context, msg *types.MsgCreatePausableHook) (*types.MsgCreatePausableHookResponse, error) {
	nextId, err := ms.k.CreatePausableHook(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreatePausableHookResponse{
		Id: nextId,
	}, nil
}

func (ms msgServer) SetPausableHookPaused(ctx context.Context, msg *types.MsgSetPausableHookPaused) (*types.MsgSetPausableHookPausedResponse, error) {
	if msg.HookId.IsZeroAddress() || msg.HookId.GetType() != uint32(types.POST_DISPATCH_HOOK_TYPE_PAUSABLE) {
		return nil, errors.Wrapf(types.ErrInvalidPausableHook, "%s", msg.HookId.String())
	}

	hook, err := ms.k.pausableHooks.Get(ctx, msg.HookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	if hook.Owner != msg.Owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "owner %s is not hook owner", msg.Owner)
	}

	hook.Paused = msg.Paused
	if err := ms.k.pausableHooks.Set(ctx, hook.Id.GetInternalId(), hook); err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetPausableHookPaused{
		PausableHookId: msg.HookId,
		Owner:          msg.Owner,
		Paused:         msg.Paused,
	})

	return &types.MsgSetPausableHookPausedResponse{}, nil
}

func (k *Keeper) CreateRateLimitedHook(ctx context.Context, msg *types.MsgCreateRateLimitedHook) (util.HexAddress, error) {
	if exists, err := k.coreKeeper.MailboxIdExists(ctx, msg.MailboxId); !exists || err != nil {
		return util.HexAddress{}, errors.Wrapf(types.ErrMailboxDoesNotExist, "%s", msg.MailboxId)
	}

	nextId, err := k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_RATE_LIMITED)
	if err != nil {
		return util.HexAddress{}, err
	}

	hook := types.RateLimitedHook{
		Id:        nextId,
		Owner:     msg.Owner,
		MailboxId: msg.MailboxId,
	}

	if err := k.rateLimitedHooks.Set(ctx, hook.Id.GetInternalId(), hook); err != nil {
		return util.HexAddress{}, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateRateLimitedHook{
		RateLimitedHookId: hook.Id,
		Owner:             hook.Owner,
		MailboxId:         hook.MailboxId,
	})

	return nextId, nil
}

func (ms msgServer) CreateRateLimitedHook(ctx context.Context, msg *types.MsgCreateRateLimitedHook) (*types.MsgCreateRateLimitedHookResponse, error) {
	nextId, err := ms.k.CreateRateLimitedHook(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateRateLimitedHookResponse{
		Id: nextId,
	}, nil
}

func (ms msgServer) SetRateLimit(ctx context.Context, msg *types.MsgSetRateLimit) (*types.MsgSetRateLimitResponse, error) {
	if msg.HookId.IsZeroAddress() || msg.HookId.GetType() != uint32(types.POST_DISPATCH_HOOK_TYPE_RATE_LIMITED) {
		return nil, errors.Wrapf(types.ErrInvalidRateLimitedHook, "%s", msg.HookId.String())
	}

	if msg.TokenId.IsZeroAddress() {
		return nil, errors.Wrap(types.ErrInvalidRateLimitedHook, "token id cannot be zero")
	}

	if msg.MaxCapacity.LT(types.RateLimitDuration) {
		return nil, errors.Wrapf(types.ErrRateLimitNotSet, "max capacity must be at least %s", types.RateLimitDuration.String())
	}

	hook, err := ms.k.rateLimitedHooks.Get(ctx, msg.HookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	if hook.Owner != msg.Owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "owner %s is not hook owner", msg.Owner)
	}

	refillRate := msg.MaxCapacity.Quo(types.RateLimitDuration)
	now := currentBlockUnix(ctx)

	tokenRateLimit := types.TokenRateLimit{
		HookId:      msg.HookId,
		TokenId:     msg.TokenId,
		MaxCapacity: msg.MaxCapacity,
		RefillRate:  refillRate,
		LastUpdated: now,
	}
	effectiveCapacity := tokenRateLimit.EffectiveCapacity()
	tokenRateLimit.FilledLevel = effectiveCapacity

	key := types.TokenRateLimitKey(msg.HookId, msg.TokenId)
	existing, err := ms.k.tokenRateLimits.Get(ctx, key)
	if err == nil {
		currentLevel := ms.k.CurrentRateLimitLevel(ctx, existing)
		if currentLevel.GT(effectiveCapacity) {
			currentLevel = effectiveCapacity
		}
		tokenRateLimit.FilledLevel = currentLevel
	} else if !stderrors.Is(err, collections.ErrNotFound) {
		return nil, err
	}

	if err := ms.k.tokenRateLimits.Set(ctx, key, tokenRateLimit); err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventSetRateLimit{
		RateLimitedHookId: msg.HookId,
		Owner:             msg.Owner,
		TokenId:           msg.TokenId,
		MaxCapacity:       msg.MaxCapacity,
		RefillRate:        refillRate,
	})

	return &types.MsgSetRateLimitResponse{}, nil
}

func (ms msgServer) RemoveRateLimit(ctx context.Context, msg *types.MsgRemoveRateLimit) (*types.MsgRemoveRateLimitResponse, error) {
	if msg.HookId.IsZeroAddress() || msg.HookId.GetType() != uint32(types.POST_DISPATCH_HOOK_TYPE_RATE_LIMITED) {
		return nil, errors.Wrapf(types.ErrInvalidRateLimitedHook, "%s", msg.HookId.String())
	}

	if msg.TokenId.IsZeroAddress() {
		return nil, errors.Wrap(types.ErrInvalidRateLimitedHook, "token id cannot be zero")
	}

	hook, err := ms.k.rateLimitedHooks.Get(ctx, msg.HookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	if hook.Owner != msg.Owner {
		return nil, errors.Wrapf(types.ErrUnauthorized, "owner %s is not hook owner", msg.Owner)
	}

	key := types.TokenRateLimitKey(msg.HookId, msg.TokenId)
	if has, err := ms.k.tokenRateLimits.Has(ctx, key); err != nil {
		return nil, err
	} else if has {
		if err := ms.k.tokenRateLimits.Remove(ctx, key); err != nil {
			return nil, err
		}
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventRemoveRateLimit{
		RateLimitedHookId: msg.HookId,
		Owner:             msg.Owner,
		TokenId:           msg.TokenId,
	})

	return &types.MsgRemoveRateLimitResponse{}, nil
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
