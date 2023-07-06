package block_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	repoMocks "github.com/abtergo/abtergo/mocks/libs/repo"
	"github.com/abtergo/abtergo/pkg/block"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("validation error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.Entity = model.Entity{}
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[block.Block])
		s := block.NewService(repoMock, loggerStub)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.Entity = model.Entity{}
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[block.Block])
		s := block.NewService(repoMock, loggerStub)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.Entity = model.Entity{}

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Create(ctxStub, mock.AnythingOfType("block.Block")).
			Once().
			Return(entityStub, nil)

		s := block.NewService(repoMock, loggerStub)

		got, err := s.Create(ctxStub, entityStub)

		assert.NoError(t, err)
		got.Entity = entityStub.Entity
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		entityStub := block.RandomBlock()

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			GetByID(ctxStub, entityStub.ID).
			Once().
			Return(block.Block{}, assert.AnError)
		s := block.NewService(repoMock, loggerStub)

		_, err := s.GetByID(ctxStub, entityStub.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock()

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			GetByID(ctxStub, entityStub.ID).
			Once().
			Return(entityStub, nil)
		s := block.NewService(repoMock, loggerStub)

		got, err := s.GetByID(ctxStub, entityStub.ID)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		filterStub := block.Filter{}

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(nil, assert.AnError)
		s := block.NewService(repoMock, loggerStub)

		_, err := s.List(ctxStub, filterStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		filterStub := block.Filter{}
		stubCollection := block.RandomBlockList(1, 3)

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(stubCollection, nil)
		s := block.NewService(repoMock, loggerStub)

		got, err := s.List(ctxStub, filterStub)

		assert.NoError(t, err)
		assert.Equal(t, stubCollection, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	const (
		idStub   model.ID   = "foo"
		eTagStub model.ETag = "bar"
	)

	t.Run("validation error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.Entity = model.Entity{}
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(repoMocks.Repository[block.Block])
		s := block.NewService(repoMock, loggerStub)

		_, err := s.Update(ctxStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("block.Block"), eTagStub).
			Once().
			Return(block.Block{}, assert.AnError)
		s := block.NewService(repoMock, loggerStub)

		_, err := s.Update(ctxStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("block.Block"), eTagStub).
			Once().
			Return(entityStub, nil)
		s := block.NewService(repoMock, loggerStub)

		got, err := s.Update(ctxStub, entityStub, eTagStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		entityStub := block.RandomBlock()

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(assert.AnError)
		s := block.NewService(repoMock, loggerStub)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock()

		repoMock := new(repoMocks.Repository[block.Block])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(nil)
		s := block.NewService(repoMock, loggerStub)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}
