package types

import (
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &DomainRoutingISM{}

func (m *DomainRoutingISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

func (m *DomainRoutingISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_ROUTING
}

func (m *DomainRoutingISM) Verify(_ context.Context, _ []byte, _ util.HyperlaneMessage) (bool, error) {
	// TODO: Discuss this solution
	panic("shouldn't be called")
}

func (m *DomainRoutingISM) Route(message util.HyperlaneMessage) (util.HexAddress, error) {
	for _, d := range m.DomainIsmMapping {
		if d.DomainId == message.Origin {
			return d.IsmId, nil
		}
	}
	return util.HexAddress{}, fmt.Errorf("domain %v not supported", message.Origin)
}
