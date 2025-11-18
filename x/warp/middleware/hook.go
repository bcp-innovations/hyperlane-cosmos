package middleware

import (
	"context"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// HandleFn is the signature of the Handle function of an Hyperlane application.
type HandleFn = func(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error

type HandleHookI interface {
	PreHandleHook(ctx context.Context, mailboxID util.HexAddress, message util.HyperlaneMessage) error
	PostHandleHook(ctx context.Context, mailboxID util.HexAddress, message util.HyperlaneMessage) error
}

type HandleHook struct {
	// preHandleHook is the behavior expected from a type that hook BEFORE executing
	// the Hyperlane application Handle method.
	preHandleFn HandleFn

	// postHandleHook is the behavior expected from a type that hook AFTER executing
	// the Hyperlane application Handle method.
	postHandleFn HandleFn
}

func (h *HandleHook) PreHandleHook(ctx context.Context, mailboxID util.HexAddress, message util.HyperlaneMessage) error {
	if h.preHandleFn != nil {
		return h.preHandleFn(ctx, mailboxID, message)
	}

	return nil
}

func (h *HandleHook) PostHandleHook(ctx context.Context, mailboxID util.HexAddress, message util.HyperlaneMessage) error {
	if h.postHandleFn != nil {
		return h.postHandleFn(ctx, mailboxID, message)
	}

	return nil
}

func NewHandleHook(opts ...HandleHookOpt) *HandleHook {
	h := &HandleHook{}

	for _, opt := range opts {
		opt(h)
	}
	return h
}

type HandleHookOpt func(*HandleHook)

func WithPreHandleHookFn(fn HandleFn) HandleHookOpt {
	return func(h *HandleHook) {
		h.preHandleFn = fn
	}
}

func WithPostHandleFn(fn HandleFn) HandleHookOpt {
	return func(h *HandleHook) {
		h.postHandleFn = fn
	}
}
