package template

import (
	"context"

	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/repo"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id string) (Template, error)
	List(ctx context.Context, filter Filter) ([]Template, error)
	Create(ctx context.Context, template Template) (Template, error)
	Update(ctx context.Context, id string, template Template, oldETag string) (Template, error)
	Delete(ctx context.Context, id, oldETag string) error
}

type Repo = repo.Repository[Template]

type service struct {
	logger *zap.Logger
	repo   Repo
}

// NewService creates a new Service instance.
func NewService(logger *zap.Logger, repository Repo) Service {
	return &service{
		logger: logger,
		repo:   repository,
	}
}

// Create persists a new entity.
func (s *service) Create(ctx context.Context, entity Template) (Template, error) {
	if entity.ID != "" {
		return Template{}, arr.New(arr.InvalidUserInput, "payload must not include an id", zap.String("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Template{}, arr.Wrap(arr.InvalidUserInput, err, "validation failed")
	}

	return s.repo.Create(ctx, entity)
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id string) (Template, error) {
	return s.repo.Retrieve(ctx, id)
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Template, error) {
	return s.repo.List(ctx, filter)
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, id string, entity Template, oldEtag string) (Template, error) {
	if entity.ID != "" && entity.ID != id {
		return Template{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match", zap.String("id in path", id), zap.String("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Template{}, arr.Wrap(arr.InvalidUserInput, err, "payload validation failed")
	}

	return s.repo.Update(ctx, entity, oldEtag)
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id string, oldEtag string) error {
	return s.repo.Delete(ctx, id, oldEtag)
}
