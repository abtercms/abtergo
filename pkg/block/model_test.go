package block_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/validation"
	"github.com/abtergo/abtergo/pkg/block"
)

func TestBlock_Clone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tt := block.RandomBlock()

		c := tt.Clone()

		assert.NotSame(t, tt, c)
		assert.Equal(t, tt, c)
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
			name:          "id is required if e-tag, updated at or created at are present",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.ID = "" },
			invalidFields: []string{"id"},
		},
		{
			name:  "id, e-tag, created at and updated at are not required if all are empty",
			block: block.RandomBlock(),
			modifier: func(c *block.Block) {
				c.ID = ""
				c.ETag = ""
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
			modifier:      func(c *block.Block) { c.Assets.HeaderJS = []html.Script{{}} },
			invalidFields: []string{"src"},
		},
		{
			name:          "assets with invalid footer js",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.Assets.FooterJS = []html.Script{{}} },
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
			invalidFields: []string{"name", "content"},
		},
		{
			name:          "e-tag is required if id, updated at or created at are present",
			block:         block.RandomBlock(),
			modifier:      func(c *block.Block) { c.ETag = "" },
			invalidFields: []string{"etag"},
		},
		{
			name:          "created at is required if id, e-tag, or updated at are present",
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
				validation.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}
