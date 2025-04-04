package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
)

// RoutingISMHandler
// The RoutingISM is a special ISM that routes messages to other ISMs based on the origin of the message.
// It enables the use of multiple ISMs for a single mailbox, allowing for more flexible and modular ISM configuration.
type RoutingISMHandler struct {
	keeper *Keeper // The ism keeper
}

// Verify implements HyperlaneInterchainSecurityModule
// Delegates the verify call to the configured ISM in the RoutingISM based on the origin of the message.
func (m *RoutingISMHandler) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	ism, err := m.keeper.isms.Get(ctx, ismId.GetInternalId())
	if err != nil {
		return false, err
	}

	// check if the ism is a routing ism
	routingIsm, ok := ism.(*types.RoutingISM)
	if !ok {
		return false, errors.Wrapf(types.ErrInvalidISMType, "ISM %s is not a routing ISM", ismId.String())
	}

	// get the ism for the registered route
	routedIsm, exists := routingIsm.GetIsm(message.Origin)
	if !exists || routedIsm == nil {
		return false, errors.Wrapf(types.ErrNoRouteFound, "no route found for message %s", message.String())
	}

	// call the top level Verify method on the core module
	// this method will then recursivly invoke the Verify method on all the sub ISMs
	return m.keeper.coreKeeper.Verify(ctx, *routedIsm, metadata, message)
}

func (m *RoutingISMHandler) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	return m.keeper.isms.Has(ctx, ismId.GetInternalId())
}
