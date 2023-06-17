package slogctx

import (
	"context"

	"golang.org/x/exp/slog"
)

type contextKey string

var defaultKey = "logging-context"

type entry struct {
	key string
	val any
}

var defaultSlogCtx = New(defaultKey)

func Default() *SlogContext {
	return defaultSlogCtx
}

func SetDefault(c *SlogContext) {
	defaultSlogCtx = c
}

func SetDefaultKey(k string) {
	defaultSlogCtx.SetKey(k)
}

func With(ctx context.Context, key string, val any) context.Context {
	return defaultSlogCtx.With(ctx, key, val)
}

func Handler() HandlerFunc {
	return defaultSlogCtx.Handler()
}

type SlogContext struct {
	ctxKey contextKey
}

func New(key string) *SlogContext {
	return &SlogContext{ctxKey: contextKey(key)}
}

func (c *SlogContext) With(ctx context.Context, key string, val any) context.Context {
	v := ctx.Value(c.ctxKey)
	if v == nil {
		return context.WithValue(ctx, c.ctxKey, []entry{{key, val}})
	}
	entries := v.([]entry)

	return context.WithValue(ctx, c.ctxKey, append(entries, entry{key, val}))
}

func (c *SlogContext) Handler() HandlerFunc {
	return func(ctx context.Context, r *slog.Record) error {
		val := ctx.Value(c.ctxKey)
		if val == nil {
			return nil
		}

		entries, ok := val.([]entry)
		if !ok {
			// maybe error
			return nil
		}

		for _, e := range entries {
			r.AddAttrs(slog.Any(e.key, e.val))
		}
		return nil
	}
}

func (c *SlogContext) SetKey(k string) {
	c.ctxKey = contextKey(k)
}
