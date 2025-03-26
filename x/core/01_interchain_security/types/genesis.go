package types

import "github.com/cosmos/cosmos-sdk/codec/types"

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Isms:                      []*types.Any{},
		ValidatorStorageLocations: []GenesisValidatorStorageLocationWrapper{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	return nil
}
