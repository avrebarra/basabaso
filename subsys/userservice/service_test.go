package userservice_test

import (
	"testing"

	"github.com/avrebarra/basabaso/subsys/userservice"
	"github.com/avrebarra/basabaso/subsys/userstore"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestUser_FromStoreData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// arrange
			// act
			data := userservice.User{}.FromStoreData(userstore.User{
				ID:             xid.New().String(),
				Username:       "randomman",
				DisplayName:    "Random Man",
				DisplayProfile: "Random Random Man Club",
			})

			// assert
			assert.NotEmpty(t, data)
		})
	})
}

func TestUser_ToStoreData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// arrange
			// act
			data := userservice.User{
				ID:             xid.New().String(),
				Username:       "nonce",
				DisplayName:    "N Once",
				DisplayProfile: "Guns N Once",
			}.ToStoreData()

			// assert
			assert.NotEmpty(t, data)
		})
	})
}
