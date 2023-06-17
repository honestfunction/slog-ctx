package slogctx

import (
	"context"

	"golang.org/x/exp/slog"
)

func SetupHandler(handler slog.Handler, fns ...HandlerFunc) slog.Handler {
	if len(fns) == 0 {
		return handler
	}

	var h = handler
	for _, fn := range fns {
		h = &decoratedHandler{
			Handler:   h,
			handlerFn: fn,
		}
	}
	return h
}

type HandlerFunc func(ctx context.Context, r *slog.Record) error

type decoratedHandler struct {
	slog.Handler

	handlerFn HandlerFunc
}

func (h *decoratedHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.handlerFn == nil {
		return h.Handler.Handle(ctx, r)
	}
	if err := h.handlerFn(ctx, &r); err != nil {
		return err
	}
	return h.Handler.Handle(ctx, r)
}
