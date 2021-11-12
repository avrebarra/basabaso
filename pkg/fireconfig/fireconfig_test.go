package fireconfig_test

import (
	"context"
	"testing"

	"github.com/avrebarra/basabaso/pkg/fireconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("err bad config", func(t *testing.T) {
		// arrange
		// act
		f, err := fireconfig.New(fireconfig.Config{
			FirebaseConfigJSON: nil,
			CollecName:         "configs",
			ContentFieldName:   "abc",
		})

		// assert
		assert.Error(t, err)
		assert.Empty(t, f)
	})

	t.Run("err bad init", func(t *testing.T) {
		// arrange
		// act
		f, err := fireconfig.New(fireconfig.Config{
			FirebaseConfigJSON: []byte("-"),
			CollecName:         "configs",
			ContentFieldName:   "abc",
		})

		// assert
		assert.Error(t, err)
		assert.Empty(t, f)
	})
}

func TestIntegration(t *testing.T) {
	const (
		EnableIntegration      = false
		FirebaseCredentialJSON = `COPY_JSON_HERE`
		FirestoreCollection    = "configs"
	)

	if !EnableIntegration {
		t.SkipNow()
		return
	}

	fc, err := fireconfig.New(fireconfig.Config{
		FirebaseConfigJSON: []byte(FirebaseCredentialJSON),
		CollecName:         FirestoreCollection,
		ContentFieldName:   "content",
	})
	require.NoError(t, err)
	defer fc.Close(context.Background())

	t.Run("ok get", func(t *testing.T) {
		// arrange
		// act
		out, err := fc.Get(context.Background(), "sample")

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, out)
	})

	t.Run("err get bad field", func(t *testing.T) {
		// arrange
		fcx, err := fireconfig.New(fireconfig.Config{
			FirebaseConfigJSON: []byte(FirebaseCredentialJSON),
			CollecName:         FirestoreCollection,
			ContentFieldName:   "contentx",
		})
		require.NoError(t, err)

		// act
		out, err := fcx.Get(context.Background(), "sample")

		// assert
		assert.Error(t, err)
		assert.Empty(t, out)
	})

	t.Run("err get not found", func(t *testing.T) {
		// arrange
		// act
		out, err := fc.Get(context.Background(), "sample-nonexistent")

		// assert
		assert.Error(t, err)
		assert.Empty(t, out)
	})
}
