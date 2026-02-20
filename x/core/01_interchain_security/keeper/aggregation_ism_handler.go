package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
)

// AggregationISMHandler
// The AggregationISM is a special ISM that aggregates verification from multiple ISMs.
// It requires a threshold number of ISMs to pass verification for the aggregation to succeed.
type AggregationISMHandler struct {
	keeper *Keeper // The ism keeper
}

// Verify implements HyperlaneInterchainSecurityModule
// Verifies the message against all configured ISMs and returns true if threshold is met.
func (m *AggregationISMHandler) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	ism, err := m.keeper.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return false, err
	}

	// check if the ism is an aggregation ism
	aggregationIsm, ok := ism.(*types.AggregationISM)
	if !ok {
		return false, errors.Wrapf(types.ErrInvalidISMType, "ISM %s is not an aggregation ISM", ismId.String())
	}

	// count how many ISMs pass verification
	passCount := uint32(0)
	var verificationErrors []error

	for _, moduleId := range aggregationIsm.Modules {
		// call the top level Verify method on the core module
		// this method will then recursively invoke the Verify method on all the sub ISMs
		verified, err := m.keeper.coreKeeper.Verify(ctx, moduleId, metadata, message)
		if err != nil {
			// Track errors for diagnostic purposes
			verificationErrors = append(verificationErrors, errors.Wrapf(err, "ISM %s verification error", moduleId.String()))
			continue
		}

		if verified {
			passCount++
			// Early exit if we've already met the threshold
			if passCount >= aggregationIsm.Threshold {
				return true, nil
			}
		}
	}

	// Build detailed error message with verification failures
	errMsg := errors.Wrapf(types.ErrInsufficientVerifications,
		"insufficient verifications: %d/%d (threshold: %d)",
		passCount, len(aggregationIsm.Modules), aggregationIsm.Threshold)

	// Append individual verification errors for debugging
	for _, verErr := range verificationErrors {
		errMsg = errors.Wrap(errMsg, verErr.Error())
	}

	return false, errMsg
}

func (m *AggregationISMHandler) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	return m.keeper.isms.Has(ctx, ismId.GetInternalId())
}
