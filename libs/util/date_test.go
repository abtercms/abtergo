package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/util"
)

func TestCloneDate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expected := time.Now()

		actual := util.CloneDate(expected)

		assert.Equal(t, expected.UnixNano(), actual.UnixNano())
		assert.NotSame(t, expected, actual)
	})
}
