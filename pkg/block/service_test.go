package block_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	repoMocks "github.com/abtergo/abtergo/mocks/libs/repo"
	"github.com/abtergo/abtergo/pkg/block"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("id provided error", func(t *testing.T) {
		entityStub := block.RandomBlock(false)

		repoMock := new(repoMocks.Repository[block.Block])

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Create(ctxStub, entityStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := block.RandomBlock(true)
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Once().
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Create(ctxStub, entityStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock(true)

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Once().
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.Create(ctxStub, entityStub)

		assert.NoError(t, err)
		got.Entity = entityStub.Entity
		assert.Equal(t, entityStub, got)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock(false)

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(nil)

		s := block.NewService(loggerStub, repoMock)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)
		require.NoError(t, err)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock(false)

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Retrieve(ctxStub, entityStub.ID).
			Once().
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.Get(ctxStub, entityStub.ID)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}

func TestService_List(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		filterStub := block.Filter{}
		stubCollection := block.RandomBlockList(1, 3)

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(stubCollection, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.List(ctxStub, filterStub)
		require.NoError(t, err)

		assert.Equal(t, stubCollection, got)
	})
}

func TestService_Update(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	const (
		idStub   = "foo"
		etagStub = "bar"
	)

	t.Run("id mismatch error", func(t *testing.T) {
		entityStub := block.RandomBlock(false)

		repoMock := new(repoMocks.Repository[block.Block])

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := block.RandomBlock(true)
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Update(ctxStub, entityStub, etagStub).
			Once().
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock(false)
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Update(ctxStub, entityStub, etagStub).
			Once().
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}
