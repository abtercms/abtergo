package template_test

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

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
	mocks "github.com/abtergo/abtergo/mocks/pkg/template"
	"github.com/abtergo/abtergo/pkg/template"
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
		resp, err := app.Test(req, 10000000)
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
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/templates", reqBody)
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
		expectedTemplate := template.RandomTemplate(false)

		// Stubs
		payloadStub := expectedTemplate.Clone()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Create(mock.Anything, payloadStub).
			Once().
			Return(template.Template{}, arr.WrapWithType(arr.ResourceIsOutdated, assert.AnError, "foo"))

		// Request
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/templates", reqBody)
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
		expectedTemplate := template.RandomTemplate(false)

		// Stubs
		payloadStub := expectedTemplate.Clone()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Create(mock.Anything, payloadStub).
			Once().
			Return(expectedTemplate, nil)

		// Request
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/templates", reqBody)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual template.Template
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedTemplate.Clone(), actual.Clone())

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
			List(mock.Anything, template.Filter{}).
			Once().
			Return(nil, arr.WrapWithType(arr.ETagMismatch, assert.AnError, "foo"))

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/templates", nil)
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
		expectedTemplates := template.RandomTemplateList(5, 5)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			List(mock.Anything, template.Filter{}).
			Once().
			Return(expectedTemplates, nil)

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/templates", nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual []template.Template
		util.ParseResponseHelper(t, resp, &actual)
		assert.Len(t, actual, 5)
		assert.Equal(t, expectedTemplates[0], actual[0])

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Get(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("error retrieving entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusConflict
		expectedTemplate := template.RandomTemplate(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Get(mock.Anything, expectedTemplate.ID).
			Once().
			Return(template.Template{}, arr.WrapWithType(arr.ResourceIsOutdated, assert.AnError, "foo"))

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		req := httptest.NewRequest(fiber.MethodGet, target, nil)
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
		expectedTemplate := template.RandomTemplate(false)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Get(mock.Anything, expectedTemplate.ID).
			Once().
			Return(expectedTemplate, nil)

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		req := httptest.NewRequest(fiber.MethodGet, target, nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual template.Template
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedTemplate, actual)

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Put(t *testing.T) {
	const (
		baseURLStub                 = "https://example.com"
		previousETagStub model.ETag = "foo"
	)

	t.Run("error parsing payload", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusBadRequest
		expectedTemplate := template.RandomTemplate(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		reqBody := util.DataToReaderHelper(t, `"foo"`)
		req := httptest.NewRequest(fiber.MethodPut, target, reqBody)
		req.Header.Set(fiber.HeaderETag, previousETagStub.String())
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
		expectedTemplate := template.RandomTemplate(false)

		// Stubs
		payloadStub := expectedTemplate.Clone()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		errStub := arr.WrapWithType(arr.UpstreamServiceUnavailable, assert.AnError, "foo")
		deps.serviceMock.EXPECT().
			Update(mock.Anything, expectedTemplate.ID, payloadStub, previousETagStub).
			Once().
			Return(template.Template{}, errStub)

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPut, target, reqBody)
		req.Header.Set(fiber.HeaderETag, previousETagStub.String())
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
		expectedTemplate := template.RandomTemplate(false)

		// Stubs
		payloadStub := expectedTemplate.Clone()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Update(mock.Anything, expectedTemplate.ID, payloadStub, previousETagStub).
			Once().
			Return(expectedTemplate, nil)

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPut, target, reqBody)
		req.Header.Set(fiber.HeaderETag, previousETagStub.String())
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual template.Template
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedTemplate, actual)

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Delete(t *testing.T) {
	const (
		baseURLStub                 = "https://example.com"
		previousEtagStub model.ETag = "foo"
	)

	t.Run("error deleting entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusTooManyRequests
		expectedTemplate := template.RandomTemplate(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Delete(mock.Anything, expectedTemplate.ID, previousEtagStub).
			Once().
			Return(arr.WrapWithType(arr.UpstreamServiceBusy, assert.AnError, "foo"))

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		req := httptest.NewRequest(fiber.MethodDelete, target, nil)
		req.Header.Set(fiber.HeaderETag, previousEtagStub.String())
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
		expectedTemplate := template.RandomTemplate(false)

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Delete(mock.Anything, expectedTemplate.ID, previousEtagStub).
			Once().
			Return(nil)

		// Request
		target := fmt.Sprintf("%s/templates/%s", baseURLStub, expectedTemplate.ID)
		req := httptest.NewRequest(fiber.MethodDelete, target, nil)
		req.Header.Set(fiber.HeaderETag, previousEtagStub.String())
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

	loggerMock := zaptest.NewLogger(t)
	serviceMock := &mocks.Service{}
	handler := template.NewHandler(serviceMock, loggerMock)
	errorHandler := fib.NewErrorHandler(loggerMock)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.Handle,
		ReadTimeout:  time.Hour,
		WriteTimeout: time.Hour,
		IdleTimeout:  time.Hour,
	})
	handler.AddAPIRoutes(app)

	return app, handlerDeps{loggerStub: loggerMock, serviceMock: serviceMock}
}
