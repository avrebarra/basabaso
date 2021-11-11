package wordservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/avrebarra/basabaso/subsys/wordservice"
	"github.com/avrebarra/basabaso/subsys/wordstore"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		// act
		svc, err := wordservice.New(wordservice.Config{
			Store: &StoreMock{},
		})

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, svc)
	})

	t.Run("err bad deps", func(t *testing.T) {
		// arrange
		// act
		svc, err := wordservice.New(wordservice.Config{
			Store: nil,
		})

		// assert
		assert.Error(t, err)
		assert.Empty(t, svc)
	})
}

func TestDefault_FindWords(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		svc, err := wordservice.New(wordservice.Config{
			Store: &StoreMock{
				FindFunc: func(ctx context.Context, in wordstore.FindInput) ([]wordstore.Word, int64, error) {
					return []wordstore.Word{{ID: xid.New().String()}}, 1, nil
				},
			},
		})
		require.NoError(t, err)

		// act
		out, err := svc.FindWords(context.Background(), wordservice.FindWordsInput{})

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, out.Total)
		assert.NotEmpty(t, out.Users)
	})

	t.Run("err store failure", func(t *testing.T) {
		// arrange
		svc, err := wordservice.New(wordservice.Config{
			Store: &StoreMock{
				FindFunc: func(ctx context.Context, in wordstore.FindInput) ([]wordstore.Word, int64, error) {
					return []wordstore.Word{}, 0, fmt.Errorf("subtle error")
				},
			},
		})
		require.NoError(t, err)

		// act
		out, err := svc.FindWords(context.Background(), wordservice.FindWordsInput{})

		// assert
		assert.Error(t, err)
		assert.Empty(t, out.Total)
		assert.Empty(t, out.Users)
	})
}

func TestDefault_PersistWord(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		svc, err := wordservice.New(wordservice.Config{
			Store: &StoreMock{
				PersistFunc: func(ctx context.Context, in *wordstore.Word) error {
					return nil
				},
			},
		})
		require.NoError(t, err)

		// act
		err = svc.PersistWord(context.Background(), wordservice.Word{
			ID:   xid.New().String(),
			Word: "Xid",
		})

		// assert
		assert.NoError(t, err)
	})

	t.Run("err bad input", func(t *testing.T) {
		// arrange
		svc, err := wordservice.New(wordservice.Config{
			Store: &StoreMock{},
		})
		require.NoError(t, err)

		// act
		err = svc.PersistWord(context.Background(), wordservice.Word{
			// ID:             xid.New().String(),
			Word: "Xid",
		})

		// assert
		assert.Error(t, err)
	})

	t.Run("err store failure", func(t *testing.T) {
		// arrange
		svc, err := wordservice.New(wordservice.Config{
			Store: &StoreMock{
				PersistFunc: func(ctx context.Context, in *wordstore.Word) error {
					return fmt.Errorf("store failure")
				},
			},
		})
		require.NoError(t, err)

		// act
		err = svc.PersistWord(context.Background(), wordservice.Word{
			ID:   xid.New().String(),
			Word: "Xid",
		})

		// assert
		assert.Error(t, err)
	})
}
