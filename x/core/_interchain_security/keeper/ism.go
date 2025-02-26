package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
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

func (h *IsmHandler) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	ism, err := h.keeper.GetIsm(ctx, ismId)
	if err != nil {
		return false, err
	}
	return ism.Verify(ctx, metadata, message)
}

func (h *IsmHandler) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	return h.keeper.isms.Has(ctx, ismId.Bytes())
}
