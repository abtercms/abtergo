package page_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	repoMocks "github.com/abtergo/abtergo/mocks/libs/repo"
	mocks "github.com/abtergo/abtergo/mocks/pkg/page"
	"github.com/abtergo/abtergo/pkg/page"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("id provided error", func(t *testing.T) {
		entityStub := page.RandomPage()

		repoMock := new(repoMocks.Repository[page.Page])
		updaterMock := new(mocks.Updater)

		sut := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := sut.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := page.RandomPage()
		entityStub.Entity = model.Entity{}
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[page.Page])
		updaterMock := new(mocks.Updater)

		sut := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := sut.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage()
		entityStub.Entity = model.Entity{}

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Once().
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Create(ctxStub, entityStub)

		assert.NoError(t, err)
		got.Entity = entityStub.Entity
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage()

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage()

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, entityStub.ID).
			Once().
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Get(ctxStub, entityStub.ID)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
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
			Once().
			Return(stubCollection, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.List(ctxStub, filterStub)

		assert.NoError(t, err)
		assert.Equal(t, stubCollection, got)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	const (
		idStub   = "foo"
		eTagStub = "bar23"
	)

	t.Run("id mismatch error", func(t *testing.T) {
		entityStub := page.RandomPage()

		repoMock := new(repoMocks.Repository[page.Page])

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := page.RandomPage()
		entityStub.Entity = model.Entity{}
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(repoMocks.Repository[page.Page])

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := page.RandomPage()
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Update(ctxStub, entityStub, eTagStub).
			Once().
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})
}

func Test_service_Transition(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	const (
		idStub   = "foo"
		eTagStub = "bar23"
	)

	t.Run("error retrieving page", func(t *testing.T) {
		t.Parallel()

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, idStub).
			Once().
			Return(page.Page{}, assert.AnError)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Transition(ctxStub, idStub, page.Activate, eTagStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("e-tag mismatch", func(t *testing.T) {
		t.Parallel()

		entityStub := page.RandomPage()
		entityStub.Status = page.StatusActive

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, idStub).
			Once().
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Transition(ctxStub, idStub, page.Activate, eTagStub)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid e-tag found")
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("error in transaction", func(t *testing.T) {
		t.Parallel()

		entityStub := page.RandomPage()
		entityStub.Status = page.StatusActive
		entityStub.ETag = eTagStub

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, idStub).
			Once().
			Return(entityStub, nil)

		updaterMock := new(mocks.Updater)
		updaterMock.EXPECT().
			Transition(page.StatusActive, page.Activate).
			Once().
			Return(page.StatusActive, assert.AnError)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Transition(ctxStub, idStub, page.Activate, eTagStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("error in updating page", func(t *testing.T) {
		t.Parallel()

		entityStub := page.RandomPage()
		entityStub.Status = page.StatusInactive
		entityStub.ETag = eTagStub

		updatedStub := entityStub.Clone().(page.Page)
		updatedStub.Status = page.StatusActive

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, idStub).
			Once().
			Return(entityStub, nil)
		repoMock.EXPECT().
			Update(ctxStub, updatedStub, eTagStub).
			Once().
			Return(page.Page{}, assert.AnError)

		updaterMock := new(mocks.Updater)
		updaterMock.EXPECT().
			Transition(page.StatusInactive, page.Activate).
			Once().
			Return(page.StatusActive, nil)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		_, err := s.Transition(ctxStub, idStub, page.Activate, eTagStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		entityStub := page.RandomPage()
		entityStub.Status = page.StatusInactive
		entityStub.ETag = eTagStub

		updatedStub := entityStub.Clone().(page.Page)
		updatedStub.Status = page.StatusActive

		repoMock := new(repoMocks.Repository[page.Page])
		repoMock.EXPECT().
			Retrieve(ctxStub, idStub).
			Once().
			Return(entityStub, nil)
		repoMock.EXPECT().
			Update(ctxStub, updatedStub, eTagStub).
			Once().
			Return(updatedStub, nil)

		updaterMock := new(mocks.Updater)
		updaterMock.EXPECT().
			Transition(page.StatusInactive, page.Activate).
			Once().
			Return(page.StatusActive, nil)

		s := page.NewService(loggerStub, repoMock, updaterMock)

		got, err := s.Transition(ctxStub, idStub, page.Activate, eTagStub)

		assert.NoError(t, err)
		assert.Equal(t, updatedStub, got)
		repoMock.AssertExpectations(t)
		updaterMock.AssertExpectations(t)
	})
}
