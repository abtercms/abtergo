package redirect_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/logtest"
	"github.com/abtergo/abtergo/libs/model"
	repoMocks "github.com/abtergo/abtergo/mocks/libs/repo"
	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestService_Create(t *testing.T) {
	loggerStub, _ := logtest.NewDefaultLogger(t)
	ctxStub := context.Background()

	t.Run("validation error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()
		entityStub.Entity = model.Entity{}
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		s := redirect.NewService(repoMock, loggerStub)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()
		entityStub.Entity = model.Entity{}

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Create(ctxStub, mock.AnythingOfType("redirect.Redirect")).
			Once().
			Return(redirect.Redirect{}, assert.AnError)
		s := redirect.NewService(repoMock, loggerStub)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()
		entityStub.Entity = model.Entity{}

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Create(ctxStub, mock.AnythingOfType("redirect.Redirect")).
			Once().
			Return(entityStub, nil)
		s := redirect.NewService(repoMock, loggerStub)

		got, err := s.Create(ctxStub, entityStub)

		assert.NoError(t, err)
		got.Entity = entityStub.Entity
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub, _ := logtest.NewDefaultLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(assert.AnError)
		s := redirect.NewService(repoMock, loggerStub)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(nil)
		s := redirect.NewService(repoMock, loggerStub)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub, _ := logtest.NewDefaultLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			GetByID(ctxStub, entityStub.ID).
			Once().
			Return(redirect.Redirect{}, assert.AnError)
		s := redirect.NewService(repoMock, loggerStub)

		_, err := s.Get(ctxStub, entityStub.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			GetByID(ctxStub, entityStub.ID).
			Once().
			Return(entityStub, nil)
		s := redirect.NewService(repoMock, loggerStub)

		got, err := s.Get(ctxStub, entityStub.ID)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	loggerStub, _ := logtest.NewDefaultLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		filterStub := redirect.Filter{}

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(nil, assert.AnError)
		s := redirect.NewService(repoMock, loggerStub)

		_, err := s.List(ctxStub, filterStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		filterStub := redirect.Filter{}
		stubCollection := redirect.RandomRedirectList(1, 3)

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(stubCollection, nil)
		s := redirect.NewService(repoMock, loggerStub)

		got, err := s.List(ctxStub, filterStub)

		assert.NoError(t, err)
		assert.Equal(t, stubCollection, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	loggerStub, _ := logtest.NewDefaultLogger(t)
	ctxStub := context.Background()

	const (
		idStub   model.ID   = "foo"
		eTagStub model.ETag = "bar"
	)

	t.Run("validation error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()
		entityStub.Entity = model.Entity{}
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		s := redirect.NewService(repoMock, loggerStub)

		_, err := s.Update(ctxStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("redirect.Redirect"), eTagStub).
			Once().
			Return(redirect.Redirect{}, assert.AnError)
		s := redirect.NewService(repoMock, loggerStub)

		_, err := s.Update(ctxStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := redirect.RandomRedirect()
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[redirect.Redirect])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("redirect.Redirect"), eTagStub).
			Once().
			Return(entityStub, nil)
		s := redirect.NewService(repoMock, loggerStub)

		got, err := s.Update(ctxStub, entityStub, eTagStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}
