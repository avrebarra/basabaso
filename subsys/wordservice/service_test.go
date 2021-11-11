package wordservice_test

import (
	"testing"

	"github.com/avrebarra/basabaso/subsys/wordservice"
	"github.com/avrebarra/basabaso/subsys/wordstore"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestUser_FromStoreData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		assert.NotPanics(t, func() {
			// arrange
			// act
			data := wordservice.Word{}.FromStoreData(wordstore.Word{
				ID:   xid.New().String(),
				Word: "Random",
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
			data := wordservice.Word{
				ID:   xid.New().String(),
				Word: "nonce",
			}.ToStoreData()

			// assert
			assert.NotEmpty(t, data)
		})
	})
}
