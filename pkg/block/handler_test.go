package block_test

import (
	"net/http/httptest"
	"testing"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/fib"
	"github.com/abtergo/abtergo/libs/problem"
	"github.com/abtergo/abtergo/libs/util"
	mocks "github.com/abtergo/abtergo/mocks/libs/ablog"
	mocks2 "github.com/abtergo/abtergo/mocks/pkg/block"
	"github.com/abtergo/abtergo/pkg/block"
)

func TestHandler_AddApiRoutes(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("Undefined route results in 404", func(t *testing.T) {
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
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/blocks", reqBody)
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
		expectedBlock := block.RandomBlock()

		// Stubs
		payloadStub := expectedBlock.AsNew()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Create(mock.Anything, payloadStub).
			Once().
			Return(block.Block{}, arr.Wrap(arr.ResourceIsOutdated, assert.AnError))

		// Request
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/blocks", reqBody)
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
		expectedBlock := block.RandomBlock()

		// Stubs
		payloadStub := expectedBlock.AsNew()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Create(mock.Anything, payloadStub).
			Once().
			Return(expectedBlock, nil)

		// Request
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPost, baseURLStub+"/blocks", reqBody)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual block.Block
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedBlock.AsNew(), actual.AsNew())

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
			List(mock.Anything, block.Filter{}).
			Once().
			Return(nil, arr.Wrap(arr.InvalidEtag, assert.AnError))

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/blocks", nil)
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
		expectedBlocks := block.RandomBlockList(5, 5)

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			List(mock.Anything, block.Filter{}).
			Once().
			Return(expectedBlocks, nil)

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/blocks", nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual []block.Block
		util.ParseResponseHelper(t, resp, &actual)
		assert.Len(t, actual, 5)
		assert.Equal(t, expectedBlocks[0].AsNew(), actual[0].AsNew())

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Get(t *testing.T) {
	const baseURLStub = "https://example.com"

	t.Run("error retrieving entity", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusConflict
		expectedBlock := block.RandomBlock()

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Get(mock.Anything, expectedBlock.ID).
			Once().
			Return(block.Block{}, arr.Wrap(arr.ResourceIsOutdated, assert.AnError))

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/blocks/"+expectedBlock.ID, nil)
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
		expectedBlock := block.RandomBlock()

		// Stubs

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.
			EXPECT().
			Get(mock.Anything, expectedBlock.ID).
			Once().
			Return(expectedBlock, nil)

		// Request
		req := httptest.NewRequest(fiber.MethodGet, baseURLStub+"/blocks/"+expectedBlock.ID, nil)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual block.Block
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedBlock.AsNew(), actual.AsNew())

		deps.serviceMock.AssertExpectations(t)
	})
}

func TestHandler_Put(t *testing.T) {
	const (
		baseURLStub = "https://example.com"

		// Stubs
		previousEtagStub = "foo"
	)

	t.Run("error parsing payload", func(t *testing.T) {
		// Expectations
		expectedStatusCode := fiber.StatusBadRequest
		expectedBlock := block.RandomBlock()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks

		// Request
		reqBody := util.DataToReaderHelper(t, `"foo"`)
		req := httptest.NewRequest(fiber.MethodPut, baseURLStub+"/blocks/"+expectedBlock.ID, reqBody)
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
		expectedBlock := block.RandomBlock()

		// Stubs
		payloadStub := expectedBlock.AsNew()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Update(mock.Anything, expectedBlock.ID, payloadStub, previousEtagStub).
			Once().
			Return(block.Block{}, arr.Wrap(arr.UpstreamServiceUnavailable, assert.AnError))

		// Request
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPut, baseURLStub+"/blocks/"+expectedBlock.ID, reqBody)
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
		expectedBlock := block.RandomBlock()

		// Stubs
		payloadStub := expectedBlock.AsNew()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Update(mock.Anything, expectedBlock.ID, payloadStub, previousEtagStub).
			Once().
			Return(expectedBlock, nil)

		// Request
		reqBody := util.DataToReaderHelper(t, payloadStub)
		req := httptest.NewRequest(fiber.MethodPut, baseURLStub+"/blocks/"+expectedBlock.ID, reqBody)
		req.Header.Set(fiber.HeaderETag, previousEtagStub)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		// Execute Test
		resp, err := app.Test(req)
		defer resp.Body.Close()

		// Asserts
		require.NoError(t, err)
		require.Equal(t, expectedStatusCode, resp.StatusCode)

		var actual block.Block
		util.ParseResponseHelper(t, resp, &actual)
		assert.Equal(t, expectedBlock.AsNew(), actual.AsNew())

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
		expectedBlock := block.RandomBlock()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().
			Delete(mock.Anything, expectedBlock.ID).
			Once().
			Return(arr.Wrap(arr.UpstreamServiceBusy, assert.AnError))

		// Request
		req := httptest.NewRequest(fiber.MethodDelete, baseURLStub+"/blocks/"+expectedBlock.ID, nil)
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
		expectedBlock := block.RandomBlock()

		// Prepare Test
		app, deps := setupHandlerMocks(t)

		// Mocks
		deps.serviceMock.EXPECT().Delete(mock.Anything, expectedBlock.ID).Once().Return(nil)

		// Request
		req := httptest.NewRequest(fiber.MethodDelete, baseURLStub+"/blocks/"+expectedBlock.ID, nil)
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
	loggerMock  *mocks.Logger
	serviceMock *mocks2.Service
}

func (hd handlerDeps) AssertExpectations(t *testing.T) {
	hd.loggerMock.AssertExpectations(t)
	hd.serviceMock.AssertExpectations(t)
}

func setupHandlerMocks(t *testing.T) (*fiber.App, handlerDeps) {
	loggerMock := &mocks.Logger{}
	serviceMock := &mocks2.Service{}
	handler := block.NewHandler(loggerMock, serviceMock)

	app := fiber.New(fiber.Config{
		ErrorHandler: fib.ErrorHandler,
	})
	handler.AddAPIRoutes(app)

	return app, handlerDeps{loggerMock: loggerMock, serviceMock: serviceMock}
}
