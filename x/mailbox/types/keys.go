package types

import "cosmossdk.io/collections"

const ModuleName = "mailbox"

var (
	ParamsKey                     = collections.NewPrefix(0)
	MailboxesKey                  = collections.NewPrefix(1)
	MailboxesSequenceKey          = collections.NewPrefix(2)
	MessagesKey                   = collections.NewPrefix(3)
	ReceiverIsmKey                = collections.NewPrefix(4)
	ValidatorsKey                 = collections.NewPrefix(5)
	ValidatorsSequencesKey        = collections.NewPrefix(6)
	IgpKey                        = collections.NewPrefix(7)
	IgpDestinationGasConfigMapKey = collections.NewPrefix(8)
	IgpSequenceKey                = collections.NewPrefix(9)
)

const (
	TokenExchangeRateScale uint64 = 1
)
