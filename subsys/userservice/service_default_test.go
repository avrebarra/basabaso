package userservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/avrebarra/basabaso/subsys/userservice"
	"github.com/avrebarra/basabaso/subsys/userstore"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		// act
		svc, err := userservice.New(userservice.Config{
			Store: &StoreMock{},
		})

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, svc)
	})

	t.Run("err bad deps", func(t *testing.T) {
		// arrange
		// act
		svc, err := userservice.New(userservice.Config{
			Store: nil,
		})

		// assert
		assert.Error(t, err)
		assert.Empty(t, svc)
	})
}

func TestDefault_FindUsers(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		svc, err := userservice.New(userservice.Config{
			Store: &StoreMock{
				FindFunc: func(ctx context.Context, in userstore.FindInput) ([]userstore.User, int64, error) {
					return []userstore.User{{ID: xid.New().String()}}, 1, nil
				},
			},
		})
		require.NoError(t, err)

		// act
		out, err := svc.FindUsers(context.Background(), userservice.FindUsersInput{})

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, out.Total)
		assert.NotEmpty(t, out.Users)
	})

	t.Run("err store failure", func(t *testing.T) {
		// arrange
		svc, err := userservice.New(userservice.Config{
			Store: &StoreMock{
				FindFunc: func(ctx context.Context, in userstore.FindInput) ([]userstore.User, int64, error) {
					return []userstore.User{}, 0, fmt.Errorf("subtle error")
				},
			},
		})
		require.NoError(t, err)

		// act
		out, err := svc.FindUsers(context.Background(), userservice.FindUsersInput{})

		// assert
		assert.Error(t, err)
		assert.Empty(t, out.Total)
		assert.Empty(t, out.Users)
	})
}

func TestDefault_PersistUser(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		svc, err := userservice.New(userservice.Config{
			Store: &StoreMock{
				PersistFunc: func(ctx context.Context, in *userstore.User) error {
					return nil
				},
			},
		})
		require.NoError(t, err)

		// act
		err = svc.PersistUser(context.Background(), userservice.User{
			ID:             xid.New().String(),
			Username:       xid.New().String(),
			DisplayName:    "Abata",
			DisplayProfile: "Sakhakho",
		})

		// assert
		assert.NoError(t, err)
	})

	t.Run("err bad input", func(t *testing.T) {
		// arrange
		svc, err := userservice.New(userservice.Config{
			Store: &StoreMock{},
		})
		require.NoError(t, err)

		// act
		err = svc.PersistUser(context.Background(), userservice.User{
			// ID:             xid.New().String(),
			Username:       xid.New().String(),
			DisplayName:    "Abata",
			DisplayProfile: "Sakhakho",
		})

		// assert
		assert.Error(t, err)
	})

	t.Run("err store failure", func(t *testing.T) {
		// arrange
		svc, err := userservice.New(userservice.Config{
			Store: &StoreMock{
				PersistFunc: func(ctx context.Context, in *userstore.User) error {
					return fmt.Errorf("store failure")
				},
			},
		})
		require.NoError(t, err)

		// act
		err = svc.PersistUser(context.Background(), userservice.User{
			ID:             xid.New().String(),
			Username:       xid.New().String(),
			DisplayName:    "Abata",
			DisplayProfile: "Sakhakho",
		})

		// assert
		assert.Error(t, err)
	})
}
