package page_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/val"
	"github.com/abtergo/abtergo/pkg/html"
	"github.com/abtergo/abtergo/pkg/page"
)

func TestPage_Clone(t *testing.T) {
	t.Run("random page can be cloned", func(t *testing.T) {
		p := page.RandomPage()

		c := p.Clone()

		assert.NotSame(t, p, c)
		assert.Equal(t, p, c)
	})

	t.Run("cloning works without temporary template data", func(t *testing.T) {
		p := page.RandomPageWithoutTemplate()

		c := p.Clone()

		assert.NotSame(t, p, c)
		assert.Equal(t, p, c)
	})

	t.Run("random page without template can be cloned", func(t *testing.T) {
		p := page.RandomPageWithoutTemplate()

		c := p.Clone()

		assert.NotSame(t, p, c)
		assert.Equal(t, p, c)
	})

	t.Run("random page with template can be cloned", func(t *testing.T) {
		p := page.RandomPageWithTemplate()

		c := p.Clone()

		assert.NotSame(t, p.Template, c.Template)
		assert.Equal(t, p.Template, c.Template)
		assert.NotSame(t, p, c)
		assert.Equal(t, p, c)
	})

	t.Run("random page with temporary template can be cloned", func(t *testing.T) {
		p := page.RandomPageWithTempTemplate()

		c := p.Clone()

		assert.NotSame(t, p.Template, c.Template)
		assert.Equal(t, p.Template, c.Template)
		assert.NotSame(t, p, c)
		assert.Equal(t, p, c)
	})
}

func TestPage_AsNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fakeit.AddPathFaker()
		fakeit.AddEtagFaker()
		expected := page.RandomPage()

		require.NotEmpty(t, expected.ID)
		require.NotEmpty(t, expected.Etag)
		require.NotEmpty(t, expected.CreatedAt)
		require.NotEmpty(t, expected.UpdatedAt)

		actual := expected.AsNew()

		expected.ID = ""
		expected.Etag = ""
		expected.CreatedAt = time.Time{}
		expected.UpdatedAt = time.Time{}

		assert.NotSame(t, expected, actual)
		assert.Equal(t, expected, actual)
	})
}

func TestPage_WithID(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		entity := page.RandomPage()

		expected := entity.ID
		updated := entity.WithID()

		assert.NotEmpty(t, expected)
		assert.Equal(t, expected, updated.ID)
	})

	t.Run("success", func(t *testing.T) {
		expected := page.Page{}

		actual := expected.WithID()

		assert.NotEmpty(t, actual.ID)
	})
}

func TestPage_WithEtag(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		entity := page.RandomPage()

		expected := entity.Etag
		updated := entity.WithEtag()

		assert.NotEmpty(t, expected)
		assert.Equal(t, expected, updated.Etag)
	})

	t.Run("success", func(t *testing.T) {
		expected := page.Page{}

		actual := expected.WithEtag()

		assert.NotEmpty(t, actual.Etag)
	})
}

func TestPage_WithTime(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		entity := page.RandomPage()

		expectedCreatedAt := entity.CreatedAt
		expectedUpdatedAt := entity.UpdatedAt

		updated := entity.WithTime()

		assert.NotEmpty(t, expectedCreatedAt)
		assert.Equal(t, expectedCreatedAt, updated.CreatedAt)
		assert.NotEmpty(t, expectedUpdatedAt)
		assert.Equal(t, expectedUpdatedAt, updated.UpdatedAt)
	})

	t.Run("success", func(t *testing.T) {
		expected := page.Page{}

		actual := expected.WithTime()

		assert.NotEmpty(t, actual.CreatedAt)
		assert.NotEmpty(t, actual.UpdatedAt)
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
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.ID = "" },
			invalidFields: []string{"id"},
		},
		{
			name:          "id, etag, created at and updated at are not required if all are empty",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.ID = ""; c.Etag = ""; c.CreatedAt = time.Time{}; c.UpdatedAt = time.Time{} },
			invalidFields: []string{},
		},
		{
			name:          "website is required",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Website = "" },
			invalidFields: []string{"website"},
		},
		{
			name:          "path is required",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Path = "" },
			invalidFields: []string{"path"},
		},
		{
			name:          "status is required",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Status = "" },
			invalidFields: []string{"status"},
		},
		{
			name:          "status must be one of Draft, Active or Inactive",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Status = "Foo" },
			invalidFields: []string{"status"},
		},
		{
			name:          "title is required",
			page:          page.RandomPageWithoutTemplate(),
			modifier:      func(c *page.Page) { c.Title = "" },
			invalidFields: []string{"title"},
		},
		{
			name:          "lead is required",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Lead = "" },
			invalidFields: []string{"lead"},
		},
		{
			name:          "body is required",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Body = "" },
			invalidFields: []string{"body"},
		},
		{
			name:          "assets with invalid header js",
			page:          page.RandomPageWithoutTemplate(),
			modifier:      func(c *page.Page) { c.Assets.HeaderJs = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid footer js",
			page:          page.RandomPageWithoutTemplate(),
			modifier:      func(c *page.Page) { c.Assets.FooterJs = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid header css",
			page:          page.RandomPageWithoutTemplate(),
			modifier:      func(c *page.Page) { c.Assets.HeaderCSS = []html.Link{{}} },
			invalidFields: []string{"rel", "href"},
		},
		{
			name:          "assets with invalid meta",
			page:          page.RandomPageWithoutTemplate(),
			modifier:      func(c *page.Page) { c.Assets.HeaderMeta = []html.Meta{{}} },
			invalidFields: []string{"name", "property", "content"},
		},
		{
			name:          "template name is required if template is not empty",
			page:          page.RandomPageWithTemplate(),
			modifier:      func(c *page.Page) { c.TemplateName = "" },
			invalidFields: []string{"template_name"},
		},
		{
			name:          "template is required if template name is not empty",
			page:          page.RandomPageWithTemplate(),
			modifier:      func(c *page.Page) { c.Template = nil },
			invalidFields: []string{"template"},
		},
		{
			name:          "template temp name and from are required if template temp until is not empty",
			page:          page.RandomPageWithTempTemplate(),
			modifier:      func(c *page.Page) { c.TemplateTempName = ""; c.TemplateTempFrom = nil },
			invalidFields: []string{"template_temp_name", "template_temp_from"},
		},
		{
			name:          "template temp name and until are required if template temp from is not empty",
			page:          page.RandomPageWithTempTemplate(),
			modifier:      func(c *page.Page) { c.TemplateTempName = ""; c.TemplateTempUntil = nil },
			invalidFields: []string{"template_temp_name", "template_temp_until"},
		},
		{
			name:          "template temp from and until are required if template temp name is not empty",
			page:          page.RandomPageWithTempTemplate(),
			modifier:      func(c *page.Page) { c.TemplateTempFrom = nil; c.TemplateTempUntil = nil },
			invalidFields: []string{"template_temp_from", "template_temp_until"},
		},
		{
			name:          "owner is required",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Owner = "" },
			invalidFields: []string{"owner"},
		},
		{
			name:          "etag is required if id, updated at or created at are present",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.Etag = "" },
			invalidFields: []string{"etag"},
		},
		{
			name:          "created at is required if id, etag, or updated at are present",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.CreatedAt = time.Time{} },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "created at must be after 2023-01-01",
			page:          page.RandomPage(),
			modifier:      func(c *page.Page) { c.CreatedAt = time.Date(2022, 10, 10, 10, 10, 10, 10, time.UTC) },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "update at must not be before created at",
			page:          page.RandomPage(),
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
				val.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}
