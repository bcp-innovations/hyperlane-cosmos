package types

import (
	"context"
	"slices"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &RoutingISM{}

// GetId implements HyperlaneInterchainSecurityModule.
func (m *RoutingISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

// ModuleType implements HyperlaneInterchainSecurityModule.
func (m *RoutingISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_ROUTING
}

// ModuleType implements HyperlaneInterchainSecurityModule.
func (m *RoutingISM) Verify(ctx context.Context, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	// This method will never be called in the routing ISM struct
	// Routing happens on the Handler level in `routing_ism_handler.go`
	return false, errors.Wrapf(ErrUnexpectedError, "Verify should not be called on RoutingISM")
}

func (m *RoutingISM) GetIsm(domainId uint32) (*util.HexAddress, bool) {
	for _, route := range m.Routes {
		if route.Domain == domainId {
			return &route.Ism, true
		}
	}
	return nil, false
}

func (m *RoutingISM) RemoveDomain(domainId uint32) bool {
	for i, route := range m.Routes {
		if route.Domain == domainId {
			m.Routes = slices.Delete(m.Routes, i, i+1)
			return true
		}
	}
	return false
}

// SetDomain adds/sets a domain to the routing ISM.
func (m *RoutingISM) SetDomain(newRoute Route) {
	for _, route := range m.Routes {
		if newRoute.Domain == route.Domain {
			route.Ism = newRoute.Ism
			return
		}
	}
	m.Routes = append(m.Routes, &newRoute)
}
