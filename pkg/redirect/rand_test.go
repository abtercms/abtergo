package redirect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestRandomRedirect(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r := redirect.RandomRedirect()

		err := r.Validate()
		assert.NoError(t, err)
	})
}

func TestRandomRedirectList(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		const (
			min = 1
			max = 3
		)

		list := redirect.RandomRedirectList(min, max)

		assert.GreaterOrEqual(t, len(list), min)
		assert.LessOrEqual(t, len(list), max)

		for _, r := range list {
			err := r.Validate()
			assert.NoError(t, err)
		}
	})
}
