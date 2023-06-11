package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
	mocks "github.com/abtergo/abtergo/mocks/libs/repo"
)

type testEntity struct {
	model.Entity

	Foo string `json:"foo"`
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
	ctx := context.Background()
	repo := repo.NewInMemoryRepo[testEntity]()

	stubEntity := testEntity{
		Entity: model.NewEntity(),
		Foo:    "bar",
	}.AsNew().(testEntity)

	storedEntity, err := repo.Create(ctx, stubEntity)
	require.NoError(t, err)

	require.Equal(t, stubEntity, storedEntity)
}

func TestInMemoryRepo_Retrieve(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "bar",
		}

		storedEntity, err := repo.Create(ctx, stubEntity)
		require.NoError(t, err)

		retrievedEntity, err := repo.Retrieve(ctx, storedEntity.ID)
		require.NoError(t, err)

		require.Equal(t, retrievedEntity, storedEntity)

		require.Equal(t, stubEntity.AsNew(), storedEntity.AsNew())
	})
}

func TestInMemoryRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "bar",
		}

		storedEntity, err := repo.Create(ctx, stubEntity)
		require.NoError(t, err)

		stubEntity2 := testEntity{
			Entity: storedEntity.Entity,
			Foo:    "baz",
		}

		updatedEntity, err := repo.Update(ctx, stubEntity2, storedEntity.GetETag())
		require.NoError(t, err)

		retrievedEntity, err := repo.Retrieve(ctx, updatedEntity.ID)
		require.NoError(t, err)

		require.NotEqual(t, retrievedEntity, storedEntity)
		require.Equal(t, retrievedEntity, updatedEntity)
	})
}

func TestInMemoryRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "bar",
		}

		storedEntity, err := repo.Create(ctx, stubEntity)
		require.NoError(t, err)

		err = repo.Delete(ctx, storedEntity.GetID(), storedEntity.GetETag())
		assert.NoError(t, err)

		_, err = repo.Retrieve(ctx, storedEntity.ID)
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
			Foo:    "bar",
		}
		_, err := sut.Create(ctxStub, entityStub)
		require.NoError(t, err)

		filterMock := new(mocks.Filter[testEntity])
		filterMock.EXPECT().Match(ctxStub, entityStub).Return(false, assert.AnError)

		_, err = sut.List(ctxStub, filterMock)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		filterMock.AssertExpectations(t)
	})

	t.Run("success empty", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := repo.NewInMemoryRepo[testEntity]()

		list, err := repo.List(ctx, testFilter{})

		assert.NoError(t, err)
		assert.Empty(t, list)
	})

	t.Run("deleted are filtered", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "bar",
		}

		storedEntity, err := repo.Create(ctx, stubEntity)
		require.NoError(t, err)

		now := time.Now()
		storedEntity.Entity.DeletedAt = &now

		_, err = repo.Update(ctx, storedEntity, storedEntity.GetETag())
		require.NoError(t, err)

		list, err := repo.List(ctx, testFilter{})

		assert.NoError(t, err)
		assert.Empty(t, list)
	})

	t.Run("success non-empty", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		repo := repo.NewInMemoryRepo[testEntity]()

		stubEntity := testEntity{
			Entity: model.NewEntity(),
			Foo:    "bar",
		}

		storedEntity, err := repo.Create(ctx, stubEntity)
		require.NoError(t, err)

		list, err := repo.List(ctx, testFilter{})

		assert.NoError(t, err)
		assert.NotEmpty(t, list)
		assert.Contains(t, list, storedEntity)
	})
}
