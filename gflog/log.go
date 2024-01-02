package gflog

import (
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"slices"
)

type _contextKey struct {
}

const LogCtxKey = "vars"

func Log(name ...string) *glog.Logger {
	l := g.Log(name...)
	if !slices.Contains(l.GetCtxKeys(), LogCtxKey) {
		l.AppendCtxKeys(LogCtxKey)
	}
	return l
}

func ContextValue(ctx context.Context, key string, val any) context.Context {
	value := ctx.Value(LogCtxKey)
	if value == nil {
		value = &gmap.Map{}
		ctx = context.WithValue(ctx, LogCtxKey, value)
	}
	gm, ok := value.(*gmap.Map)
	if !ok {
		gm = &gmap.Map{}
		value = gm
		ctx = context.WithValue(ctx, LogCtxKey, value)
	}
	gm.Set(key, val)
	return ctx
}
