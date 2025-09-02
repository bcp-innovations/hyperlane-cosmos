package types

const (
	ModuleName       = "hyperlane"
	ModuleId   uint8 = 0
)

var (
	MailboxesKey          = []byte{ModuleId, 1}
	MailboxesSequenceKey  = []byte{ModuleId, 2}
	MessagesKey           = []byte{ModuleId, 3}
	IsmRouterKey          = []byte{ModuleId, 4}
	PostDispatchRouterKey = []byte{ModuleId, 5}
	AppRouterKey          = []byte{ModuleId, 6}
	MessagesMapKey        = []byte{ModuleId, 7}

	// Leave 0 in case we add params in the future.
)
