package util

import (
	"context"
	"encoding/binary"
	"fmt"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/collections"
)

const (
	nameLength       = 20
	moduleTypeLength = 4
	sequenceLength   = 8
)

type InterchainSecurityModule interface {
	Exists(ctx context.Context, ismId HexAddress) (bool, error)
	Verify(ctx context.Context, ismId HexAddress, metadata []byte, message HyperlaneMessage) (bool, error)
}

type PostDispatchModule interface {
	Exists(ctx context.Context, hookId HexAddress) (bool, error)
	PostDispatch(ctx context.Context, mailboxId, hookId HexAddress, metadata []byte, message HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error)
	QuoteDispatch(ctx context.Context, mailboxId, hookId HexAddress, metadata []byte, message HyperlaneMessage) (sdk.Coins, error)
	HookType() uint8
	SupportsMetadata(metadata []byte) bool
}

type HyperlaneApp interface {
	Exists(ctx context.Context, recipient HexAddress) (bool, error)
	Handle(ctx context.Context, mailboxId HexAddress, message HyperlaneMessage) error
	ReceiverIsmId(ctx context.Context, recipient HexAddress) (HexAddress, error)
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

func (r *Router[T]) GetModule(ctx context.Context, id HexAddress) (*T, error) {
	// the first byte of the id are the module id
	moduleId := id.GetType()
	module, ok := r.modules[moduleId]
	if !ok {
		return nil, fmt.Errorf("id %s not found", id.String())
	}
	return &module, nil
}

// GetNextSequence returns the next sequence number and maps it to the given module id.
//
// The is is a 32 byte array encoded as follows:
// - 0:20 bytes are the module name
// - 20:24 bytes are the module id
// - 24:32 bytes are the sequence number
// - the rest of the bytes are reserved for future use
func (r *Router[T]) GetNextSequence(ctx context.Context, moduleId uint8) (HexAddress, error) {
	id := uint32(moduleId)
	next, err := r.sequence.Next(ctx)
	if err != nil {
		return HexAddress{}, fmt.Errorf("failed to get next sequence: %w", err)
	}

	if _, ok := r.modules[id]; !ok {
		return HexAddress{}, fmt.Errorf("module with id %d not found", moduleId)
	}

	address := GenerateHexAddress(r.name, id, next)
	return address, nil
}

/*
TODO: refactor hex addresses, settle on one type
SPEC: HexAddress

The HexAddress mimics an evm-compatible address for a smart contract.
Due to the nature of cosmos, addresses must be created differently.

Requirements:
- HexAddresses must be unique across all cosmos modules interacting with Hyperlane

Structure
- The HexAddress has 32 bytes and is used for external communication
- For internal usage and storage an uint64 is totally sufficient

HexAddress: <module-specifier (20 byte)> <type (4 byte)> <internal-id (8 byte)>

The struct provides functions to encode and decode the information stored within the address.

To ensure global uniqueness, the HexAddressFactory should be used. It is initialized once per Keeper
and keeps global track of all registered module specifiers.

*/
// Hex Address Factory
func GenerateHexAddress(moduleSpecifier [20]byte, internalType uint32, internalId uint64) HexAddress {
	internalTypeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(internalTypeBytes, internalType)

	internalIdBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(internalIdBytes, internalId)

	return HexAddress(slices.Concat(moduleSpecifier[:], internalTypeBytes, internalIdBytes))
}
