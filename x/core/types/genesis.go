package types

// NewGenesisState creates a new genesis state with default values.
func NewGenesisState() *GenesisState {
	return &GenesisState{
		IsmGenesis:           nil, // TODO add
		PostDispatchGenesis:  nil, // TODO add
		Mailboxes:            nil,
		Messages:             nil,
		IsmSequence:          0,
		PostDispatchSequence: 0,
		AppSequence:          0,
	}
}

// Validate performs basic genesis state validation returning an error upon any
func (gs *GenesisState) Validate() error {
	// TODO validate

	return nil
}
