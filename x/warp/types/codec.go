package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateCollateralToken{}, "hyperlane/warp/MsgCreateCollateralToken", nil)
	cdc.RegisterConcrete(&MsgCreateSyntheticToken{}, "hyperlane/warp/MsgCreateSyntheticToken", nil)
	cdc.RegisterConcrete(&MsgSetTokenResponse{}, "hyperlane/warp/MsgSetTokenResponse", nil)
	cdc.RegisterConcrete(&MsgEnrollRemoteRouter{}, "hyperlane/warp/MsgEnrollRemoteRouter", nil)
	cdc.RegisterConcrete(&MsgUnrollRemoteRouter{}, "hyperlane/warp/MsgUnrollRemoteRouter", nil)
	cdc.RegisterConcrete(&MsgRemoteTransfer{}, "hyperlane/warp/MsgRemoteTransfer", nil)
}

// RegisterInterfaces registers the interfaces types with the interface registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCollateralToken{},
		&MsgCreateSyntheticToken{},
		&MsgSetTokenResponse{},
		&MsgEnrollRemoteRouter{},
		&MsgUnrollRemoteRouter{},
		&MsgRemoteTransfer{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var Amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
