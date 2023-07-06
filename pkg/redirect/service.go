package redirect

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id model.ID) (Redirect, error)
	List(ctx context.Context, filter Filter) ([]Redirect, error)
	Create(ctx context.Context, redirect Redirect) (Redirect, error)
	Update(ctx context.Context, redirect Redirect, oldETag model.ETag) (Redirect, error)
	Delete(ctx context.Context, id model.ID, oldETag model.ETag) error
}

type Repo = repo.Repository[Redirect]

type service struct {
	logger *zap.Logger
	repo   Repo
}

// NewService creates a new Service instance.
func NewService(repo Repo, logger *zap.Logger) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

// Create persists a new entity.
func (s *service) Create(ctx context.Context, entity Redirect) (Redirect, error) {
	if err := entity.Validate(); err != nil {
		return Redirect{}, arr.WrapWithType(arr.InvalidUserInput, err, "validation failed")
	}

	entity.Entity = model.NewEntity()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return Redirect{}, errors.Wrap(err, "failed to create redirect")
	}

	return entity, nil
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id model.ID) (Redirect, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Redirect{}, errors.Wrap(err, "failed to retrieve redirect")
	}

	return entity, nil
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Redirect, error) {
	entities, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list redirects")
	}

	return entities, nil
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, entity Redirect, oldETag model.ETag) (Redirect, error) {
	if err := entity.Validate(); err != nil {
		return Redirect{}, arr.WrapWithType(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.ETag = ""
	entity.UpdatedAt = time.Now()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Update(ctx, entity, oldETag)
	if err != nil {
		return Redirect{}, errors.Wrap(err, "failed to update redirect")
	}

	return entity, nil
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id model.ID, oldETag model.ETag) error {
	err := s.repo.Delete(ctx, id, oldETag)
	if err != nil {
		return errors.Wrap(err, "failed to delete redirect")
	}

	return nil
}
