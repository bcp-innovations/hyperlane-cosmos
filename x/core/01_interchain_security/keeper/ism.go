package keeper

import (
	"context"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// IsmDefaultHandler is used to handle all current implementations of Isms and implements the
// Go HyperlaneInterchainSecurityModule. Every current ISM does not require any outside keeper
// and can therefore all be handled by the same handler. If an ISM needs to access state
// in the future, one needs to provide another IsmHandler which holds the keeper and can access state.
type IsmDefaultHandler struct {
	keeper *Keeper
}

func NewIsmHandler(keeper *Keeper) *IsmDefaultHandler {
	return &IsmDefaultHandler{
		keeper: keeper,
	}
}

// Verify checks if the metadata has signed the message correctly.
func (h *IsmDefaultHandler) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	ism, err := h.keeper.isms.Get(ctx, ismId.Bytes())
	if err != nil {
		return false, err
	}

	return ism.Verify(ctx, metadata, message)
}

// Exists checks if the given ISM id does exist.
func (h *IsmDefaultHandler) Exists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	return h.keeper.isms.Has(ctx, ismId.Bytes())
}
