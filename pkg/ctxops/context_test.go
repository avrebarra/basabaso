package ctxops_test

import (
	"context"
	"testing"

	"github.com/avrebarra/basabaso/pkg/ctxops"
	"github.com/stretchr/testify/assert"
)

func TestContextOps_Default(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		// act
		ctx := ctxops.ContextOps{}.Default()

		// assert
		assert.NotNil(t, ctx)
	})
}

func TestCreateWith(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		// act
		ctx := ctxops.CreateWith(context.Background())

		// assert
		assert.NotNil(t, ctx)
	})
}

func TestExtractFrom(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		ctx := ctxops.CreateWith(context.Background())

		// act
		out, ok := ctxops.ExtractFrom(ctx)

		// assert
		assert.True(t, ok)
		assert.NotNil(t, out)
	})

	t.Run("ok", func(t *testing.T) {
		// arrange
		// act
		out, ok := ctxops.ExtractFrom(context.Background())

		// assert
		assert.False(t, ok)
		assert.NotNil(t, out)
	})
}
