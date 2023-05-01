package renderer_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/fib"
	mocks "github.com/abtergo/abtergo/mocks/libs/ablog"
	mocks2 "github.com/abtergo/abtergo/mocks/pkg/renderer"
	"github.com/abtergo/abtergo/pkg/renderer"
)

func TestHandler_AddRoutes(t *testing.T) {
	const baseURLStub = "http://example.com"

	t.Run("Catch All", func(t *testing.T) {
		// Expectations

		// Stubs
		pathStub := "/does-not-exist"

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+pathStub, nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Prepare Test
		app, deps := setupHandlerMocks(t)
		deps.loggerMock.EXPECT().Infow(mock.Anything, "Method", http.MethodGet, "Path", pathStub).Once()
		deps.serviceMock.EXPECT().Get(mock.Anything, baseURLStub, pathStub).Return("", assert.AnError)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

type handlerDeps struct {
	loggerMock  *mocks.Logger
	serviceMock *mocks2.Service
}

func (hd handlerDeps) AssertExpectations(t *testing.T) {
	hd.loggerMock.AssertExpectations(t)
}

func setupHandlerMocks(t *testing.T) (*fiber.App, handlerDeps) {
	loggerMock := &mocks.Logger{}
	serviceMock := &mocks2.Service{}
	handler := renderer.NewHandler(loggerMock, serviceMock)

	app := fiber.New(fiber.Config{
		ErrorHandler: fib.ErrorHandler,
	})
	handler.AddRoutes(app)

	return app, handlerDeps{loggerMock: loggerMock, serviceMock: serviceMock}
}
