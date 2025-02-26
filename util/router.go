package util

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InterchainSecurityModule interface {
	Verify(ctx sdk.Context, ismId uint64, metadata []byte, message HyperlaneMessage) (bool, error)
}

type PostDispatchModule interface {
	PostDispatch(ctx sdk.Context, hookId uint64, metadata any, message HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
}

type Router[T any] struct {
	modules map[uint16]T
	sealed  bool

	sequence collections.Sequence
	mapping  collections.Map[uint64, uint16]
}

func NewRouter[T any](keyPrefix []byte, builder *collections.SchemaBuilder) *Router[T] {
	sequenceKey := append(keyPrefix, 0)
	mappingKey := append(keyPrefix, 1)

	sequence := collections.NewSequence(builder, sequenceKey, "router_sequence")
	mapping := collections.NewMap(builder, mappingKey, "router", collections.Uint64Key, collections.Uint16Value)

	return &Router[T]{
		modules:  make(map[uint16]T),
		sealed:   false,
		sequence: sequence,
		mapping:  mapping,
	}
}

func (r *Router[T]) RegisterModule(moduleId uint8, module T) {
	if r.sealed {
		panic("cannot add module to sealed router")
	}

	if _, ok := r.modules[uint16(moduleId)]; ok {
		panic("module already registered")
	}
	r.modules[uint16(moduleId)] = module
}

func (r *Router[T]) GetModule(ctx context.Context, id uint64) (*T, error) {
	moduleId, err := r.mapping.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	module, ok := r.modules[moduleId]
	if !ok {
		return nil, fmt.Errorf("module with id %d not found", moduleId)
	}
	return &module, nil
}

func (r *Router[T]) Seal() {
	r.sealed = true
}

func (r *Router[T]) IsSealed() bool {
	return r.sealed
}

func (r *Router[T]) Exists(ctx context.Context, id uint64) (bool, error) {
	exists, err := r.mapping.Has(ctx, id)
	return exists, err
}

// GetNextSequence returns the next sequence number and maps it to the given module id.
func (r *Router[T]) GetNextSequence(ctx context.Context, moduleId uint8) (uint64, error) {
	next, err := r.sequence.Next(ctx)
	if err != nil {
		return 0, err
	}

	if _, ok := r.modules[uint16(moduleId)]; !ok {
		return 0, fmt.Errorf("module with id %d not found", moduleId)
	}

	err = r.mapping.Set(ctx, next, uint16(moduleId))
	if err != nil {
		return 0, err
	}

	return next, nil
}
