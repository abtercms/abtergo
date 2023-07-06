package redirect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/validation"
	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestRedirect_Clone(t *testing.T) {
	t.Run("random page can be cloned", func(t *testing.T) {
		r := redirect.RandomRedirect()

		c := r.Clone()

		assert.NotSame(t, r, c)
		assert.Equal(t, r, c)
	})
}

func TestRedirect_Validate(t *testing.T) {
	tests := []struct {
		name          string
		redirect      redirect.Redirect
		modifier      func(c *redirect.Redirect)
		invalidFields []string
	}{
		{
			name:          "valid redirect",
			redirect:      redirect.RandomRedirect(),
			modifier:      func(c *redirect.Redirect) {},
			invalidFields: []string{},
		},
		{
			name:     "many missing fields",
			redirect: redirect.RandomRedirect(),
			modifier: func(c *redirect.Redirect) {
				c.ID = ""
				c.Website = ""
				c.Path = ""
				c.Target = ""
			},
			invalidFields: []string{"id", "website", "path", "target"},
		},
		{
			name:     "missing id field",
			redirect: redirect.RandomRedirect(),
			modifier: func(c *redirect.Redirect) {
				c.ID = ""
			},
			invalidFields: []string{"id"},
		},
		{
			name:     "missing website field",
			redirect: redirect.RandomRedirect(),
			modifier: func(c *redirect.Redirect) {
				c.Website = ""
			},
			invalidFields: []string{"website"},
		},
		{
			name:     "missing target field",
			redirect: redirect.RandomRedirect(),
			modifier: func(c *redirect.Redirect) {
				c.Target = ""
			},
			invalidFields: []string{"target"},
		},
		{
			name:     "invalid website field",
			redirect: redirect.RandomRedirect(),
			modifier: func(c *redirect.Redirect) {
				c.Website = "foo"
			},
			invalidFields: []string{"website"},
		},
		{
			name:     "invalid target field",
			redirect: redirect.RandomRedirect(),
			modifier: func(c *redirect.Redirect) {
				c.Target = "foo"
			},
			invalidFields: []string{"target"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(&tt.redirect)

			res := tt.redirect.Validate()

			if len(tt.invalidFields) == 0 {
				assert.NoError(t, res)
			} else {
				validation.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}
