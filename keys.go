package mailbox

import "cosmossdk.io/collections"

const ModuleName = "mailbox"

var (
	ParamsKey  = collections.NewPrefix(0)
	CounterKey = collections.NewPrefix(1)
)
