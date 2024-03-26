package rlog

import (
	"context"
	"github.com/v587labs/robin/rlog/adapter"
	"github.com/v587labs/robin/rlog/adapter/log15"
)

type logContextKey struct {
}

var root adapter.Logger = log15.New()

func New(ctx ...interface{}) adapter.Logger {
	return root.New(ctx...)
}

// Log a message at the given level with context key/value pairs
func Debug(msg string, ctx ...interface{}) {
	root.Debug(msg, ctx...)
}
func Info(msg string, ctx ...interface{}) {
	root.Info(msg, ctx...)
}
func Warn(msg string, ctx ...interface{}) {
	root.Warn(msg, ctx...)
}
func Error(msg string, ctx ...interface{}) {
	root.Error(msg, ctx...)
}
func Crit(msg string, ctx ...interface{}) {
	root.Crit(msg, ctx...)
}

func L(ctx context.Context) adapter.Logger {
	value := ctx.Value(logContextKey{})
	if value == nil {
		return root
	}
	l, ok := value.(adapter.Logger)
	if !ok {
		return root
	}
	return l
}

func WithContext(ctx context.Context, logger adapter.Logger) context.Context {
	return context.WithValue(ctx, logContextKey{}, logger)
}

func With(ctx context.Context, keyvals ...interface{}) (context.Context, adapter.Logger) {
	logger := L(ctx).New(keyvals...)
	ctx = WithContext(ctx, logger)
	return ctx, logger
}
