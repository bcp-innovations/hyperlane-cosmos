package keeper

import (
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// special ism handler that is used for HyperlaneInterchainSecurityModule
// HyperlaneInterchainSecurityModule implements the Verify method themselves and don't need any outside keeper
type IsmHandler struct {
	keeper *Keeper
}

func NewIsmHandler(keeper *Keeper) *IsmHandler {
	return &IsmHandler{
		keeper: keeper,
	}
}

func (h *IsmHandler) Verify(ctx sdk.Context, ismId uint64, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	ism, err := h.keeper.GetIsm(ctx, ismId)
	if err != nil {
		return false, err
	}
	return ism.Verify(ctx, metadata, message)
}
