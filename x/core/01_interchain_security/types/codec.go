package types

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateMessageIdMultisigIsm{},
		&MsgCreateMerkleRootMultisigIsm{},
		&MsgCreateNoopIsm{},
		&MsgAnnounceValidator{},
		&MsgCreateRoutingIsm{},
		&MsgSetRoutingIsmDomain{},
		&MsgRemoveRoutingIsmDomain{},
		&MsgUpdateRoutingIsmOwner{},
		&MsgCreateAggregationIsm{},
		&MsgSetAggregationIsmModules{},
		&MsgUpdateAggregationIsmOwner{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterInterface(
		"hyperlane.core.interchain_security.v1.HyperlaneInterchainSecurityModule",
		(*HyperlaneInterchainSecurityModule)(nil),
		&NoopISM{},
		&MessageIdMultisigISM{},
		&MerkleRootMultisigISM{},
		&RoutingISM{},
		&AggregationISM{},
	)
}
