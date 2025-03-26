package types

import (
	"fmt"
)

func NewGenesisState() *GenesisState {
	return &GenesisState{
		Igps:            []InterchainGasPaymaster{},
		IgpGasConfigs:   []GenesisDestinationGasConfigWrapper{},
		MerkleTreeHooks: []MerkleTreeHook{},
		NoopHooks:       []NoopHook{},
	}
}

func (gs *GenesisState) Validate() error {
	igpMap := make(map[uint64]struct{})
	for _, igp := range gs.Igps {
		if _, ok := igpMap[igp.Id.GetInternalId()]; ok {
			return fmt.Errorf("duplicate igp: %s", igp.Id)
		}
		igpMap[igp.Id.GetInternalId()] = struct{}{}
	}

	for _, config := range gs.IgpGasConfigs {
		if _, ok := igpMap[config.IgpId]; !ok {
			return fmt.Errorf("igp does not exist: %d", config.IgpId)
		}
	}

	return nil
}
