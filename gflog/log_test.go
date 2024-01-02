package gflog

import (
	"context"
	"testing"
)

func TestContextValue(t *testing.T) {
	ctx := context.Background()
	ctx = ContextValue(ctx, "1", "2")
	ctx = ContextValue(ctx, "222", "2")
	Log().Info(ContextValue(ctx, "1111", "2"), "test", "start", "test", "asdf", "asdf")
	Log().Info(ContextValue(ctx, "qweqwe", "2"), "test", "start", "test", "asdf", "asdf")

}
