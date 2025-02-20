package types

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ HyperlaneInterchainSecurityModule = &MerkleRootMultiSigISM{}

func (m *MerkleRootMultiSigISM) GetId() uint64 {
	return m.Id
}

func (m *MerkleRootMultiSigISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG
}

func (m *MerkleRootMultiSigISM) Verify(ctx sdk.Context, metadata any, message util.HyperlaneMessage) (bool, error) {
	// TODO implement me

	panic("implement me")
}
