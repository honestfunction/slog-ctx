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

func TestSimpleLayers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := bytes.NewBuffer([]byte{})
	reader := bufio.NewReader(buf)

	handler := slogctx.Setup(slog.NewJSONHandler(buf, nil))
	slog.SetDefault(slog.New(handler))

	commonKey := "category"

	rootCtx := slogctx.With(ctx, commonKey, "root")

	slog.InfoCtx(rootCtx, "from root")

	line, _, _ := reader.ReadLine()
	jsonKeyValue := map[string]string{}
	json.Unmarshal(line, &jsonKeyValue)
	if jsonKeyValue[commonKey] != "root" {
		t.Error("not the expected value: root", jsonKeyValue[commonKey])
	}

	func(ctxIn context.Context) {
		ctx := slogctx.With(ctxIn, commonKey, "fun1")
		ctx = slogctx.With(ctx, "fun1", "fun1")
		slog.InfoCtx(ctx, "from fun1")

		line, _, _ := reader.ReadLine()
		jsonKeyValue := map[string]string{}
		json.Unmarshal(line, &jsonKeyValue)
		if jsonKeyValue[commonKey] != "fun1" {
			t.Error("not the expected value: fun1", jsonKeyValue[commonKey])
		}
		if jsonKeyValue["fun1"] != "fun1" {
			t.Error("not the expected value: fun1", jsonKeyValue["fun1"])
		}

	}(rootCtx)

	func(ctxIn context.Context) {
		ctx := slogctx.With(ctxIn, commonKey, "fun2")
		ctx = slogctx.With(ctx, "fun2", "fun2")
		slog.InfoCtx(ctx, "from fun2")

		line, _, _ := reader.ReadLine()
		jsonKeyValue := map[string]string{}
		json.Unmarshal(line, &jsonKeyValue)
		if jsonKeyValue[commonKey] != "fun2" {
			t.Error("not the expected value: fun2", jsonKeyValue[commonKey])
		}
		if jsonKeyValue["fun2"] != "fun2" {
			t.Error("not the expected value: fun2", jsonKeyValue["fun2"])
		}

		if jsonKeyValue["fun1"] != "" {
			t.Error("not the expected value: {}", jsonKeyValue["fun1"])
		}

	}(rootCtx)
}
