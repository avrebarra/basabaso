package logutil_test

import (
	"io/ioutil"
	"testing"

	"github.com/avrebarra/basabaso/pkg/logutil"
	"github.com/stretchr/testify/assert"
)

func TestPrettyPrinter_Write(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// arrange
		pp := logutil.PrettyPrinter{Enable: true, Out: ioutil.Discard}

		// act
		len, err := pp.Write([]byte(`{"level":"debug","message":"ok","time":"2019-07-04T13:33:03.969Z"}`))

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, len)
	})

	t.Run("ok non json", func(t *testing.T) {
		// arrange
		pp := logutil.PrettyPrinter{Enable: true, Out: ioutil.Discard}

		// act
		len, err := pp.Write([]byte("ok"))

		// assert
		assert.NoError(t, err)
		assert.NotEmpty(t, len)
	})

	t.Run("ok disabled", func(t *testing.T) {
		// arrange
		pp := logutil.PrettyPrinter{Enable: false, Out: ioutil.Discard}

		// act
		len, err := pp.Write([]byte(""))

		// assert
		assert.NoError(t, err)
		assert.Empty(t, len)
	})
}
