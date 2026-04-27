package types

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &AggregationISM{}

// GetId implements HyperlaneInterchainSecurityModule.
func (m *AggregationISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

// ModuleType implements HyperlaneInterchainSecurityModule.
func (m *AggregationISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION
}

// Verify implements HyperlaneInterchainSecurityModule, but should not be called on AggregationISM.
func (m *AggregationISM) Verify(ctx context.Context, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	// This method will never be called in the aggregation ISM struct
	// Aggregation happens on the Handler level in `aggregation_ism_handler.go`
	return false, errors.Wrapf(ErrUnexpectedError, "Verify should not be called on AggregationISM")
}

// ValidateAggregationISM validates the AggregationISM configuration.
func ValidateAggregationISM(modules []util.HexAddress, threshold uint32) error {
	if threshold == 0 {
		return errors.Wrapf(ErrInvalidThreshold, "threshold must be greater than 0")
	}

	if len(modules) == 0 {
		return errors.Wrapf(ErrInvalidISM, "modules list cannot be empty")
	}

	if threshold > uint32(len(modules)) {
		return errors.Wrapf(ErrInvalidThreshold, "threshold (%d) cannot exceed number of modules (%d)", threshold, len(modules))
	}

	// Check for duplicate modules
	seen := make(map[string]bool)
	for _, module := range modules {
		key := module.String()
		if seen[key] {
			return errors.Wrapf(ErrInvalidISM, "duplicate module found: %s", key)
		}
		seen[key] = true
	}

	return nil
}
