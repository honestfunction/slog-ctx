package tests

import (
	"context"
	"os"
	"testing"

	slogctx "github.com/honestfunction/slog-ctx"
	"golang.org/x/exp/slog"
)

func TestReadme(t *testing.T) {

	handler := slogctx.Setup(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(slog.New(handler))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = slogctx.With(ctx, "requestID", "3456789527")

	call(ctx)
}

func call(ctx context.Context) {
	slog.InfoCtx(ctx, "do something")
}
