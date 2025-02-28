package types

import fmt "fmt"

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if gs.Params.Domain == 0 {
		return fmt.Errorf("local domain cannot be 0")
	}

	// TODO validate

	return nil
}
