package block

import (
	"context"
	"fmt"

	"github.com/abtergo/abtergo/libs/ablog"
	"github.com/abtergo/abtergo/libs/arr"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id string) (Block, error)
	List(ctx context.Context, filter Filter) ([]Block, error)
	Create(ctx context.Context, block Block) (Block, error)
	Update(ctx context.Context, id string, block Block, etag string) (Block, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	logger ablog.Logger
	repo   Repo
}

// NewService creates a new Service instance.
func NewService(logger ablog.Logger, repo Repo) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

// Create persists a new entity.
func (s *service) Create(ctx context.Context, entity Block) (Block, error) {
	if entity.ID != "" {
		return Block{}, arr.New(arr.InvalidUserInput, "payload must not include an id. id in payload: '%s'", entity.ID)
	}

	if err := entity.Validate(); err != nil {
		return Block{}, arr.Wrap(arr.InvalidUserInput, err)
	}

	return s.repo.Create(ctx, entity.AsNew())
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id string) (Block, error) {
	return s.repo.Retrieve(ctx, id)
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Block, error) {
	return s.repo.List(ctx, filter)
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, id string, entity Block, etag string) (Block, error) {
	if entity.ID != "" && entity.ID != id {
		return Block{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match. id in path: '%s', id in payload: '%s'", id, entity.ID)
	}

	if err := entity.Validate(); err != nil {
		return Block{}, fmt.Errorf("payload validation failed. err: %w", arr.Wrap(arr.InvalidUserInput, err))
	}

	return s.repo.Update(ctx, id, entity.AsNew(), etag)
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
