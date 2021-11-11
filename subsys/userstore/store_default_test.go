package userstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/avrebarra/basabaso/subsys/userstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewMongo(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		// act
		x, err := userstore.NewMongo(userstore.ConfigMongo{
			DB:               &mongo.Database{},
			CollecNamePrefix: "",
		})

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, x)
	})

	t.Run("err bad deps", func(t *testing.T) {
		// arrange
		// act
		x, err := userstore.NewMongo(userstore.ConfigMongo{
			DB:               nil,
			CollecNamePrefix: "",
		})

		// assert
		assert.Empty(t, x)
		assert.Error(t, err)
	})
}

func TestIntegration(t *testing.T) {
	const (
		EnableIntegration     = true
		MongoURL              = "mongodb://root:rootpw@localhost:27017"
		MongoDBName           = "testing"
		MongoCollectionPrefix = "basabaso_"
	)

	if !EnableIntegration {
		t.SkipNow()
		return
	}

	// arrange
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoconn, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoURL))
	require.NoError(t, err)

	store, err := userstore.NewMongo(userstore.ConfigMongo{
		DB:               mongoconn.Database(MongoDBName),
		CollecNamePrefix: MongoCollectionPrefix,
	})
	require.NoError(t, err)

	t.Run("err persist bad input", func(t *testing.T) {
		// arrange
		// act
		err = store.Persist(ctx, nil)

		// assert
		assert.Error(t, err)
	})

	t.Run("ok persist user", func(t *testing.T) {
		// arrange
		usr := userstore.User{
			ID:             "testing-id",
			Username:       "@lalamove",
			DisplayName:    "lalalala",
			DisplayProfile: "Lala Lalala",
		}

		// act
		err = store.Persist(ctx, &usr)

		// assert
		assert.NoError(t, err)
	})

	t.Run("err find bad input", func(t *testing.T) {
		// arrange
		// act
		out, tot, err := store.Find(ctx, userstore.FindInput{})

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, int(tot))
		assert.Empty(t, out)
	})

	t.Run("ok find empty result", func(t *testing.T) {
		// arrange
		// act
		out, tot, err := store.Find(ctx, userstore.FindInput{
			PagingLimit:  10,
			PagingOffset: 0,
			ID:           "testing",
		})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, int(tot))
		assert.Empty(t, out)
	})

	t.Run("ok find result hit", func(t *testing.T) {
		// arrange
		// act
		out, tot, err := store.Find(ctx, userstore.FindInput{
			PagingLimit:  10,
			PagingOffset: 0,
			ID:           "testing-id",
		})

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, out)
		assert.Equal(t, 1, int(tot))
		assert.Equal(t, 1, len(out))
	})

	t.Run("ok find result skipped by offset", func(t *testing.T) {
		// arrange
		// act
		out, tot, err := store.Find(ctx, userstore.FindInput{
			PagingLimit:  10,
			PagingOffset: 2,
			ID:           "testing-id",
		})

		// assert
		assert.NoError(t, err)
		assert.Empty(t, out)
		assert.Equal(t, 1, int(tot))
		assert.Equal(t, 0, len(out))
	})
}
