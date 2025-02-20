package types

import (
	"cosmossdk.io/collections"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type HyperlaneInterchainSecurityModule interface {
	proto.Message

	GetId() uint64
	ModuleType() uint8
	Verify(ctx sdk.Context, metadata any, message util.HyperlaneMessage) (bool, error)
}

var (
	IsmsKey         = collections.NewPrefix(1)
	IsmsSequenceKey = collections.NewPrefix(2)
)

const (
	SubModuleName = "ism"

	HEX_ADDRESS_CLASS_IDENTIFIER = "coreism"
)

const (
	INTERCHAIN_SECURITY_MODULE_TPYE_UNUSED uint8 = iota
	INTERCHAIN_SECURITY_MODULE_TPYE_ROUTING
	INTERCHAIN_SECURITY_MODULE_TPYE_AGGREGATION
	INTERCHAIN_SECURITY_MODULE_TPYE_LEGACY_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TPYE_MERKLE_ROOT_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TPYE_MESSAGE_ID_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TPYE_NULL // used with relayer carrying no metadata
	INTERCHAIN_SECURITY_MODULE_TPYE_CCIP_READ
	INTERCHAIN_SECURITY_MODULE_TPYE_ARB_L2_TO_L1
	INTERCHAIN_SECURITY_MODULE_TPYE_WEIGHTED_MERKLE_ROOT_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TPYE_WEIGHTED_MESSAGE_ID_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TPYE_OP_L2_TO_L1
)
