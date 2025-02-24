package types

import (
	"cosmossdk.io/math"
)

const (
	ModuleName       = "hyperlane"
	ModuleId   uint8 = 0
)

var (
	ParamsKey                     = []byte{ModuleId, 0}
	MailboxesKey                  = []byte{ModuleId, 1}
	MailboxesSequenceKey          = []byte{ModuleId, 2}
	MessagesKey                   = []byte{ModuleId, 3}
	ReceiverIsmKey                = []byte{ModuleId, 4}
	ValidatorsKey                 = []byte{ModuleId, 5}
	ValidatorsSequencesKey        = []byte{ModuleId, 6}
	StorageLocationsKey           = []byte{ModuleId, 7}
	IgpKey                        = []byte{ModuleId, 8}
	IgpDestinationGasConfigMapKey = []byte{ModuleId, 9}
	IgpSequenceKey                = []byte{ModuleId, 10}
	IsmsKey                       = []byte{ModuleId, 11}
	IsmsSequencesKey              = []byte{ModuleId, 12}
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
