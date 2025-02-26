package types

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
	IsmRouterKey                  = []byte{ModuleId, 13}
	PostDispatchRouterKey         = []byte{ModuleId, 14}
)
