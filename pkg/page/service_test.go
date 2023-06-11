package page_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	repoMocks "github.com/abtergo/abtergo/mocks/libs/repo"
	mocks "github.com/abtergo/abtergo/mocks/pkg/page"
	"github.com/abtergo/abtergo/pkg/page"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("id provided error", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(repoMocks.Repository[page.Page])

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := page.RandomPage(true)
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[page.Page])

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(true)

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Create(ctxStub, entityStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Return(nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, entityStub.ID).
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

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
		filterStub := page.Filter{}
		stubCollection := page.RandomPageList(1, 3)

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Return(stubCollection, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

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
		eTagStub = "bar"
	)

	t.Run("id mismatch error", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(repoMocks.Repository[page.Page])

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := page.RandomPage(true)
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(repoMocks.Repository[page.Page])

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(false)
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Update(ctxStub, entityStub, eTagStub).
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}
