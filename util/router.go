package util

import (
	"context"
	"encoding/binary"
	"fmt"

	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InterchainSecurityModule interface {
	Exists(ctx context.Context, ismId HexAddress) (bool, error)
	Verify(ctx context.Context, ismId HexAddress, metadata []byte, message HyperlaneMessage) (bool, error)
}

type PostDispatchModule interface {
	Exists(ctx context.Context, hookId HexAddress) (bool, error)
	PostDispatch(ctx context.Context, hookId HexAddress, metadata any, message HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
}

type Router[T any] struct {
	modules  map[uint8]T
	sequence collections.Sequence
}

// TODO: custom address prefix
func NewRouter[T any](keyPrefix []byte, builder *collections.SchemaBuilder) *Router[T] {
	sequence := collections.NewSequence(builder, keyPrefix, "router_sequence")

	return &Router[T]{
		modules:  make(map[uint8]T),
		sequence: sequence,
	}
}

func (r *Router[T]) RegisterModule(moduleId uint8, module T) {
	if _, ok := r.modules[moduleId]; ok {
		panic("module already registered")
	}
	r.modules[moduleId] = module
}

func (r *Router[T]) GetModule(ctx context.Context, id HexAddress) (*T, error) {
	// the first byte of the id are the module id
	moduleId := id[0]
	module, ok := r.modules[moduleId]
	if !ok {
		return nil, fmt.Errorf("module with id %d not found", moduleId)
	}
	return &module, nil
}

// GetNextSequence returns the next sequence number and maps it to the given module id.
//
// The is is a 32 byte array encoded as follows:
// - 0:1 bytes are the module id
// - 24:32 bytes are the sequence number
// - the rest of the bytes are reserved for future use
func (r *Router[T]) GetNextSequence(ctx context.Context, moduleId uint8) (HexAddress, error) {
	next, err := r.sequence.Next(ctx)
	if err != nil {
		return HexAddress{}, err
	}

	if _, ok := r.modules[moduleId]; !ok {
		return HexAddress{}, fmt.Errorf("module with id %d not found", moduleId)
	}

	id := [32]byte{}
	id[0] = moduleId
	binary.BigEndian.PutUint64(id[24:32], next)

	return id, nil
}
