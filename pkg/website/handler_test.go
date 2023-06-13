package website_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/fib"
	mocks "github.com/abtergo/abtergo/mocks/pkg/website"
	"github.com/abtergo/abtergo/pkg/website"
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
		deps.serviceMock.EXPECT().
			Get(mock.Anything, baseURLStub, pathStub).
			Once().
			Return("", assert.AnError)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

type handlerDeps struct {
	loggerStub  *zap.Logger
	serviceMock *mocks.Service
}

func (hd handlerDeps) AssertExpectations(t *testing.T) {
	t.Helper()

	hd.serviceMock.AssertExpectations(t)
}

func setupHandlerMocks(t *testing.T) (*fiber.App, handlerDeps) {
	t.Helper()

	loggerStub := zaptest.NewLogger(t)
	serviceMock := &mocks.Service{}
	handler := website.NewHandler(serviceMock, loggerStub)
	errorHandler := fib.NewErrorHandler(loggerStub)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.Handle,
		ReadTimeout:  time.Hour,
		WriteTimeout: time.Hour,
		IdleTimeout:  time.Hour,
	})
	handler.AddRoutes(app)

	return app, handlerDeps{loggerStub: loggerStub, serviceMock: serviceMock}
}
