package page_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/validation"
	"github.com/abtergo/abtergo/pkg/page"
)

func TestNewPage(t *testing.T) {
	p := page.NewPage()
	assert.NotEmpty(t, p.Entity)
	assert.NotEmpty(t, p.Entity.ID)
	assert.NotEmpty(t, p.Entity.CreatedAt)
	assert.NotEmpty(t, p.Entity.UpdatedAt)
	assert.Empty(t, p.Entity.DeletedAt)

	// TODO: fix etag issue
	// assert.NotEmpty(t, p.Entity.Etag)
}

func TestPage_Clone(t *testing.T) {
	t.Run("random page can be cloned", func(t *testing.T) {
		p := page.RandomPage(false)

		c := p.Clone()

		assert.NotSame(t, p, c)
		assert.Equal(t, p, c)
	})
}

func TestPage_Validate(t *testing.T) {
	tests := []struct {
		name          string
		page          page.Page
		modifier      func(c *page.Page)
		invalidFields []string
	}{
		{
			name:          "id is required if etag, updated at or created at are present",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.ID = "" },
			invalidFields: []string{"id"},
		},
		{
			name:          "id, etag, created at and updated at are not required if all are empty",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.ID = ""; c.Etag = ""; c.CreatedAt = time.Time{}; c.UpdatedAt = time.Time{} },
			invalidFields: []string{},
		},
		{
			name:          "website is required",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Website = "" },
			invalidFields: []string{"website"},
		},
		{
			name:          "path is required",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Path = "" },
			invalidFields: []string{"path"},
		},
		{
			name:          "status is required",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Status = "" },
			invalidFields: []string{"status"},
		},
		{
			name:          "status must be one of Draft, Active or Inactive",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Status = "Foo" },
			invalidFields: []string{"status"},
		},
		{
			name:          "title is required",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Title = "" },
			invalidFields: []string{"title"},
		},
		{
			name:          "lead is required",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Lead = "" },
			invalidFields: []string{"lead"},
		},
		{
			name:          "body is required",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Body = "" },
			invalidFields: []string{"body"},
		},
		{
			name:          "assets with invalid header js",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Assets.HeaderJS = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid footer js",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Assets.FooterJS = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid header css",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Assets.HeaderCSS = []html.Link{{}} },
			invalidFields: []string{"rel", "href"},
		},
		{
			name:          "assets with invalid meta",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Assets.HeaderMeta = []html.Meta{{}} },
			invalidFields: []string{"name", "content"},
		},
		{
			name:          "etag is required if id, updated at or created at are present",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.Etag = "" },
			invalidFields: []string{"etag"},
		},
		{
			name:          "created at is required if id, etag, or updated at are present",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.CreatedAt = time.Time{} },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "created at must be after 2023-01-01",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.CreatedAt = time.Date(2022, 10, 10, 10, 10, 10, 10, time.UTC) },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "update at must not be before created at",
			page:          page.RandomPage(false),
			modifier:      func(c *page.Page) { c.UpdatedAt = c.CreatedAt.Add(-1) },
			invalidFields: []string{"updated_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(&tt.page)

			res := tt.page.Validate()

			if len(tt.invalidFields) == 0 {
				assert.NoError(t, res)
			} else {
				validation.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}
