package redirect

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
	"github.com/abtergo/abtergo/libs/util"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id string) (Redirect, error)
	List(ctx context.Context, filter Filter) ([]Redirect, error)
	Create(ctx context.Context, redirect Redirect) (Redirect, error)
	Update(ctx context.Context, id string, redirect Redirect, oldETag string) (Redirect, error)
	Delete(ctx context.Context, id, oldETag string) error
}

type Repo = repo.Repository[Redirect]

type service struct {
	logger *zap.Logger
	repo   Repo
}

// NewService creates a new Service instance.
func NewService(logger *zap.Logger, repo Repo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

// Create persists a new entity.
func (s *service) Create(ctx context.Context, entity Redirect) (Redirect, error) {
	if err := entity.Validate(); err != nil {
		return Redirect{}, arr.Wrap(arr.InvalidUserInput, err, "validation failed")
	}

	entity.Entity = model.NewEntity()
	entity.ETag = util.ETagAny(entity)

	return s.repo.Create(ctx, entity)
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id string) (Redirect, error) {
	return s.repo.Retrieve(ctx, id)
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Redirect, error) {
	return s.repo.List(ctx, filter)
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, id string, entity Redirect, oldETag string) (Redirect, error) {
	if entity.ID != "" && entity.ID != id {
		return Redirect{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match", zap.String("id in path", id), zap.String("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Redirect{}, arr.Wrap(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.ID = id
	entity.ETag = ""
	entity.UpdatedAt = time.Now()
	entity.ETag = util.ETagAny(entity)

	return s.repo.Update(ctx, entity, oldETag)
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id, oldETag string) error {
	return s.repo.Delete(ctx, id, oldETag)
}
