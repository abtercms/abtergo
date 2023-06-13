package block

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
	"github.com/abtergo/abtergo/libs/util"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id string) (Block, error)
	List(ctx context.Context, filter Filter) ([]Block, error)
	Create(ctx context.Context, block Block) (Block, error)
	Update(ctx context.Context, id string, block Block, oldETag string) (Block, error)
	Delete(ctx context.Context, id, oldETag string) error
}

type Repo = repo.Repository[Block]

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
func (s *service) Create(ctx context.Context, entity Block) (Block, error) {
	if err := entity.Validate(); err != nil {
		return Block{}, arr.WrapWithType(arr.InvalidUserInput, err, "validation failed")
	}

	entity.Entity = model.NewEntity()
	entity.ETag = util.ETagAny(entity)

	entity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return Block{}, errors.Wrap(err, "failed to create entity")
	}

	return entity, nil
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id string) (Block, error) {
	entity, err := s.repo.Retrieve(ctx, id)
	if err != nil {
		return Block{}, errors.Wrap(err, "failed to retrieve entity")
	}

	return entity, nil
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Block, error) {
	entities, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list entities")
	}

	return entities, nil
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, id string, entity Block, etag string) (Block, error) {
	if entity.ID != "" && entity.ID != id {
		return Block{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match", zap.String("id in path", id), zap.String("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Block{}, arr.WrapWithType(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.ID = id
	entity.UpdatedAt = time.Now()
	entity.ETag = util.ETagAny(entity)

	entity, err := s.repo.Update(ctx, entity, etag)
	if err != nil {
		return Block{}, errors.Wrap(err, "failed to update entity")
	}

	return entity, nil
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id, oldETag string) error {
	err := s.repo.Delete(ctx, id, oldETag)
	if err != nil {
		return errors.Wrap(err, "failed to update entity")
	}

	return nil
}
