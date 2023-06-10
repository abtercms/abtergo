package renderer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	blockMocks "github.com/abtergo/abtergo/mocks/pkg/block"
	pageMocks "github.com/abtergo/abtergo/mocks/pkg/page"
	redirectMocks "github.com/abtergo/abtergo/mocks/pkg/redirect"
	mocks "github.com/abtergo/abtergo/mocks/pkg/renderer"
	templateMocks "github.com/abtergo/abtergo/mocks/pkg/template"
	"github.com/abtergo/abtergo/pkg/page"
	"github.com/abtergo/abtergo/pkg/renderer"
)

func TestService_Get(t *testing.T) {
	ctxStub := context.Background()

	t.Run("error retrieving page", func(t *testing.T) {
		s, deps := setupServiceMocks(t)

		entity := page.RandomPage(false)

		deps.pageRepoMock.
			EXPECT().
			RetrieveByWebsiteAndPage(mock.Anything, entity.Website, entity.Path).
			Once().
			Return(page.Page{}, assert.AnError)

		got, err := s.Get(ctxStub, entity.Website, entity.Path)
		require.Error(t, err)

		assert.Empty(t, got)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("success", func(t *testing.T) {
		expectedBody := "foo"
		expectedError := assert.AnError

		s, deps := setupServiceMocks(t)

		entity := page.RandomPage(false)

		deps.pageRepoMock.
			EXPECT().
			RetrieveByWebsiteAndPage(mock.Anything, entity.Website, entity.Path).
			Once().
			Return(entity, nil)

		deps.rendererMock.
			EXPECT().
			Render(entity.Body).
			Once().
			Return(expectedBody, expectedError)

		got, err := s.Get(ctxStub, entity.Website, entity.Path)

		assert.Equal(t, expectedBody, got)
		assert.Equal(t, err, expectedError)
	})
}

type serviceDeps struct {
	rendererMock     *mocks.Renderer
	pageRepoMock     *pageMocks.Repo
	templateRepoMock *templateMocks.Repo
	blockRepoMock    *blockMocks.Repo
	redirectRepoMock *redirectMocks.Repo
}

func (sd serviceDeps) AssertExpectations(t *testing.T) {
	sd.pageRepoMock.AssertExpectations(t)
	sd.templateRepoMock.AssertExpectations(t)
	sd.blockRepoMock.AssertExpectations(t)
	sd.redirectRepoMock.AssertExpectations(t)
}

func setupServiceMocks(t *testing.T) (renderer.Service, serviceDeps) {
	sd := serviceDeps{
		rendererMock:     &mocks.Renderer{},
		pageRepoMock:     &pageMocks.Repo{},
		templateRepoMock: &templateMocks.Repo{},
		blockRepoMock:    &blockMocks.Repo{},
		redirectRepoMock: &redirectMocks.Repo{},
	}

	serviceMock := renderer.NewService(sd.rendererMock, sd.pageRepoMock, sd.templateRepoMock, sd.blockRepoMock, sd.redirectRepoMock)

	return serviceMock, sd
}
