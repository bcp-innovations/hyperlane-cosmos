package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateIgp{},
		&MsgSetIgpOwner{},
		&MsgSetDestinationGasConfig{},
		&MsgPayForGas{},
		&MsgClaim{},
		&MsgCreateMerkleTreeHook{},
		&MsgCreateNoopHook{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
