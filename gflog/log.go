package gflog

import (
	"context"
	"slices"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
)

type _contextKey struct {
}

type Logger struct {
	*glog.Logger
}

const LogCtxKey = "vars"

func Log(name ...string) *glog.Logger {
	l := g.Log(name...)
	if !slices.Contains(l.GetCtxKeys(), LogCtxKey) {
		l.AppendCtxKeys(LogCtxKey)
		l.SetHandlers(func(ctx context.Context, in *glog.HandlerInput) {
			in.Next(ctx)
		})
	}
	return l
}

func ContextValue(ctx context.Context, key string, val any) context.Context {
	// todo fix this
	value := ctx.Value(LogCtxKey)
	gm, ok := value.(gmap.Map)
	if !ok {
		gm = gmap.Map{}
	}
	gm.Set(key, val)
	return context.WithValue(ctx, LogCtxKey, gm)
}
