package slogctx

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"testing"

	"golang.org/x/exp/slog"
)

func TestContextLog(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	buf := bytes.NewBuffer([]byte{})

	handler := Setup(slog.NewJSONHandler(buf, nil))
	slog.SetDefault(slog.New(handler))

	input := []entry{
		{randStr(5), randStr(10)},
		{randStr(5), randStr(10)},
		{randStr(5), randStr(10)},
		{randStr(5), randStr(10)},
		{randStr(5), randStr(10)},
	}

	for _, i := range input {
		ctx = With(ctx, i.key, i.val)
	}

	testContextByCallee(t, ctx, input, buf)
}

func testContextByCallee(t *testing.T, ctx context.Context, input []entry, buf *bytes.Buffer) {
	slog.InfoCtx(ctx, "hello slog-ctx")
	line, _, _ := bufio.NewReader(buf).ReadLine()

	jsonKeyValue := map[string]string{}
	json.Unmarshal(line, &jsonKeyValue)

	t.Log(jsonKeyValue)

	for _, i := range input {
		if jsonKeyValue[i.key] != i.val {
			t.Error("Not equal:", i.key, jsonKeyValue[i.key], i.val)
		}
	}
}

func randStr(size int) string {
	if size <= 0 {
		return ""
	}
	charset := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
