package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/validation"
)

func TestEntity_SetCreatedAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		n := time.Unix(time.Now().Unix(), 0)
		sut := model.NewEntity().SetCreatedAt(n).(model.Entity)

		assert.Equal(t, n, sut.GetCreatedAt())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := model.NewEntity().SetCreatedAt(time.Now()).SetETag("foo").(model.Entity)

		assert.Panics(t, func() { sut.SetCreatedAt(time.Now()) })
	})
}

func TestEntity_SetUpdatedAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		n := time.Unix(time.Now().Unix(), 0)
		sut := model.NewEntity().SetUpdatedAt(n).(model.Entity)

		assert.Equal(t, n, sut.GetUpdatedAt())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		n := time.Unix(time.Now().Unix(), 0)
		sut := model.NewEntity().SetUpdatedAt(n).SetETag("foo").(model.Entity)

		assert.Panics(t, func() { sut.SetUpdatedAt(n) })
	})
}

func TestEntity_SetDeletedAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		n := time.Unix(time.Now().Unix(), 0)
		sut := model.NewEntity().SetDeletedAt(&n).(model.Entity)

		assert.Equal(t, &n, sut.GetDeletedAt())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		n := time.Unix(time.Now().Unix(), 0)
		sut := model.NewEntity().SetDeletedAt(&n).(model.Entity).SetETag("foo").(model.Entity)

		assert.Panics(t, func() { sut.SetDeletedAt(&n) })
	})
}

func TestEntity_SetID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := model.NewEntity().SetID("foo").(model.Entity)

		assert.Equal(t, model.ID("foo"), sut.GetID())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := model.NewEntity().SetID("foo").SetETag("foo").(model.Entity)

		assert.Panics(t, func() { sut.SetID("foo") })
	})
}

func TestEntity_SetETag(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := model.NewEntity().SetETag("foo").(model.Entity)

		assert.Equal(t, model.ETag("foo"), sut.GetETag())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := model.NewEntity().SetETag("foo").(model.Entity)

		assert.Panics(t, func() { sut.SetETag("foo") })
	})
}

func TestEntity_Validate(t *testing.T) {
	tests := []struct {
		name          string
		entity        model.Entity
		modifier      func(c *model.Entity)
		invalidFields []string
	}{
		{
			name:          "valid empty",
			entity:        model.Entity{},
			modifier:      func(c *model.Entity) {},
			invalidFields: []string{},
		},
		{
			name:          "valid random",
			entity:        model.RandomEntity(),
			modifier:      func(c *model.Entity) {},
			invalidFields: []string{},
		},
		{
			name:          "invalid without etag",
			entity:        model.NewEntity(),
			modifier:      func(c *model.Entity) {},
			invalidFields: []string{"etag"},
		},
		{
			name:   "invalid without id",
			entity: model.NewEntity(),
			modifier: func(c *model.Entity) {
				c.ETag = "foo23"
				c.ID = ""
			},
			invalidFields: []string{"id"},
		},
		{
			name:   "invalid etag",
			entity: model.NewEntity(),
			modifier: func(c *model.Entity) {
				c.ID = "foo"
				c.ETag = "bar"
			},
			invalidFields: []string{"etag"},
		},
		{
			name:   "invalid missing created at",
			entity: model.RandomEntity(),
			modifier: func(c *model.Entity) {
				c.CreatedAt = time.Time{}
			},
			invalidFields: []string{"created_at"},
		},
		{
			name:   "invalid missing updated at",
			entity: model.RandomEntity(),
			modifier: func(c *model.Entity) {
				c.UpdatedAt = time.Time{}
			},
			invalidFields: []string{"updated_at"},
		},
		{
			name:   "created at before 20223",
			entity: model.RandomEntity(),
			modifier: func(c *model.Entity) {
				c.CreatedAt = time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
			},
			invalidFields: []string{"created_at"},
		},
		{
			name:   "invalid created at after updated at",
			entity: model.RandomEntity(),
			modifier: func(c *model.Entity) {
				c.UpdatedAt = c.CreatedAt.Add(-1 * time.Second)
			},
			invalidFields: []string{"updated_at"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifier(&tt.entity)

			res := tt.entity.Validate()

			if len(tt.invalidFields) == 0 {
				assert.NoError(t, res)
			} else {
				validation.AssertFieldErrorsOn(t, res, tt.invalidFields)
			}
		})
	}
}

func TestEntity_Clone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		stubEntity := model.NewEntity()

		clone := stubEntity.Clone()

		assert.Equal(t, stubEntity, clone)
	})
}

func TestEntity_IsComplete(t *testing.T) {
	t.Run("invalid empty entity", func(t *testing.T) {
		t.Parallel()

		stubEntity := model.NewEntity()

		assert.False(t, stubEntity.IsComplete())
	})

	t.Run("invalid non-empty entity, E-Tag missing", func(t *testing.T) {
		t.Parallel()

		stubEntity := model.NewEntity()

		assert.False(t, stubEntity.IsComplete())
	})

	t.Run("invalid non-empty entity, ID missing", func(t *testing.T) {
		t.Parallel()

		stubEntity := model.NewEntity()
		stubEntity.ID = ""
		stubEntity.ETag = "foo"

		assert.False(t, stubEntity.IsComplete())
	})

	t.Run("valid non-empty entity", func(t *testing.T) {
		t.Parallel()

		stubEntity := model.NewEntity()
		stubEntity.ETag = "foo"

		assert.True(t, stubEntity.IsComplete())
	})
}
