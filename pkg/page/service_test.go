package page_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	mocks2 "github.com/abtergo/abtergo/mocks/pkg/page"
	"github.com/abtergo/abtergo/pkg/page"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("id provided error", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(mocks2.Repo)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Create(ctxStub, entityStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := page.RandomPage(true)
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Return(entityStub, nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Create(ctxStub, entityStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(true)

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Return(entityStub, nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Create(ctxStub, entityStub)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID).
			Return(nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		err := s.Delete(ctxStub, entityStub.ID)
		require.NoError(t, err)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(false)

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Retrieve(ctxStub, entityStub.ID).
			Return(entityStub, nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Get(ctxStub, entityStub.ID)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}

func TestService_List(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		filterStub := page.Filter{}
		stubCollection := page.RandomPageList(1, 3)

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Return(stubCollection, nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

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
		entityStub := page.RandomPage(false)

		repoMock := new(mocks2.Repo)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := page.RandomPage(true)
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Update(ctxStub, idStub, entityStub, etagStub).
			Return(entityStub, nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
	})

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage(false)
		entityStub.ID = idStub

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Update(ctxStub, idStub, entityStub, etagStub).
			Return(entityStub, nil)

		updaterMock := new(mocks2.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}
