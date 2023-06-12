package redirect_test

import (
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/fib"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/problem"
	"github.com/abtergo/abtergo/libs/util"
	mocks "github.com/abtergo/abtergo/mocks/pkg/redirect"
	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestHandler_AddApiRoutes(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("undefined route results in 404", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusNotFound

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/does-not-exist", nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Prepare Test
		app, _ := setupHandlerMocks(t)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)
	})
}

func TestHandler_Post(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("error parsing body", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusBadRequest

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks

		// Request
		reqBody := util.DataToReaderHelper(t, `"foo"`)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/redirects", reqBody)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		deps.AssertExpectations(t)
	})

	t.Run("error persisting entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusConflict
		expectedRedirect := redirect.RandomRedirect(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Create(mock.Anything, expectedRedirect).
			Once().
			Return(redirect.Redirect{}, arr.Wrap(arr.ResourceIsOutdated, assert.AnError, "foo"))

		// Request
		reqBody := util.DataToReaderHelper(t, expectedRedirect)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/redirects", reqBody)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		deps.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusCreated
		expectedRedirect := redirect.RandomRedirect(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Create(mock.Anything, expectedRedirect).
			Once().
			Return(expectedRedirect, nil)

		// Request
		reqBody := util.DataToReaderHelper(t, expectedRedirect)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/redirects", reqBody)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual redirect.Redirect
		util.ParseResponseHelper(t, resp, &actual)
		actual.Entity = expectedRedirect.Entity
		assert.Equal(t, expectedRedirect, actual)

		deps.AssertExpectations(t)
	})
}

func TestHandler_List(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("error retrieving collection", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusPreconditionFailed

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			List(mock.Anything, redirect.Filter{}).
			Once().
			Return(nil, arr.Wrap(arr.ETagMismatch, assert.AnError, "foo"))

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/redirects", nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual problem.Problem
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedStatusCode, actual.Status)

		deps.serviceMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusOK
		expectedRedirects := redirect.RandomRedirectList(5, 5)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			List(mock.Anything, redirect.Filter{}).
			Once().
			Return(expectedRedirects, nil)

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/redirects", nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual []redirect.Redirect
		util.ParseResponseHelper(t, resp, &actual)
		assert.Len(t, actual, 5)

		actual[0].Entity = model.Entity{}
		expectedRedirects[0].Entity = model.Entity{}
		assert.Equal(t, expectedRedirects[0], actual[0])

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Get(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("error retrieving entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusConflict
		expectedRedirect := redirect.RandomRedirect(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Get(mock.Anything, expectedRedirect.ID).
			Once().
			Return(redirect.Redirect{}, arr.Wrap(arr.ResourceIsOutdated, assert.AnError, "foo"))

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/redirects/"+expectedRedirect.ID, nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual problem.Problem
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedStatusCode, actual.Status)

		deps.serviceMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusOK
		expectedRedirect := redirect.RandomRedirect(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Get(mock.Anything, expectedRedirect.ID).
			Once().
			Return(expectedRedirect, nil)

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/redirects/"+expectedRedirect.ID, nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual redirect.Redirect
		util.ParseResponseHelper(t, resp, &actual)
		actual.Entity = expectedRedirect.Entity
		assert.Equal(t, expectedRedirect, actual)

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Put(t *testing.T) {
	const (
		baseURLStub      = "https://example.com"
		previousEtagStub = "foo"
	)

	t.Run("error parsing payload", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusBadRequest
		expectedRedirect := redirect.RandomRedirect(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks

		// Request
		reqBody := util.DataToReaderHelper(t, `"foo"`)
		req := httptest.NewRequest(fiber.MethodPut, baseURLStub+"/redirects/"+expectedRedirect.ID, reqBody)
		req.Header.Set(fiber.HeaderETag, previousEtagStub)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual problem.Problem
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedStatusCode, actual.Status)

		deps.serviceMock.AssertExpectations(t)
	})

	t.Run("error updating entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusBadGateway
		expectedRedirect := redirect.RandomRedirect(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Update(mock.Anything, expectedRedirect.ID, expectedRedirect, previousEtagStub).
			Once().
			Return(redirect.Redirect{}, arr.Wrap(arr.UpstreamServiceUnavailable, assert.AnError, "foo"))

		// Request
		reqBody := util.DataToReaderHelper(t, expectedRedirect)
		req := httptest.NewRequest(fiber.MethodPut, baseURLStub+"/redirects/"+expectedRedirect.ID, reqBody)
		req.Header.Set(fiber.HeaderETag, previousEtagStub)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual problem.Problem
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedStatusCode, actual.Status)

		deps.serviceMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusOK
		expectedRedirect := redirect.RandomRedirect(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Update(mock.Anything, expectedRedirect.ID, expectedRedirect, previousEtagStub).
			Once().
			Return(expectedRedirect, nil)

		// Request
		reqBody := util.DataToReaderHelper(t, expectedRedirect)
		req := httptest.NewRequest(fiber.MethodPut, baseURLStub+"/redirects/"+expectedRedirect.ID, reqBody)
		req.Header.Set(fiber.HeaderETag, previousEtagStub)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual redirect.Redirect
		util.ParseResponseHelper(t, resp, &actual)
		actual.Entity = expectedRedirect.Entity
		assert.Equal(t, expectedRedirect, actual)

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Delete(t *testing.T) {
	const (
		baseURLStub      = "https://example.com"
		previousEtagStub = "foo"
	)

	t.Run("error deleting entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusTooManyRequests
		expectedRedirect := redirect.RandomRedirect(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Delete(mock.Anything, expectedRedirect.ID, previousEtagStub).
			Once().
			Return(arr.Wrap(arr.UpstreamServiceBusy, assert.AnError, "foo"))

		// Request
		req := httptest.NewRequest(fiber.MethodDelete, baseURLStub+"/redirects/"+expectedRedirect.ID, nil)
		req.Header.Set(fiber.HeaderETag, previousEtagStub)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		deps.serviceMock.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusNoContent
		expectedRedirect := redirect.RandomRedirect(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Delete(mock.Anything, expectedRedirect.ID, previousEtagStub).
			Once().
			Return(nil)

		// Request
		req := httptest.NewRequest(fiber.MethodDelete, baseURLStub+"/redirects/"+expectedRedirect.ID, nil)
		req.Header.Set(fiber.HeaderETag, previousEtagStub)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		deps.serviceMock.AssertExpectations(t)
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
	handler := redirect.NewHandler(loggerStub, serviceMock)
	errorHandler := fib.NewErrorHandler(loggerStub)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.Error,
	})
	handler.AddAPIRoutes(app)

	return app, handlerDeps{loggerStub: loggerStub, serviceMock: serviceMock}
}
