# slog-ctx
The slog-ctx is a library helps to handle the log with context using slog

## Quick Start
 ```go
import (
	"context"
    "os"

	slogctx "github.com/honestfunction/slog-ctx"
	"golang.org/x/exp/slog"
)

func main() {

	handler := slogctx.Setup(slog.NewJSONHandler(os.Stdout, nil), slogctx.Handler())
	slog.SetDefault(slog.New(handler))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = slogctx.With(ctx, "requestID", "3456789527")

	call(ctx)
}

func call(ctx context.Context) {
	slog.InfoCtx(ctx, "do something")
}

```

```text
{"time":"2023-06-17T14:04:51.886745+08:00","level":"INFO","msg":"do something","requestID":"3456789527"}
```