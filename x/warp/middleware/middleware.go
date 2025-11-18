package middleware

import (
	"context"
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

// AppMiddleware defines the behavior an Hyperlane application middleware must implement.
type AppMiddleware interface {
	util.HyperlaneApp
}

var _ AppMiddleware = (*Middleware)(nil)

// Middleware is the default concrete type used to implement an Hyperlane
// application middleware.
type Middleware struct {
	util.HyperlaneApp
	hook HandleHookI
}

func NewMiddleware(inner util.HyperlaneApp, hook HandleHookI) *Middleware {
	return &Middleware{
		HyperlaneApp: inner,
		hook:         hook,
	}
}

func (m *Middleware) Exists(ctx context.Context, recipient util.HexAddress) (bool, error) {
	return m.HyperlaneApp.Exists(ctx, recipient)
}

func (m *Middleware) Handle(ctx context.Context, mailboxID util.HexAddress, message util.HyperlaneMessage) error {

	err := m.hook.PreHandleHook(ctx, mailboxID, message)
	if err != nil {
		return fmt.Errorf("failed to execute pre handle hook: %w", err)
	}

	err = m.HyperlaneApp.Handle(ctx, mailboxID, message)
	if err != nil {
		return err
	}

	err = m.hook.PostHandleHook(ctx, mailboxID, message)
	if err != nil {
		return fmt.Errorf("failed to execute post handle hook: %w", err)
	}
	return nil
}

func (m *Middleware) ReceiverIsmId(ctx context.Context, recipient util.HexAddress) (*util.HexAddress, error) {
	return m.HyperlaneApp.ReceiverIsmId(ctx, recipient)
}
