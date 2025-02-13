package types

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/math"
)

const ModuleName = "hyperlane"

var (
	ParamsKey                     = collections.NewPrefix(0)
	MailboxesKey                  = collections.NewPrefix(1)
	MailboxesSequenceKey          = collections.NewPrefix(2)
	MessagesKey                   = collections.NewPrefix(3)
	ReceiverIsmKey                = collections.NewPrefix(4)
	ValidatorsKey                 = collections.NewPrefix(5)
	ValidatorsSequencesKey        = collections.NewPrefix(6)
	StorageLocationsKey           = collections.NewPrefix(7)
	IgpKey                        = collections.NewPrefix(8)
	IgpDestinationGasConfigMapKey = collections.NewPrefix(9)
	IgpSequenceKey                = collections.NewPrefix(10)
	IsmsKey                       = collections.NewPrefix(11)
	IsmsSequencesKey              = collections.NewPrefix(12)
)

var TokenExchangeRateScale = math.NewInt(1e10)

const (
	UNUSED uint32 = iota
	ROUTING
	AGGREGATION
	LEGACY_MULTISIG
	MERKLE_ROOT_MULTISIG
	MESSAGE_ID_MULTISIG
	NULL // used with relayer carrying no metadata
	CCIP_READ
	ARB_L2_TO_L1
	WEIGHTED_MERKLE_ROOT_MULTISIG
	WEIGHTED_MESSAGE_ID_MULTISIG
	OP_L2_TO_L1
)
