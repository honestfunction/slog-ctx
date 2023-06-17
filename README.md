# slog-ctx

The slog-ctx library is designed to facilitate handling context using slog.

slog is a structured logging package that has been accepted and is planned to be included in the standard library starting from Go 1.21. slog introduces the Ctx family of APIs, such as DebugCtx, InfoCtx, WarnCtx, and ErrorCtx. However, the default slog handler currently does not implement the context processing functionality.

With slog-ctx, setting up the slog default handler and appending context values becomes straightforward. You can then start logging using the standard slog Ctx family of APIs.

Please note that since slog is not yet part of the Go language's standard library, make sure you have installed the slog package correctly before using slog-ctx.

## Quick Start
 ```go
import (
	"context"
    "os"

	slogctx "github.com/honestfunction/slog-ctx"
	"golang.org/x/exp/slog"
)

func main() {

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

```

```text
{"time":"2023-06-17T14:04:51.886745+08:00","level":"INFO","msg":"do something","requestID":"3456789527"}
```
