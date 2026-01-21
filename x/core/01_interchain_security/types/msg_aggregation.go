package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic performs stateless validation of MsgCreateAggregationIsm
func (m *MsgCreateAggregationIsm) ValidateBasic() error {
	// Validate creator address
	if m.Creator == "" {
		return errors.Wrap(ErrInvalidOwner, "creator cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return errors.Wrap(ErrInvalidOwner, "invalid creator address")
	}

	// Validate aggregation configuration
	if err := ValidateAggregationISM(m.Modules, m.Threshold); err != nil {
		return err
	}

	return nil
}

// ValidateBasic performs stateless validation of MsgSetAggregationIsmModules
func (m *MsgSetAggregationIsmModules) ValidateBasic() error {
	// Validate owner address
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "owner cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return errors.Wrap(ErrInvalidOwner, "invalid owner address")
	}

	// Validate ISM ID is not zero
	if m.IsmId.IsZero() {
		return errors.Wrap(ErrInvalidISM, "ism_id cannot be zero")
	}

	// Validate aggregation configuration
	if err := ValidateAggregationISM(m.Modules, m.Threshold); err != nil {
		return err
	}

	return nil
}

// ValidateBasic performs stateless validation of MsgUpdateAggregationIsmOwner
func (m *MsgUpdateAggregationIsmOwner) ValidateBasic() error {
	// Validate owner address
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "owner cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return errors.Wrap(ErrInvalidOwner, "invalid owner address")
	}

	// Validate ISM ID is not zero
	if m.IsmId.IsZero() {
		return errors.Wrap(ErrInvalidISM, "ism_id cannot be zero")
	}

	// Validate new owner address format if provided
	if m.NewOwner != "" {
		if _, err := sdk.AccAddressFromBech32(m.NewOwner); err != nil {
			return errors.Wrap(ErrInvalidOwner, "invalid new_owner address")
		}
	}

	// Cannot both set new owner and renounce
	if m.RenounceOwnership && m.NewOwner != "" {
		return errors.Wrap(ErrInvalidOwner, "cannot set new owner and renounce ownership at the same time")
	}

	// Cannot set empty owner without renouncing
	if !m.RenounceOwnership && m.NewOwner == "" {
		return errors.Wrap(ErrInvalidOwner, "cannot set owner to empty address without renouncing ownership")
	}

	return nil
}
