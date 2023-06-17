package tests

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"testing"

	slogctx "github.com/honestfunction/slog-ctx"
	"golang.org/x/exp/slog"
)

func TestDupKey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := bytes.NewBuffer([]byte{})

	handler := slogctx.Setup(slog.NewJSONHandler(buf, nil))
	slog.SetDefault(slog.New(handler))

	lastVal := "3333"

	ctx = slogctx.With(ctx, "key", "1")
	ctx = slogctx.With(ctx, "key", "2")
	ctx = slogctx.With(ctx, "key", lastVal)

	slog.InfoCtx(ctx, "dup key")

	line, _, _ := bufio.NewReader(buf).ReadLine()
	jsonKeyValue := map[string]string{}
	json.Unmarshal(line, &jsonKeyValue)

	if jsonKeyValue["key"] != lastVal {
		t.Error("not the last value:", lastVal, jsonKeyValue["key"])
	}
}
