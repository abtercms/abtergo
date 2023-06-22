package template_test

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
	"github.com/abtergo/abtergo/pkg/template"
)

func TestService_Create(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("validation error", func(t *testing.T) {
		entityStub := template.RandomTemplate(true)
		entityStub.ID = ""
		entityStub.Website = ""

		repoMock := new(repoMocks.Repository[template.Template])
		s := template.NewService(repoMock, loggerStub)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		entityStub := template.RandomTemplate(true)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			Create(ctxStub, mock.AnythingOfType("template.Template")).
			Once().
			Return(template.Template{}, assert.AnError)
		s := template.NewService(repoMock, loggerStub)

		_, err := s.Create(ctxStub, entityStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := template.RandomTemplate(true)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			Create(ctxStub, mock.AnythingOfType("template.Template")).
			Once().
			Return(entityStub, nil)
		s := template.NewService(repoMock, loggerStub)

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

	t.Run("repo error", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(assert.AnError)
		s := template.NewService(repoMock, loggerStub)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID, entityStub.ETag).
			Once().
			Return(nil)
		s := template.NewService(repoMock, loggerStub)

		err := s.Delete(ctxStub, entityStub.ID, entityStub.ETag)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			GetByID(ctxStub, entityStub.ID).
			Once().
			Return(template.Template{}, assert.AnError)
		s := template.NewService(repoMock, loggerStub)

		_, err := s.Get(ctxStub, entityStub.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			GetByID(ctxStub, entityStub.ID).
			Once().
			Return(entityStub, nil)
		s := template.NewService(repoMock, loggerStub)

		got, err := s.Get(ctxStub, entityStub.ID)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	loggerStub := zaptest.NewLogger(t)
	ctxStub := context.Background()

	t.Run("repo error", func(t *testing.T) {
		filterStub := template.Filter{}

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(nil, assert.AnError)
		s := template.NewService(repoMock, loggerStub)

		_, err := s.List(ctxStub, filterStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		filterStub := template.Filter{}
		stubCollection := template.RandomTemplateList(1, 3)

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Once().
			Return(stubCollection, nil)
		s := template.NewService(repoMock, loggerStub)

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

	t.Run("id mismatch error", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)

		repoMock := new(repoMocks.Repository[template.Template])

		s := template.NewService(repoMock, loggerStub)
		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := template.RandomTemplate(true)
		entityStub.Website = ""
		entityStub.ID = ""

		repoMock := new(repoMocks.Repository[template.Template])
		s := template.NewService(repoMock, loggerStub)

		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFromError(err))
		repoMock.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("template.Template"), eTagStub).
			Once().
			Return(template.Template{}, assert.AnError)
		s := template.NewService(repoMock, loggerStub)

		_, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		repoMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		entityStub := template.RandomTemplate(false)
		entityStub.ID = idStub

		repoMock := new(repoMocks.Repository[template.Template])
		repoMock.EXPECT().
			Update(ctxStub, mock.AnythingOfType("template.Template"), eTagStub).
			Once().
			Return(entityStub, nil)
		s := template.NewService(repoMock, loggerStub)

		got, err := s.Update(ctxStub, idStub, entityStub, eTagStub)

		assert.NoError(t, err)
		assert.Equal(t, entityStub, got)
		repoMock.AssertExpectations(t)
	})
}
