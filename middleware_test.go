package slogctx

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"testing"

	"golang.org/x/exp/slog"

	"github.com/samber/lo"
)

func TestMiddleWare(t *testing.T) {
	appendKey := lo.RandomString(5, lo.LettersCharset)
	appendVal := lo.RandomString(10, lo.AllCharset)

	buf := bytes.NewBuffer([]byte{})
	handler := Setup(slog.NewJSONHandler(buf, nil),
		func(ctx context.Context, r *slog.Record) error {
			r.AddAttrs(slog.String(appendKey, appendVal))
			return nil
		})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	slog.Info("hello")
	line, _, _ := bufio.NewReader(buf).ReadLine()

	jsonKeyValue := map[string]string{}
	json.Unmarshal(line, &jsonKeyValue)

	if jsonKeyValue[appendKey] != appendVal {
		t.Fail()
	}
}

func SourceHandleFunc(ctx context.Context, r *slog.Record) error {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	r.AddAttrs(slog.String("logger", fmt.Sprintf("%s:%d", f.Func.Name(), f.Line)))
	r.AddAttrs(slog.String("file", path.Base(f.File)))
	return nil
}
