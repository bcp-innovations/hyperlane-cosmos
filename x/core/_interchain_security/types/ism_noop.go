package types

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ HyperlaneInterchainSecurityModule = &NoopISM{}

func (m *NoopISM) GetId() uint64 {
	return m.Id
}

func (m *NoopISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED
}

func (m *NoopISM) Verify(_ sdk.Context, _ []byte, _ util.HyperlaneMessage) (bool, error) {
	return true, nil
}
