package website_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/abtergo/abtergo/libs/model"
	templMocks "github.com/abtergo/abtergo/mocks/libs/templ"
	mocks "github.com/abtergo/abtergo/mocks/pkg/website"
	"github.com/abtergo/abtergo/pkg/page"
	templatePkg "github.com/abtergo/abtergo/pkg/template"
	"github.com/abtergo/abtergo/pkg/website"
)

func TestService_Get(t *testing.T) {
	ctxStub := context.Background()

	t.Run("error retrieving content", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		randomPage := page.RandomPage()
		key := model.KeyFromStrings(randomPage.Website, randomPage.Path)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, key).
			Once().
			Return(nil, assert.AnError)

		got, err := s.Get(ctxStub, randomPage.Website, randomPage.Path)

		assert.Error(t, err)
		assert.Empty(t, got)
		assert.ErrorIs(t, err, assert.AnError)
		deps.AssertExpectations(t)
	})

	t.Run("error retrieving template", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		randomPage := page.RandomPage()
		key := model.KeyFromStrings(randomPage.Website, randomPage.Path)

		contentMock := new(templMocks.CacheableContent)
		contentMock.EXPECT().Render().Once().Return(randomPage.Body, nil)
		contentMock.EXPECT().GetContext().Once().Return(nil)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, key).
			Once().
			Return(contentMock, nil)

		deps.templateRetriever.
			EXPECT().
			Retrieve(mock.Anything, randomPage.Website, randomPage.Path).
			Once().
			Return(nil, assert.AnError)

		got, err := s.Get(ctxStub, randomPage.Website, randomPage.Path)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, got)
		deps.AssertExpectations(t)
	})

	t.Run("error rendering content", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		randomPage := page.RandomPage()
		randomTemplate := templatePkg.RandomTemplate()
		key := model.KeyFromStrings(randomPage.Website, randomPage.Path)

		contentMock := new(templMocks.CacheableContent)
		contentMock.EXPECT().Render().Once().Return(randomPage.Body, nil)
		contentMock.EXPECT().GetContext().Once().Return(nil)

		templateMock := new(templMocks.CacheableContent)
		templateMock.EXPECT().Render().Once().Return(randomTemplate.Body, nil)
		templateMock.EXPECT().GetContext().Once().Return(nil)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, key).
			Once().
			Return(contentMock, nil)

		deps.templateRetriever.
			EXPECT().
			Retrieve(mock.Anything, randomPage.Website, randomPage.Path).
			Once().
			Return(templateMock, nil)

		deps.rendererMock.
			EXPECT().
			RenderInLayout(randomPage.Body, randomTemplate.Body).
			Once().
			Return("", assert.AnError)

		got, err := s.Get(ctxStub, randomPage.Website, randomPage.Path)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, got)
		deps.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		expectedBody := "foo"

		s, deps := setupServiceMocks(t)

		randomPage := page.RandomPage()
		randomTemplate := templatePkg.RandomTemplate()
		key := model.KeyFromStrings(randomPage.Website, randomPage.Path)

		contentMock := new(templMocks.CacheableContent)
		contentMock.EXPECT().Render().Once().Return(randomPage.Body, nil)
		contentMock.EXPECT().GetContext().Once().Return(nil)

		templateMock := new(templMocks.CacheableContent)
		templateMock.EXPECT().Render().Once().Return(randomTemplate.Body, nil)
		templateMock.EXPECT().GetContext().Once().Return(nil)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, key).
			Once().
			Return(contentMock, nil)

		deps.templateRetriever.
			EXPECT().
			Retrieve(mock.Anything, randomPage.Website, randomPage.Path).
			Once().
			Return(templateMock, nil)

		deps.rendererMock.
			EXPECT().
			RenderInLayout(randomPage.Body, randomTemplate.Body).
			Once().
			Return(expectedBody, nil)

		got, err := s.Get(ctxStub, randomPage.Website, randomPage.Path)

		assert.NoError(t, err)
		assert.Equal(t, expectedBody, got)
		deps.AssertExpectations(t)
	})
}

type serviceDeps struct {
	contentRetriever  *mocks.ContentRetriever
	templateRetriever *mocks.TemplateRetriever
	rendererMock      *templMocks.Renderer
}

func (sd serviceDeps) AssertExpectations(t *testing.T) {
	t.Helper()

	sd.contentRetriever.AssertExpectations(t)
	sd.templateRetriever.AssertExpectations(t)
	sd.rendererMock.AssertExpectations(t)
}

func setupServiceMocks(t *testing.T) (website.Service, serviceDeps) {
	t.Helper()

	sd := serviceDeps{
		contentRetriever:  &mocks.ContentRetriever{},
		templateRetriever: &mocks.TemplateRetriever{},
		rendererMock:      &templMocks.Renderer{},
	}

	serviceMock := website.NewService(sd.contentRetriever, sd.templateRetriever, sd.rendererMock)

	return serviceMock, sd
}
