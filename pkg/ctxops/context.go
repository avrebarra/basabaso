package ctxops

import (
	"context"
	"sync"
)

type ref string

var referencekey = ref(".")

type ContextOps struct {
	data sync.Map
}

func (c ContextOps) Default() (out ContextOps) {
	return ContextOps{data: sync.Map{}}
}

func CreateWith(in context.Context) (out context.Context) {
	ctx := ContextOps{}.Default()
	out = context.WithValue(in, referencekey, &ctx)

	return
}

func ExtractFrom(ctx context.Context) (ctxRequest *ContextOps, ok bool) {
	// prevent panic due to nil on failed type conversion
	ctxRequest, ok = ctx.Value(referencekey).(*ContextOps)
	if !ok {
		ctx := ContextOps{}.Default()
		ctxRequest = &ctx
	}
	return
}
