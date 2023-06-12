package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/validation"
)

func TestEntity_CalculateETag(t *testing.T) {
	const stubPayload = "foo"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.CalculateETag(stubPayload).(Entity)

		assert.NotEmpty(t, sut.GetETag())
	})
}

func TestEntity_IsModified(t *testing.T) {
	const stubPayload = "foo"

	t.Run("not modified", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.CalculateETag(stubPayload).(Entity)

		got := sut.IsModified(stubPayload)

		assert.False(t, got)
	})

	t.Run("modified", func(t *testing.T) {
		t.Parallel()

		modifiedPayload := "bar"

		stubEntity := NewEntity()
		stubEntity.CalculateETag(stubPayload)

		got := stubEntity.IsModified(modifiedPayload)

		assert.True(t, got)
	})
}

func TestEntity_SetCreatedAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		n := time.Unix(time.Now().Unix(), 0)
		sut := NewEntity()

		sut = sut.SetCreatedAt(n).(Entity)

		assert.Equal(t, n, sut.GetCreatedAt())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.SetCreatedAt(time.Now()).(Entity)

		sut = sut.SetETag("foo").(Entity)

		assert.Panics(t, func() { sut.SetCreatedAt(time.Now()) })
	})
}

func TestEntity_SetUpdatedAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		n := time.Unix(time.Now().Unix(), 0)
		sut = sut.SetUpdatedAt(n).(Entity)

		assert.Equal(t, n, sut.GetUpdatedAt())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		n := time.Unix(time.Now().Unix(), 0)
		sut = sut.SetUpdatedAt(n).(Entity)

		sut = sut.SetETag("foo").(Entity)

		assert.Panics(t, func() { sut.SetUpdatedAt(n) })
	})
}

func TestEntity_SetDeletedAt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		n := time.Unix(time.Now().Unix(), 0)
		sut = sut.SetDeletedAt(&n).(Entity)

		assert.Equal(t, &n, sut.GetDeletedAt())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		n := time.Unix(time.Now().Unix(), 0)
		sut = sut.SetDeletedAt(&n).(Entity)

		sut = sut.SetETag("foo").(Entity)

		assert.Panics(t, func() { sut.SetDeletedAt(&n) })
	})
}

func TestEntity_SetID(t *testing.T) {
	const stubID = "foo"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.SetID(stubID).(Entity)

		assert.Equal(t, stubID, sut.GetID())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.SetID(stubID).(Entity)

		sut = sut.SetETag("foo").(Entity)

		assert.Panics(t, func() { sut.SetID(stubID) })
	})
}

func TestEntity_SetETag(t *testing.T) {
	const stubETag = "foo"

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.SetETag(stubETag).(Entity)

		assert.Equal(t, stubETag, sut.GetETag())
	})

	t.Run("error setting updated after e-tag was added", func(t *testing.T) {
		t.Parallel()

		sut := NewEntity()
		sut = sut.SetETag(stubETag).(Entity)

		assert.Panics(t, func() { sut.SetETag(stubETag) })
	})
}

func TestEntity_Validate(t *testing.T) {
	tests := []struct {
		name          string
		entity        Entity
		modifier      func(c *Entity)
		invalidFields []string
	}{
		{
			name:          "valid empty",
			entity:        Entity{},
			modifier:      func(c *Entity) {},
			invalidFields: []string{},
		},
		{
			name:          "valid random",
			entity:        RandomEntity(),
			modifier:      func(c *Entity) {},
			invalidFields: []string{},
		},
		{
			name:          "invalid without etag",
			entity:        NewEntity(),
			modifier:      func(c *Entity) {},
			invalidFields: []string{"etag"},
		},
		{
			name:   "invalid without id",
			entity: NewEntity(),
			modifier: func(c *Entity) {
				c.ETag = "foo23"
				c.ID = ""
			},
			invalidFields: []string{"id"},
		},
		{
			name:   "invalid etag",
			entity: NewEntity(),
			modifier: func(c *Entity) {
				c.ID = "foo"
				c.ETag = "bar"
			},
			invalidFields: []string{"etag"},
		},
		{
			name:   "invalid missing created at",
			entity: RandomEntity(),
			modifier: func(c *Entity) {
				c.CreatedAt = time.Time{}
			},
			invalidFields: []string{"created_at"},
		},
		{
			name:   "invalid missing updated at",
			entity: RandomEntity(),
			modifier: func(c *Entity) {
				c.UpdatedAt = time.Time{}
			},
			invalidFields: []string{"updated_at"},
		},
		{
			name:   "created at before 20223",
			entity: RandomEntity(),
			modifier: func(c *Entity) {
				c.CreatedAt = time.Date(2022, 12, 31, 0, 0, 0, 0, time.UTC)
			},
			invalidFields: []string{"created_at"},
		},
		{
			name:   "invalid created at after updated at",
			entity: RandomEntity(),
			modifier: func(c *Entity) {
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
