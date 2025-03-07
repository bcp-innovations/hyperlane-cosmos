package keeper

import (
	"context"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

// Keeper is used to handle all core implementations of ISMs and implements the
// Go HyperlaneInterchainSecurityModule. Every core ISM does not require any outside keeper
// and can therefore all be handled by the same handler. If an ISM needs to access state
// in the future, one needs to provide another IsmHandler which holds the keeper and can access state.
type Keeper struct {
	// isms is a map from internal ISM ID to the ISM.
	isms collections.Map[uint64, types.HyperlaneInterchainSecurityModule]
	// storageLocations is a map from a key composed of a (mailbox ID, validator
	// address, and storage location index) to the storage location. A storage location
	// is a string that that describes where a validator persists their signatures.
	storageLocations collections.Map[collections.Triple[uint64, []byte, uint64], string]
	schema           collections.Schema

	coreKeeper types.CoreKeeper
}

func NewKeeper(cdc codec.BinaryCodec, storeService storetypes.KVStoreService) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		isms:             collections.NewMap(sb, types.IsmsKey, "isms", collections.Uint64Key, codec.CollInterfaceValue[types.HyperlaneInterchainSecurityModule](cdc)),
		storageLocations: collections.NewMap(sb, types.StorageLocationsKey, "storage_locations", collections.TripleKeyCodec(collections.Uint64Key, collections.BytesKey, collections.Uint64Key), collections.StringValue),
		coreKeeper:       nil,
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema

	return k
}

func (k *Keeper) SetCoreKeeper(coreKeeper types.CoreKeeper) {
	if k.coreKeeper != nil {
		panic("core keeper already set")
	}

	k.coreKeeper = coreKeeper

	// set the router from the core keeper
	router := coreKeeper.IsmRouter()
	// add default modules
	router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TYPE_UNUSED, k)
	router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TYPE_MERKLE_ROOT_MULTISIG, k)
	router.RegisterModule(types.INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG, k)
}

// Verify checks if the metadata has signed the message correctly.
func (k *Keeper) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	ism, err := k.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return false, err
	}

	return ism.Verify(ctx, metadata, message)
}

// Exists checks if the given ISM id does exist.
func (k *Keeper) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	return k.isms.Has(ctx, ismId.GetInternalId())
}
