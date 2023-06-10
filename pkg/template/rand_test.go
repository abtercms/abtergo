package template_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/template"
)

func TestRandomTemplate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		tt := template.RandomTemplate(false)

		err := tt.Validate()
		assert.NoError(t, err)
	})
}

func TestRandomTemplateList(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		const (
			min = 1
			max = 3
		)

		list := template.RandomTemplateList(min, max)

		assert.GreaterOrEqual(t, len(list), min)
		assert.LessOrEqual(t, len(list), max)

		for _, r := range list {
			err := r.Validate()
			assert.NoError(t, err)
		}
	})
}
