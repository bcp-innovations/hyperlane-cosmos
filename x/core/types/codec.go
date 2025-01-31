package types

import (
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	"cosmossdk.io/core/registry"
	"cosmossdk.io/core/transaction"
)

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry registry.InterfaceRegistrar) {
	registry.RegisterImplementations((*transaction.Msg)(nil),
		&MsgUpdateParams{},
		&MsgCreateMailbox{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
