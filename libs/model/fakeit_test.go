package model_test

import (
	"regexp"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/model"
)

func TestAddEtagFaker(t *testing.T) {
	model.AddETagFaker()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Foo string `fake:"{word}"`
			Bar string `fake:"{etag}"`
		}

		f := foo{}

		err := gofakeit.Struct(&f)
		require.NoError(t, err)

		assert.NotEmpty(t, f.Bar)
		assert.Regexp(t, regexp.MustCompile("^[a-z0-9]{5}$"), f.Bar)
	})
}
