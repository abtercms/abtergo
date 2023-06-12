package block_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/block"
)

func TestRandomBlock(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tt := block.RandomBlock()

		err := tt.Validate()
		assert.NoError(t, err)
	})
}

func TestRandomBlockList(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		const (
			min = 1
			max = 3
		)

		list := block.RandomBlockList(min, max)

		assert.GreaterOrEqual(t, len(list), min)
		assert.LessOrEqual(t, len(list), max)

		for _, r := range list {
			err := r.Validate()
			assert.NoError(t, err)
		}
	})
}
