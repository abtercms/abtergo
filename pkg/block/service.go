package block

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
)

type Service interface {
	GetByID(ctx context.Context, id model.ID) (Block, error)
	List(ctx context.Context, filter Filter) ([]Block, error)
	Create(ctx context.Context, block Block) (Block, error)
	Update(ctx context.Context, block Block, oldETag model.ETag) (Block, error)
	Delete(ctx context.Context, id model.ID, oldETag model.ETag) error
}

type Repo = repo.Repository[Block]

type service struct {
	logger *slog.Logger
	repo   Repo
}

func NewService(repo Repo, logger *slog.Logger) Service {
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

func (s *service) Update(ctx context.Context, entity Block, eTag model.ETag) (Block, error) {
	if err := entity.Validate(); err != nil {
		return Block{}, arr.WrapWithType(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.UpdatedAt = time.Now()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Update(ctx, entity, eTag)
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
