package repo_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mocks "github.com/abtergo/abtergo/mocks/libs/repo"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
)

type testEntity struct {
	model.Entity

	Website string `json:"website,omitempty"`
	Path    string `json:"path,omitempty"`
	Foo     string `json:"foo"`
}

func (e testEntity) GetUniqueKey() model.Key {
	return model.KeyFromStrings(e.Website, e.Path)
}

func (e testEntity) Clone() model.EntityInterface {
	c := e.AsNew().(testEntity)
	c.Entity = e.Entity.Clone().(model.Entity)

	return e
}

func (e testEntity) AsNew() model.EntityInterface {
	return testEntity{
		Entity: model.Entity{},
		Foo:    e.Foo,
	}
}

type testFilter struct{}

func (f testFilter) Match(ctx context.Context, e testEntity) (bool, error) {
	_, _ = ctx, e

	return true, nil
}

func TestInMemoryRepo_Create(t *testing.T) {
	t.Run("incomplete entity", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}.AsNew().(testEntity)

		_, err := sut.Create(ctx, stubEntity)
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		assert.NoError(t, err)

		assert.Equal(t, stubEntity, storedEntity)
	})
}

func TestInMemoryRepo_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		retrievedEntity, err := sut.GetByID(ctx, storedEntity.ID)

		assert.NoError(t, err)
		assert.Equal(t, retrievedEntity, storedEntity)
		assert.Equal(t, stubEntity.AsNew(), storedEntity.AsNew())
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		retrievedEntity, err := sut.GetByID(ctx, storedEntity.ID)

		assert.NoError(t, err)
		assert.Equal(t, retrievedEntity, storedEntity)
		assert.Equal(t, stubEntity.AsNew(), storedEntity.AsNew())
	})
}

func TestInMemoryRepo_Retrieve(t *testing.T) {
	t.Run("missing index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity:  model.NewEntity(),
			Foo:     "foo",
			Website: "bar",
			Path:    "quix",
		}
		stubEntity.ETag = "baz"

		_, err := sut.GetByKey(ctx, stubEntity.GetUniqueKey())

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity:  model.NewEntity(),
			Foo:     "foo",
			Website: "bar",
			Path:    "quix",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		retrievedEntity, err := sut.GetByKey(ctx, stubEntity.GetUniqueKey())

		assert.NoError(t, err)
		assert.Equal(t, retrievedEntity, storedEntity)
		assert.Equal(t, stubEntity.AsNew(), storedEntity.AsNew())
	})
}

func TestInMemoryRepo_Update(t *testing.T) {
	t.Run("incomplete entity", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}.AsNew().(testEntity)

		_, err := sut.Update(ctx, stubEntity, "bar")

		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, arr.HTTPStatusFromError(err))
	})

	t.Run("error not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		_, err := sut.Update(ctx, stubEntity, "bar")

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, arr.HTTPStatusFromError(err))
	})

	t.Run("error e-tag mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		_, err = sut.Update(ctx, storedEntity, "bar")

		assert.Error(t, err)
		assert.Equal(t, http.StatusPreconditionFailed, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		stubEntity2 := testEntity{
			Entity: storedEntity.Entity,
			Foo:    "baz",
		}

		updatedEntity, err := sut.Update(ctx, stubEntity2, storedEntity.GetETag())
		assert.NoError(t, err)

		retrievedEntity, err := sut.GetByID(ctx, updatedEntity.ID)

		assert.NoError(t, err)
		assert.NotEqual(t, retrievedEntity, storedEntity)
		assert.Equal(t, retrievedEntity, updatedEntity)
	})
}

func TestInMemoryRepo_Delete(t *testing.T) {
	t.Run("error entity not found", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		err := sut.Delete(ctx, stubEntity.GetID(), stubEntity.GetETag())

		assert.Error(t, err)
		assert.Equal(t, http.StatusNotFound, arr.HTTPStatusFromError(err))
	})

	t.Run("error e-tag mismatch", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		err = sut.Delete(ctx, storedEntity.GetID(), "bar")

		assert.Error(t, err)
		assert.Equal(t, http.StatusPreconditionFailed, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		err = sut.Delete(ctx, storedEntity.GetID(), storedEntity.GetETag())

		assert.NoError(t, err)
		_, err = sut.GetByID(ctx, storedEntity.ID)
		assert.Error(t, err)
	})
}

func TestInMemoryRepo_List(t *testing.T) {
	ctxStub := context.Background()

	t.Run("error matching entity", func(t *testing.T) {
		t.Parallel()

		sut := repo.NewInMemoryRepo[testEntity]()

		entityStub := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		entityStub.ETag = "baz"
		_, err := sut.Create(ctxStub, entityStub)
		require.NoError(t, err)

		filterMock := new(mocks.Filter[testEntity])
		filterMock.EXPECT().
			Match(ctxStub, entityStub).
			Once().
			Return(false, assert.AnError)

		_, err = sut.List(ctxStub, filterMock)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		filterMock.AssertExpectations(t)
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		list, err := sut.List(ctx, testFilter{})

		assert.NoError(t, err)
		assert.Empty(t, list)
	})

	t.Run("deleted are filtered", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		now := time.Now()
		storedEntity.Entity.DeletedAt = &now

		_, err = sut.Update(ctx, storedEntity, storedEntity.GetETag())
		require.NoError(t, err)

		list, err := sut.List(ctx, testFilter{})

		assert.NoError(t, err)
		assert.Empty(t, list)
	})

	t.Run("success non-empty", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sut := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "foo",
		}
		stubEntity.ETag = "baz"

		storedEntity, err := sut.Create(ctx, stubEntity)
		require.NoError(t, err)

		list, err := sut.List(ctx, testFilter{})

		assert.NoError(t, err)
		assert.NotEmpty(t, list)
		assert.Contains(t, list, storedEntity)
	})
}
