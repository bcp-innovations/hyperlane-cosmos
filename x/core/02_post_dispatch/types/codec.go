package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateIgp{}, "hyperlane/post_dispatch/MsgCreateIgp", nil)
	cdc.RegisterConcrete(&MsgSetIgpOwner{}, "hyperlane/post_dispatch/MsgSetIgpOwner", nil)
	cdc.RegisterConcrete(&MsgSetDestinationGasConfig{}, "hyperlane/post_dispatch/MsgSetDestinationGasConfig", nil)
	cdc.RegisterConcrete(&MsgPayForGas{}, "hyperlane/post_dispatch/MsgPayForGas", nil)
	cdc.RegisterConcrete(&MsgClaim{}, "hyperlane/post_dispatch/MsgClaim", nil)
	cdc.RegisterConcrete(&MsgCreateMerkleTreeHook{}, "hyperlane/post_dispatch/MsgCreateMerkleTreeHook", nil)
	cdc.RegisterConcrete(&MsgCreateNoopHook{}, "hyperlane/post_dispatch/MsgCreateNoopHook", nil)
}

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

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
