package types

import (
	"encoding/binary"
	"slices"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/crypto"
)

type HyperlaneInterchainSecurityModule interface {
	proto.Message

	GetId() uint64
	ModuleType() uint8
	Verify(ctx sdk.Context, metadata []byte, message util.HyperlaneMessage) (bool, error)
}

var (
	IsmsKey             = []byte{SubModuleId, 0}
	IsmsSequenceKey     = []byte{SubModuleId, 1}
	StorageLocationsKey = []byte{SubModuleId, 2}
)

const (
	SubModuleName       = "ism"
	SubModuleId   uint8 = 1

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

func GetAnnouncementDigest(storageLocation string, domainId uint32, mailbox []byte) [32]byte {
	var domainHashBytes []byte

	domainIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(domainIdBytes, domainId)

	// TODO: Check if all of them are required
	domainHashBytes = slices.Concat(
		domainIdBytes,
		mailbox,
		[]byte("HYPERLANE_ANNOUNCEMENT"),
	)

	domainHash := crypto.Keccak256Hash(domainHashBytes)

	announcementDigestBytes := slices.Concat(
		domainHash.Bytes(),
		[]byte(storageLocation),
	)

	return crypto.Keccak256Hash(announcementDigestBytes)
}
