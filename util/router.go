package util

import (
	"context"
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/collections"
)

const (
	nameLength = 20
)

type InterchainSecurityModule interface {
	Exists(ctx context.Context, ismId HexAddress) (bool, error)
	Verify(ctx context.Context, ismId HexAddress, metadata []byte, message HyperlaneMessage) (bool, error)
}

type PostDispatchModule interface {
	Exists(ctx context.Context, hookId HexAddress) (bool, error)
	// PostDispatch returns the charged coins.
	PostDispatch(ctx context.Context, mailboxId, hookId HexAddress, metadata StandardHookMetadata, message HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
	QuoteDispatch(ctx context.Context, mailboxId, hookId HexAddress, metadata StandardHookMetadata, message HyperlaneMessage) (sdk.Coins, error)
	HookType() uint8
}

type HyperlaneApp interface {
	Exists(ctx context.Context, recipient HexAddress) (bool, error)
	Handle(ctx context.Context, mailboxId HexAddress, message HyperlaneMessage) error
	ReceiverIsmId(ctx context.Context, recipient HexAddress) (*HexAddress, error)
}

type Router[T any] struct {
	modules  map[uint32]T
	sequence collections.Sequence
	name     [20]byte
}

func NewRouter[T any](keyPrefix []byte, name string, builder *collections.SchemaBuilder) *Router[T] {
	nameBytes := []byte(name)
	if len(nameBytes) > nameLength {
		panic(fmt.Sprintf("router name '%s' is too long, must be at most %d bytes", name, nameLength))
	}

	sequence := collections.NewSequence(builder, keyPrefix, name)

	fixedName := [20]byte{}
	copy(fixedName[:], nameBytes)

	return &Router[T]{
		modules:  make(map[uint32]T),
		sequence: sequence,
		name:     fixedName,
	}
}

func (r *Router[T]) RegisterModule(moduleId uint8, module T) {
	id := uint32(moduleId)
	if _, ok := r.modules[id]; ok {
		panic("module already registered")
	}
	r.modules[id] = module
}

func (r *Router[T]) GetModule(id HexAddress) (*T, error) {
	// the first byte of the id are the module id
	moduleId := id.GetType()
	module, ok := r.modules[moduleId]
	if !ok {
		return nil, fmt.Errorf("id %s not found", id.String())
	}
	return &module, nil
}

func (r *Router[T]) GetModuleIds() (moduleIds []uint32) {
	for moduleId := range r.modules {
		moduleIds = append(moduleIds, moduleId)
	}

	slices.Sort(moduleIds)
	return
}

// GetNextSequence returns the next sequence number and maps it to the given module id
func (r *Router[T]) GetNextSequence(ctx context.Context, moduleId uint8) (HexAddress, error) {
	id := uint32(moduleId)
	next, err := r.sequence.Next(ctx)
	if err != nil {
		return HexAddress{}, fmt.Errorf("failed to get next sequence: %w", err)
	}

	if _, ok := r.modules[id]; !ok {
		return HexAddress{}, fmt.Errorf("module with id not found: %d", moduleId)
	}

	address := GenerateHexAddress(r.name, id, next)
	return address, nil
}

func (r *Router[T]) GetInternalSequence(ctx context.Context) (uint64, error) {
	return r.sequence.Peek(ctx)
}

func (r *Router[T]) SetInternalSequence(ctx context.Context, value uint64) error {
	return r.sequence.Set(ctx, value)
}
