package types

import (
	"context"
	"encoding/binary"
	"slices"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/cosmos/gogoproto/proto"
	"github.com/ethereum/go-ethereum/crypto"
)

type HyperlaneInterchainSecurityModule interface {
	proto.Message

	ModuleType() uint8
	GetId() (util.HexAddress, error)
	Verify(ctx context.Context, metadata []byte, message util.HyperlaneMessage) (bool, error)
}

var (
	IsmsKey             = []byte{SubModuleId, 0}
	StorageLocationsKey = []byte{SubModuleId, 2}
)

const (
	SubModuleName       = "ism"
	SubModuleId   uint8 = 1
)

// Constants defined by the Hyperlane spec
const (
	INTERCHAIN_SECURITY_MODULE_TYPE_UNUSED uint8 = iota
	INTERCHAIN_SECURITY_MODULE_TYPE_ROUTING
	INTERCHAIN_SECURITY_MODULE_TYPE_AGGREGATION
	INTERCHAIN_SECURITY_MODULE_TYPE_LEGACY_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TYPE_MERKLE_ROOT_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TYPE_NULL // used with relayer carrying no metadata
	INTERCHAIN_SECURITY_MODULE_TYPE_CCIP_READ
	INTERCHAIN_SECURITY_MODULE_TYPE_ARB_L2_TO_L1
	INTERCHAIN_SECURITY_MODULE_TYPE_WEIGHTED_MERKLE_ROOT_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TYPE_WEIGHTED_MESSAGE_ID_MULTISIG
	INTERCHAIN_SECURITY_MODULE_TYPE_OP_L2_TO_L1
)

func GetAnnouncementDigest(storageLocation string, domainId uint32, mailbox []byte) [32]byte {
	var domainHashBytes []byte

	domainIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(domainIdBytes, domainId)

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
