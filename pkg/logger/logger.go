package logger

import (
	"context"
	"io"
	"log/slog"
	"maps"
	"os"
)

func SetupLog(path string, level slog.Level) {
	w := io.Writer(os.Stdout)

	if path != "" {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o666)
		if err != nil {
			panic(err)
		}

		w = io.MultiWriter(os.Stdout, file)
	}

	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		for k, v := range c.m {
			rec.Add(k, v)
		}
	}

	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)}
}

type logCtx struct {
	m map[string]any
}

type keyType int

const key = keyType(0)

func WithContext(ctx context.Context, k string, value any) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.m = maps.Clone(c.m)
		c.m[k] = value
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{map[string]any{k: value}})
}

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := logCtx{}

	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}

	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	if e, ok := err.(*errorWithLogCtx); ok {
		return context.WithValue(ctx, key, e.ctx)
	}

	return ctx
}
