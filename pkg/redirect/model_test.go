package redirect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/validation"
	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestNewRedirect(t *testing.T) {
	r := redirect.NewRedirect()

	assert.NotEmpty(t, r.Entity)
	assert.NotEmpty(t, r.Entity.ID)
	assert.NotEmpty(t, r.Entity.CreatedAt)
	assert.NotEmpty(t, r.Entity.UpdatedAt)
	assert.Empty(t, r.Entity.DeletedAt)

	// TODO: fix etag issue
	// assert.NotEmpty(t, r.Entity.Etag)
}

func TestRedirect_Clone(t *testing.T) {
	t.Run("random page can be cloned", func(t *testing.T) {
		r := redirect.RandomRedirect(false)

		c := r.Clone()

		assert.NotSame(t, r, c)
		assert.Equal(t, r, c)
	})
}

func TestRedirect_AsNew(t *testing.T) {
	t.Run("clone with empty entity", func(t *testing.T) {
		t.Parallel()

		sut := redirect.RandomRedirect(false)
		clone := sut.AsNew()

		assert.NotSame(t, sut, clone)
		assert.NotEqual(t, sut, clone)

		clone.Entity = sut.Entity.Clone().(model.Entity)
		assert.Equal(t, sut, clone)
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
			redirect:      redirect.RandomRedirect(false),
			modifier:      func(c *redirect.Redirect) {},
			invalidFields: []string{},
		},
		{
			name:     "many missing fields",
			redirect: redirect.RandomRedirect(false),
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
			redirect: redirect.RandomRedirect(false),
			modifier: func(c *redirect.Redirect) {
				c.ID = ""
			},
			invalidFields: []string{"id"},
		},
		{
			name:     "missing website field",
			redirect: redirect.RandomRedirect(false),
			modifier: func(c *redirect.Redirect) {
				c.Website = ""
			},
			invalidFields: []string{"website"},
		},
		{
			name:     "missing target field",
			redirect: redirect.RandomRedirect(false),
			modifier: func(c *redirect.Redirect) {
				c.Target = ""
			},
			invalidFields: []string{"target"},
		},
		{
			name:     "invalid website field",
			redirect: redirect.RandomRedirect(false),
			modifier: func(c *redirect.Redirect) {
				c.Website = "foo"
			},
			invalidFields: []string{"website"},
		},
		{
			name:     "invalid target field",
			redirect: redirect.RandomRedirect(false),
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
