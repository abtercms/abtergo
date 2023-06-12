package redirect_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	repoMocks "github.com/abtergo/abtergo/mocks/libs/repo"
	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("validation error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(true)
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[redirect.Redirect])

		s := redirect.NewService(loggerStub, repoMock)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(true)

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Create(ctxStub, mock.AnythingOfType("redirect.Redirect")).
			Once().
			Return(entityStub, nil)

		s := redirect.NewService(loggerStub, repoMock)

		got, err := s.Create(ctxStub, entityStub)

		assert.NoError(t, err)
		got.Entity = entityStub.Entity
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(false)

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(nil)

		s := redirect.NewService(loggerStub, repoMock)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(false)

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Retrieve(ctxStub, entityStub.ID).
			Once().
			Return(entityStub, nil)

		s := redirect.NewService(loggerStub, repoMock)

		got, err := s.Get(ctxStub, entityStub.ID)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		filterStub := redirect.Filter{}
		stubCollection := redirect.RandomRedirectList(1, 3)

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(stubCollection, nil)

		s := redirect.NewService(loggerStub, repoMock)

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
		idStub   = "foo"
		etagStub = "bar"
	)

	t.Run("id mismatch error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(false)

		repoMock := new(repoMocks.Repository[redirect.Redirect])

		s := redirect.NewService(loggerStub, repoMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(true)
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(repoMocks.Repository[redirect.Redirect])

		s := redirect.NewService(loggerStub, repoMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect(false)
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("redirect.Redirect"), etagStub).
			Once().
			Return(entityStub, nil)

		s := redirect.NewService(loggerStub, repoMock)

		got, err := s.Update(ctxStub, idStub, entityStub, etagStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}
