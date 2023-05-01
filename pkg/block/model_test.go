package block_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/val"
	"github.com/abtergo/abtergo/pkg/block"
	"github.com/abtergo/abtergo/pkg/html"
)

func TestBlock_Clone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tt := block.RandomBlock()

		c := tt.Clone()

		assert.NotSame(t, tt, c)
		assert.Equal(t, tt, c)
	})
}

func TestBlock_AsNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fakeit.AddEtagFaker()
		expected := block.RandomBlock()

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

func TestBlock_WithID(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		entity := block.RandomBlock()

		expected := entity.ID
		updated := entity.WithID()

		assert.NotEmpty(t, expected)
		assert.Equal(t, expected, updated.ID)
	})

	t.Run("success", func(t *testing.T) {
		expected := block.Block{}

		actual := expected.WithID()

		assert.NotEmpty(t, actual.ID)
	})
}

func TestBlock_WithEtag(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		entity := block.RandomBlock()

		expected := entity.Etag
		updated := entity.WithEtag()

		assert.NotEmpty(t, expected)
		assert.Equal(t, expected, updated.Etag)
	})

	t.Run("success", func(t *testing.T) {
		expected := block.Block{}

		actual := expected.WithEtag()

		assert.NotEmpty(t, actual.Etag)
	})
}

func TestBlock_WithTime(t *testing.T) {
	t.Run("skip", func(t *testing.T) {
		entity := block.RandomBlock()

		expectedCreatedAt := entity.CreatedAt
		expectedUpdatedAt := entity.UpdatedAt

		updated := entity.WithTime()

		assert.NotEmpty(t, expectedCreatedAt)
		assert.Equal(t, expectedCreatedAt, updated.CreatedAt)
		assert.NotEmpty(t, expectedUpdatedAt)
		assert.Equal(t, expectedUpdatedAt, updated.UpdatedAt)
	})

	t.Run("success", func(t *testing.T) {
		expected := block.Block{}

		actual := expected.WithTime()

		assert.NotEmpty(t, actual.CreatedAt)
		assert.NotEmpty(t, actual.UpdatedAt)
	})
}

func TestBlock_Validate(t *testing.T) {
	tests := []struct {
		name          string
		block         block.Block
		modifier      func(c *block.Block)
		invalidFields []string
	}{
		{
			name:          "id is required if etag, updated at or created at are present",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.ID = "" },
			invalidFields: []string{"id"},
		},
		{
			name:  "id, etag, created at and updated at are not required if all are empty",
			block: block.RandomBlock(),
			modifier: func(c *block.Block) {
				c.ID = ""
				c.Etag = ""
				c.CreatedAt = time.Time{}
				c.UpdatedAt = time.Time{}
			},
			invalidFields: []string{},
		},
		{
			name:          "website is required",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Website = "" },
			invalidFields: []string{"website"},
		},
		{
			name:          "name is required",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Name = "" },
			invalidFields: []string{"name"},
		},
		{
			name:          "body is required",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Body = "" },
			invalidFields: []string{"body"},
		},
		{
			name:          "assets with invalid header js",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Assets.HeaderJs = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid footer js",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Assets.FooterJs = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid header css",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Assets.HeaderCSS = []html.Link{{}} },
			invalidFields: []string{"rel", "href"},
		},
		{
			name:          "assets with invalid meta",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Assets.HeaderMeta = []html.Meta{{}} },
			invalidFields: []string{"name", "property", "content"},
		},
		{
			name:          "owner is required",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Owner = "" },
			invalidFields: []string{"owner"},
		},
		{
			name:          "etag is required if id, updated at or created at are present",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Etag = "" },
			invalidFields: []string{"etag"},
		},
		{
			name:          "created at is required if id, etag, or updated at are present",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.CreatedAt = time.Time{} },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "created at must be after 2023-01-01",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.CreatedAt = time.Date(2022, 10, 10, 10, 10, 10, 10, time.UTC) },
			invalidFields: []string{"created_at"},
		},
		{
			name:          "update at must not be before created at",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.UpdatedAt = c.CreatedAt.Add(-1) },
			invalidFields: []string{"updated_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(&tt.block)
			res := tt.block.Validate()
			if len(tt.invalidFields) == 0 {
				assert.NoError(t, res)
			} else {
				val.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}
