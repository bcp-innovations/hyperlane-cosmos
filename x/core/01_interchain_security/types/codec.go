package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateMessageIdMultisigIsm{}, "hyperlane/interchain_security/MsgCreateMessageIdMultisigIsm", nil)
	cdc.RegisterConcrete(&MsgCreateMerkleRootMultisigIsm{}, "hyperlane/interchain_security/MsgCreateMerkleRootMultisigIsm", nil)
	cdc.RegisterConcrete(&MsgCreateNoopIsm{}, "hyperlane/interchain_security/MsgCreateNoopIsm", nil)
	cdc.RegisterConcrete(&MsgAnnounceValidator{}, "hyperlane/interchain_security/MsgAnnounceValidator", nil)
}

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateMessageIdMultisigIsm{},
		&MsgCreateMerkleRootMultisigIsm{},
		&MsgCreateNoopIsm{},
		&MsgAnnounceValidator{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)

	registry.RegisterInterface(
		"hyperlane.core._interchain_security.HyperlaneInterchainSecurityModule",
		(*HyperlaneInterchainSecurityModule)(nil),
		&NoopISM{},
		&MessageIdMultisigISM{},
		&MerkleRootMultisigISM{},
	)
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
