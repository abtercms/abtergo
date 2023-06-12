package template_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/validation"
	"github.com/abtergo/abtergo/pkg/template"
)

func TestTemplate_Clone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tt := template.RandomTemplate(false)

		c := tt.Clone()

		assert.NotSame(t, tt, c)
		assert.Equal(t, tt, c)
	})
}

func TestTemplate_Validate(t *testing.T) {
	tests := []struct {
		name          string
		template      template.Template
		modifier      func(c *template.Template)
		invalidFields []string
	}{
		{
			name:          "id is required if etag, updated at or created at are present",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.ID = "" },
			invalidFields: []string{"id"},
		},
		{
			name:     "id, etag, created at and updated at are not required if all are empty",
			template: template.RandomTemplate(false),
			modifier: func(c *template.Template) {
				c.ID = ""
				c.ETag = ""
				c.CreatedAt = time.Time{}
				c.UpdatedAt = time.Time{}
			},
			invalidFields: []string{},
		},
		{
			name:          "website is required",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Website = "" },
			invalidFields: []string{"website"},
		},
		{
			name:          "name is required",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Name = "" },
			invalidFields: []string{"name"},
		},
		{
			name:          "body is required",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Body = "" },
			invalidFields: []string{"body"},
		},
		{
			name:          "assets with invalid header js",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Assets.HeaderJS = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid footer js",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Assets.FooterJS = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid header css",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Assets.HeaderCSS = []html.Link{{}} },
			invalidFields: []string{"rel", "href"},
		},
		{
			name:          "assets with invalid meta",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.Assets.HeaderMeta = []html.Meta{{}} },
			invalidFields: []string{"name", "content"},
		},
		{
			name:          "etag is required if id, updated at or created at are present",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.ETag = "" },
			invalidFields: []string{"etag"},
		},
		{
			name:          "created at is required if id, etag, or updated at are present",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.CreatedAt = time.Time{} },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "created at must be after 2023-01-01",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.CreatedAt = time.Date(2022, 10, 10, 10, 10, 10, 10, time.UTC) },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "update at must not be before created at",
			template:      template.RandomTemplate(false),
			modifier:      func(c *template.Template) { c.UpdatedAt = c.CreatedAt.Add(-1) },
			invalidFields: []string{"updated_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(&tt.template)
			res := tt.template.Validate()
			if len(tt.invalidFields) == 0 {
				assert.NoError(t, res)
			} else {
				validation.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}
