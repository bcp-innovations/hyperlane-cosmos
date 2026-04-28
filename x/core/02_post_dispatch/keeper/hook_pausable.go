package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PausableHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = PausableHookHandler{}

// Exists returns whether a pausable hook exists for the given hook ID.
func (p PausableHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := p.k.pausableHooks.Has(ctx, hookId.GetInternalId())
	if err != nil {
		return false, err
	}
	return has, nil
}

// HookType returns the Hyperlane post-dispatch hook type for pausable hooks.
func (p PausableHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_PAUSABLE
}

// QuoteDispatch returns the required fee for the hook.
// Pausing does not charge fees, so it always returns zero coins.
func (p PausableHookHandler) QuoteDispatch(_ context.Context, _, _ util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage) (sdk.Coins, error) {
	return sdk.NewCoins(), nil
}

// PostDispatch rejects dispatch when the hook is paused.
func (p PausableHookHandler) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage, _ sdk.Coins) (sdk.Coins, error) {
	hook, err := p.k.pausableHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return sdk.NewCoins(), err
	}

	if hook.MailboxId != mailboxId {
		return sdk.NewCoins(), errors.Wrapf(types.ErrSenderIsNotDesignatedMailbox, "required mailbox id: %s, sender mailbox id: %s", hook.MailboxId, mailboxId.String())
	}

	if hook.Paused {
		return sdk.NewCoins(), errors.Wrapf(types.ErrPausableHookPaused, "hook: %s", hookId.String())
	}

	return sdk.NewCoins(), nil
}
