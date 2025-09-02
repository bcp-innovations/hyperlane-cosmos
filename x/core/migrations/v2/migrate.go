package v2

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MigrateStore migrates all messages from KeySet to Map using `true` as value.
func MigrateStore(ctx sdk.Context, v1KeySet collections.KeySet[collections.Pair[uint64, []byte]], v2Map collections.Map[collections.Pair[uint64, []byte], bool]) error {
	it, err := v1KeySet.Iterate(ctx, nil)
	if err != nil {
		return err
	}
	defer it.Close()

	for ; it.Valid(); it.Next() {
		key, err := it.Key()
		if err != nil {
			return err
		}
		if err := v2Map.Set(ctx, key, true); err != nil {
			return err
		}
		// TODO: Decide if we want to clean up directly.
		//if err := v1KeySet.Remove(ctx, key); err != nil {
		//	return err
		//}
	}
	return nil
}
