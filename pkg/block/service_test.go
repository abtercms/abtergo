package block_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abtergo/abtergo/libs/ablog"
	"github.com/abtergo/abtergo/libs/arr"
	mocks "github.com/abtergo/abtergo/mocks/libs/ablog"
	mocks2 "github.com/abtergo/abtergo/mocks/pkg/block"
	"github.com/abtergo/abtergo/pkg/block"
)

func TestService_Create(t *testing.T) {
	mockZapLogger := &mocks.ZapLogger{}
	loggerStub := ablog.Wrap(mockZapLogger)

	t.Run("id provided error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		ctxStub := context.Background()

		repoMock := new(mocks2.Repo)

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Create(ctxStub, entityStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFrom(err))
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.ID = ""
		entityStub.Website = ""
		ctxStub := context.Background()

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Create(ctxStub, entityStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFrom(err))
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock().AsNew()
		ctxStub := context.Background()

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Create(ctxStub, entityStub).
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.Create(ctxStub, entityStub)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}

func TestService_Delete(t *testing.T) {
	mockZapLogger := &mocks.ZapLogger{}
	loggerStub := ablog.Wrap(mockZapLogger)

	t.Run("success", func(t *testing.T) {
		ctxStub := context.Background()
		entityStub := block.RandomBlock()

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Delete(ctxStub, entityStub.ID).
			Return(nil)

		s := block.NewService(loggerStub, repoMock)

		err := s.Delete(ctxStub, entityStub.ID)
		require.NoError(t, err)
	})
}

func TestService_Get(t *testing.T) {
	mockZapLogger := &mocks.ZapLogger{}
	loggerStub := ablog.Wrap(mockZapLogger)

	t.Run("success", func(t *testing.T) {
		ctxStub := context.Background()
		entityStub := block.RandomBlock()

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Retrieve(ctxStub, entityStub.ID).
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.Get(ctxStub, entityStub.ID)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}

func TestService_List(t *testing.T) {
	mockZapLogger := &mocks.ZapLogger{}
	loggerStub := ablog.Wrap(mockZapLogger)

	t.Run("success", func(t *testing.T) {
		ctxStub := context.Background()
		filterStub := block.Filter{}
		stubCollection := block.RandomBlockList(1, 3)

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			List(ctxStub, filterStub).
			Return(stubCollection, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.List(ctxStub, filterStub)
		require.NoError(t, err)

		assert.Equal(t, stubCollection, got)
	})
}

func TestService_Update(t *testing.T) {
	mockZapLogger := &mocks.ZapLogger{}
	loggerStub := ablog.Wrap(mockZapLogger)

	const (
		idStub   = "foo"
		etagStub = "bar"
	)

	t.Run("id mismatch error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		ctxStub := context.Background()

		repoMock := new(mocks2.Repo)

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFrom(err))
	})

	t.Run("validation error", func(t *testing.T) {
		entityStub := block.RandomBlock()
		entityStub.Website = ""
		entityStub.ID = ""
		ctxStub := context.Background()

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Update(ctxStub, idStub, entityStub, etagStub).
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		_, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.Error(t, err)

		assert.Equal(t, http.StatusBadRequest, arr.HTTPStatusFrom(err))
	})

	t.Run("success", func(t *testing.T) {
		entityStub := block.RandomBlock().AsNew()
		ctxStub := context.Background()

		repoMock := new(mocks2.Repo)
		repoMock.EXPECT().
			Update(ctxStub, idStub, entityStub, etagStub).
			Return(entityStub, nil)

		s := block.NewService(loggerStub, repoMock)

		got, err := s.Update(ctxStub, idStub, entityStub, etagStub)
		require.NoError(t, err)

		assert.Equal(t, entityStub, got)
	})
}
