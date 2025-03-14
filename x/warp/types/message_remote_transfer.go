package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsTypes "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRemoteTransfer{}
)

func (msg *MsgRemoteTransfer) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoteTransfer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoteTransfer) Route() string {
	return ModuleName
}

func (msg *MsgRemoteTransfer) Type() string {
	return "hyperlane/warp/MsgRemoteTransfer"
}

func (msg *MsgRemoteTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrapf(errorsTypes.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
