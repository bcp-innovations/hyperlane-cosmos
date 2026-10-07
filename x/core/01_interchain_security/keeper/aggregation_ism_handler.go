package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
)

// AggregationISMHandler
// The AggregationISM is a special ISM that requires a threshold of its sub-ISMs to verify a message.
// It mirrors the semantics of the AbstractAggregationIsm on EVM.
type AggregationISMHandler struct {
	keeper *Keeper // The ism keeper
}

// Verify implements HyperlaneInterchainSecurityModule
// Every sub-ISM for which metadata is provided must verify the message, and exactly
// `threshold` sub-ISMs must be provided with metadata.
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

	subMetadata, err := types.ParseAggregationMetadata(metadata, len(aggregationIsm.Modules))
	if err != nil {
		return false, errors.Wrap(types.ErrInvalidAggregationMetadata, err.Error())
	}

	var verifiedCount uint32
	for i, module := range aggregationIsm.Modules {
		if subMetadata[i] == nil {
			continue
		}

		// call the top level Verify method on the core module
		// this method will then recursively invoke the Verify method on all the sub ISMs
		verified, err := m.keeper.coreKeeper.Verify(ctx, module, subMetadata[i], message)
		if err != nil {
			return false, errors.Wrapf(err, "sub-ISM %s at index %d", module.String(), i)
		}
		if !verified {
			return false, nil
		}
		verifiedCount++
	}

	if verifiedCount != aggregationIsm.Threshold {
		return false, errors.Wrapf(types.ErrInvalidAggregationMetadata, "expected metadata for %d modules, got %d", aggregationIsm.Threshold, verifiedCount)
	}

	return true, nil
}

func (m *AggregationISMHandler) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	return m.keeper.isms.Has(ctx, ismId.GetInternalId())
}
