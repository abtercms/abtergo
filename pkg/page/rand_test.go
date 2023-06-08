package page_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/page"
)

func TestRandomPage(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		p := page.RandomPage()

		err := p.Validate()
		assert.NoError(t, err)
	})
}

func TestRandomPageList(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		const (
			min = 1
			max = 3
		)

		list := page.RandomPageList(min, max)

		assert.GreaterOrEqual(t, len(list), min)
		assert.LessOrEqual(t, len(list), max)

		for _, r := range list {
			err := r.Validate()
			assert.NoError(t, err)
		}
	})
}
