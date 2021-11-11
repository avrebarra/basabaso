package ctxops

import (
	"context"
	"strings"
	"time"

	"github.com/rs/xid"
)

func SetData(ctx context.Context, key string, val interface{}) {
	c, _ := ExtractFrom(ctx)
	c.data.Store(key, val)
}

func GetData(ctx context.Context, key string) (val interface{}) {
	c, _ := ExtractFrom(ctx)
	val, _ = c.data.Load(key)
	return
}

// ***

var prefixKeyOpsID = "*ctxkey:core:opsid::"

func SetOpsID(ctx context.Context, in string) {
	if val := GetData(ctx, prefixKeyOpsID); val != nil {
		return
	}
	SetData(ctx, prefixKeyOpsID, in)
}

func GetOpsID(ctx context.Context) (out string) {
	var val interface{}
	if val = GetData(ctx, prefixKeyOpsID); val == nil {
		return
	}
	if strval, ok := val.(string); ok {
		return strval
	}
	return
}

// ***

var prefixKeyWarning = "*ctxkey:core:warning::"

func AddWarning(ctx context.Context, val string) {
	key := prefixKeyWarning + xid.New().String()
	SetData(ctx, key, val)
}

func ListWarnings(ctx context.Context) (out []string) {
	c, _ := ExtractFrom(ctx)
	out = []string{}

	c.data.Range(func(key, val interface{}) bool {
		if strkey, ok := key.(string); !ok || !strings.HasPrefix(strkey, prefixKeyWarning) {
			return true
		}
		if strval, ok := val.(string); ok {
			out = append(out, strval)
		}
		return true
	})

	return
}

// ***

var prefixKeyVars = "*ctxkey:core:vars::"

type kv struct {
	key string
	val interface{}
}

func AddVar(ctx context.Context, varkey string, varval interface{}) {
	key := prefixKeyVars + xid.New().String()
	SetData(ctx, key, kv{key: varkey, val: varval})
}

func ListVars(ctx context.Context) (out map[string]interface{}) {
	out = map[string]interface{}{}

	c, _ := ExtractFrom(ctx)
	c.data.Range(func(key, val interface{}) bool {
		if strkey, ok := key.(string); !ok || !strings.HasPrefix(strkey, prefixKeyVars) {
			return true
		}
		if kvval, ok := val.(kv); ok {
			out[kvval.key] = kvval.val
		}
		return true
	})

	return
}

// ***

var prefixKeyProcess = "*ctxkey:core:process::"

type Process struct {
	Name     string
	Data     map[string]interface{}
	StartAt  time.Time
	FinishAt time.Time
}

func AddProcess(ctx context.Context, proc Process) {
	key := prefixKeyProcess + xid.New().String()
	SetData(ctx, key, proc)
}

func ListProcesses(ctx context.Context) (out []Process) {
	out = []Process{}

	c, _ := ExtractFrom(ctx)
	c.data.Range(func(key, val interface{}) bool {
		if strkey, ok := key.(string); !ok || !strings.HasPrefix(strkey, prefixKeyProcess) {
			return true
		}
		if kvval, ok := val.(Process); ok {
			out = append(out, kvval)
		}
		return true
	})

	return
}
