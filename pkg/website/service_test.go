package website_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mocks "github.com/abtergo/abtergo/mocks/pkg/website"
	"github.com/abtergo/abtergo/pkg/page"
	templatePkg "github.com/abtergo/abtergo/pkg/template"
	"github.com/abtergo/abtergo/pkg/website"
)

func TestService_Get(t *testing.T) {
	ctxStub := context.Background()

	t.Run("error retrieving content", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		pageStub := page.RandomPage()

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(nil, assert.AnError)

		got, err := s.Get(ctxStub, pageStub.Website, pageStub.Path)

		assert.Error(t, err)
		assert.Empty(t, got)
		assert.ErrorIs(t, err, assert.AnError)
		deps.AssertExpectations(t)
	})

	t.Run("error retrieving template", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		pageStub := page.RandomPage()

		contentStub := new(mocks.Content)
		contentStub.
			EXPECT().
			Render().
			Once().
			Return(pageStub.Body, nil)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(contentStub, nil)

		deps.templateRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(templatePkg.Template{}, assert.AnError)

		got, err := s.Get(ctxStub, pageStub.Website, pageStub.Path)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, got)
		deps.AssertExpectations(t)
	})

	t.Run("error rendering content", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		pageStub := page.RandomPage()
		template := templatePkg.RandomTemplate(false)

		contentStub := new(mocks.Content)
		contentStub.
			EXPECT().
			Render().
			Once().
			Return(pageStub.Body, nil)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(contentStub, nil)

		deps.templateRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(template, nil)

		deps.rendererMock.
			EXPECT().
			Render(template.Body, pageStub.Body).
			Once().
			Return("", assert.AnError)

		got, err := s.Get(ctxStub, pageStub.Website, pageStub.Path)

		assert.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, got)
		deps.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		expectedBody := "foo"

		s, deps := setupServiceMocks(t)

		pageStub := page.RandomPage()
		template := templatePkg.RandomTemplate(false)

		contentStub := new(mocks.Content)
		contentStub.
			EXPECT().
			Render().
			Once().
			Return(pageStub.Body, nil)

		deps.contentRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(contentStub, nil)

		deps.templateRetriever.
			EXPECT().
			Retrieve(mock.Anything, pageStub.Website, pageStub.Path).
			Once().
			Return(template, nil)

		deps.rendererMock.
			EXPECT().
			Render(template.Body, pageStub.Body).
			Once().
			Return(expectedBody, nil)

		got, err := s.Get(ctxStub, pageStub.Website, pageStub.Path)

		assert.NoError(t, err)
		assert.Equal(t, expectedBody, got)
		deps.AssertExpectations(t)
	})
}

type serviceDeps struct {
	contentRetriever  *mocks.ContentRetriever
	templateRetriever *mocks.TemplateRetriever
	rendererMock      *mocks.Renderer
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
		rendererMock:      &mocks.Renderer{},
	}

	serviceMock := website.NewService(sd.contentRetriever, sd.templateRetriever, sd.rendererMock)

	return serviceMock, sd
}
