package log

import (
	"context"
	"github.com/v587labs/robin/log/adapter/log15"
)

type logContextKey struct {
}

var root Logger = log15.New()

func New(ctx ...interface{}) Logger {
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

func L(ctx context.Context) Logger {
	value := ctx.Value(logContextKey{})
	if value == nil {
		return root
	}
	l, ok := value.(Logger)
	if !ok {
		return root
	}
	return l
}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, logContextKey{}, logger)
}

func With(ctx context.Context, keyvals ...interface{}) (context.Context, Logger) {
	logger := L(ctx).New(keyvals...)
	ctx = WithContext(ctx, logger)
	return ctx, logger
}