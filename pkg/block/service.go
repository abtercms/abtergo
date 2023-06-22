package block

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
)

type Service interface {
	GetByID(ctx context.Context, id model.ID) (Block, error)
	List(ctx context.Context, filter Filter) ([]Block, error)
	Create(ctx context.Context, block Block) (Block, error)
	Update(ctx context.Context, id model.ID, block Block, oldETag model.ETag) (Block, error)
	Delete(ctx context.Context, id model.ID, oldETag model.ETag) error
}

type Repo = repo.Repository[Block]

type service struct {
	logger *zap.Logger
	repo   Repo
}

func NewService(repo Repo, logger *zap.Logger) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Create(ctx context.Context, entity Block) (Block, error) {
	if err := entity.Validate(); err != nil {
		return Block{}, arr.WrapWithType(arr.InvalidUserInput, err, "validation failed")
	}

	entity.Entity = model.NewEntity()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return Block{}, errors.Wrap(err, "failed to create entity")
	}

	return entity, nil
}

func (s *service) GetByID(ctx context.Context, id model.ID) (Block, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Block{}, errors.Wrap(err, "failed to retrieve entity")
	}

	return entity, nil
}

func (s *service) List(ctx context.Context, filter Filter) ([]Block, error) {
	entities, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list entities")
	}

	return entities, nil
}

func (s *service) Update(ctx context.Context, id model.ID, entity Block, etag model.ETag) (Block, error) {
	if entity.ID != "" && entity.ID != id {
		return Block{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match", zap.Stringer("id in path", id), zap.Stringer("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Block{}, arr.WrapWithType(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.ID = id
	entity.UpdatedAt = time.Now()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Update(ctx, entity, etag)
	if err != nil {
		return Block{}, errors.Wrap(err, "failed to update entity")
	}

	return entity, nil
}

func (s *service) Delete(ctx context.Context, id model.ID, oldETag model.ETag) error {
	err := s.repo.Delete(ctx, id, oldETag)
	if err != nil {
		return errors.Wrap(err, "failed to update entity")
	}

	return nil
}
