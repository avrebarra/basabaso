package ctxops_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/avrebarra/basabaso/pkg/ctxops"
	"github.com/stretchr/testify/assert"
)

func TestSetData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())

		// act
		ctxops.SetData(ctx, "a", "a")
		ctxops.SetData(ctx, "b", "b")

		// assert
		out, ok := ctxops.ExtractFrom(ctx)
		assert.True(t, ok)
		assert.NotEmpty(t, out)
	})
}

func TestGetData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())
		ctxops.SetData(ctx, "a", "a")
		ctxops.SetData(ctx, "b", "b")

		// act
		outa := ctxops.GetData(ctx, "a")
		outb := ctxops.GetData(ctx, "b")

		// assert
		assert.Equal(t, outa, "a")
		assert.Equal(t, outb, "b")
	})
}

func TestGetSetOpsID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		opsid := "abc"
		ctx := ctxops.CreateWith(context.Background())

		// act
		ctxops.SetOpsID(ctx, opsid)
		outopsid := ctxops.GetOpsID(ctx)

		// assert
		assert.Equal(t, opsid, outopsid)
	})

	t.Run("ok multiple sets", func(t *testing.T) {
		// arrange
		opsid := "abc"
		ctx := ctxops.CreateWith(context.Background())

		// act
		ctxops.SetOpsID(ctx, opsid)
		ctxops.SetOpsID(ctx, "abc")
		outopsid := ctxops.GetOpsID(ctx)

		// assert
		assert.Equal(t, opsid, outopsid)
	})

	t.Run("ok not set", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())

		// act
		outopsid := ctxops.GetOpsID(ctx)

		// assert
		assert.Equal(t, "", outopsid)
	})

	t.Run("coverage: wrong type", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())
		ctxops.SetData(ctx, "*ctxkey:core:opsid::", 123)

		// act
		outopsid := ctxops.GetOpsID(ctx)

		// assert
		assert.Equal(t, "", outopsid)
	})
}

func TestAddWarningListWarnings(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())
		err := fmt.Errorf("random error")

		ctxops.SetData(ctx, "abcd", "lalala") // add pollution

		// act
		ctxops.AddWarning(ctx, err.Error())
		ctxops.AddWarning(ctx, err.Error())
		out := ctxops.ListWarnings(ctx)

		// assert
		assert.NotEmpty(t, out)
		assert.Equal(t, 2, len(out))
	})

	t.Run("ok empty context", func(t *testing.T) {
		// arrange
		ctx := context.Background()
		err := fmt.Errorf("random error")

		// act
		ctxops.AddWarning(ctx, err.Error())
		out := ctxops.ListWarnings(ctx)

		// assert
		assert.Empty(t, out)
	})
}

func TestAddVarAndListVars(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())
		ctxops.SetData(ctx, "abcd", "lalala") // add pollution

		// act
		ctxops.AddVar(ctx, "k1", "v1")
		ctxops.AddVar(ctx, "k2", "v2")
		out := ctxops.ListVars(ctx)

		// assert
		assert.NotEmpty(t, out)
		assert.Equal(t, 2, len(out))
		assert.Equal(t, "v1", out["k1"])
	})
}

func TestAddProcessAndListProcesses(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())
		ctxops.SetData(ctx, "abcd", "lalala") // add pollution
		ctxops.AddVar(ctx, "k1", "v1")        // add pollution

		px := ctxops.Process{
			Name: "outgoing/firebase",
			Data: map[string]interface{}{
				"in/id":     "1234",
				"in/target": "@resource",
			},
			StartAt:  time.Now(),
			FinishAt: time.Now(),
		}

		px.Data["out/status"] = "200"
		px.Data["out/message"] = "finished success"
		px.FinishAt = time.Now()

		// act
		ctxops.AddProcess(ctx, px)
		out := ctxops.ListProcesses(ctx)

		// assert
		assert.NotEmpty(t, out)
		assert.Equal(t, "200", out[0].Data["out/status"])
	})
}
