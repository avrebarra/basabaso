package jsoncast_test

import (
	"testing"

	"github.com/avrebarra/basabaso/pkg/jsoncast"
	"github.com/stretchr/testify/assert"
)

func TestCast(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		type A struct {
			Abra    string
			Cadabra string
		}
		type B struct {
			Abra    string
			Cadabra string
		}

		a := A{Abra: "a", Cadabra: "c"}
		b := B{}

		// act
		err := jsoncast.Cast(&a, &b)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, b.Abra, a.Abra)
	})

	t.Run("ok diff type A", func(t *testing.T) {
		// arrange
		type A struct {
			Abra    string
			Cadabra string
		}
		type B struct {
			Abra     string
			Cadabra  string
			Alakazam string
		}

		a := A{Abra: "a", Cadabra: "c"}
		b := B{}

		// act
		err := jsoncast.Cast(&a, &b)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, b.Abra, a.Abra)
	})

	t.Run("ok strict diff type A", func(t *testing.T) {
		// arrange
		type A struct {
			Abra    string
			Cadabra string
		}
		type B struct {
			Abra     string
			Cadabra  string
			Alakazam string
		}

		a := A{Abra: "a", Cadabra: "c"}
		b := B{}

		// act
		err := jsoncast.CastWithOpts(&a, &b, jsoncast.CastOpts{Strict: true})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, b.Abra, a.Abra)
	})

	t.Run("ok diff type B", func(t *testing.T) {
		// arrange
		type A struct {
			Abra     string
			Cadabra  string
			Alakazam string
		}
		type B struct {
			Abra    string
			Cadabra string
		}

		a := A{Abra: "a", Cadabra: "c", Alakazam: "d"}
		b := B{}

		// act
		err := jsoncast.Cast(&a, &b)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, b.Abra, a.Abra)
	})

	t.Run("err strict diff type B", func(t *testing.T) {
		// arrange
		type A struct {
			Abra     string
			Cadabra  string
			Alakazam string
		}
		type B struct {
			Abra    string
			Cadabra string
		}

		a := A{Abra: "a", Cadabra: "c", Alakazam: "d"}
		b := B{}

		// act
		err := jsoncast.CastWithOpts(&a, &b, jsoncast.CastOpts{Strict: true})

		// assert
		assert.Error(t, err)
	})

	t.Run("cov boost err marshal", func(t *testing.T) {
		// arrange
		type A struct {
			Abra     string
			Cadabra  string
			Alakazam chan (int)
		}
		type B struct {
			Abra    string
			Cadabra string
		}

		a := A{Abra: "a", Cadabra: "c", Alakazam: make(chan int)}
		b := B{}

		// act
		err := jsoncast.CastWithOpts(&a, &b, jsoncast.CastOpts{Strict: true})

		// assert
		assert.Error(t, err)
	})
}
