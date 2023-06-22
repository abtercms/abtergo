package page

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
	Get(ctx context.Context, id string) (Page, error)
	List(ctx context.Context, filter Filter) ([]Page, error)
	Create(ctx context.Context, page Page) (Page, error)
	Update(ctx context.Context, id string, page Page, oldETag string) (Page, error)
	Delete(ctx context.Context, id, oldETag string) error
	Transition(ctx context.Context, id string, trigger Trigger, oldEtag string) (Page, error)
}

type Repo = repo.Repository[Page]

type service struct {
	logger  *zap.Logger
	repo    Repo
	updater Updater
}

// NewService creates a new Service instance.
func NewService(repo Repo, updater Updater, logger *zap.Logger) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		updater: updater,
	}
}

// Create persists a new entity.
func (s *service) Create(ctx context.Context, entity Page) (Page, error) {
	if entity.ID != "" {
		return Page{}, arr.New(arr.InvalidUserInput, "payload must not include an id", zap.String("id", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Page{}, arr.WrapWithType(arr.InvalidUserInput, err, "validation failed")
	}

	entity.Entity = model.NewEntity()
	entity.ETag = util.ETagAny(entity)

	entity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return Page{}, errors.Wrap(err, "failed to create page")
	}

	return entity, nil
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id string) (Page, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Page{}, errors.Wrap(err, "getting entity failed")
	}

	return entity, nil
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Page, error) {
	entities, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list entities")
	}

	return entities, nil
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, id string, entity Page, oldETag string) (Page, error) {
	if entity.ID != "" && entity.ID != id {
		return Page{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match", zap.String("id in path", id), zap.String("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Page{}, arr.WrapWithType(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.ID = id
	entity.UpdatedAt = time.Now()
	entity.ETag = util.ETagAny(entity)

	entity, err := s.repo.Update(ctx, entity, oldETag)
	if err != nil {
		return Page{}, errors.Wrap(err, "failed to update page")
	}

	return entity, nil
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id, oldETag string) error {
	err := s.repo.Delete(ctx, id, oldETag)
	if err != nil {
		return errors.Wrap(err, "failed to delete page")
	}

	return nil
}

// Transition changes the status of an existing entity.
func (s *service) Transition(ctx context.Context, id string, trigger Trigger, oldEtag string) (Page, error) {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Page{}, errors.Wrap(err, "page not found")
	}

	if page.ETag != oldEtag {
		return Page{}, arr.New(arr.ETagMismatch, "invalid e-tag found", zap.String("id", id), zap.String("request etag", oldEtag), zap.String("found etag", page.ETag))
	}

	newStatus, err := s.updater.Transition(page.Status, trigger)
	if err != nil {
		return Page{}, errors.Wrap(err, "transition failed")
	}

	page.Status = newStatus
	// TODO: update the e-tag

	newPage, err := s.repo.Update(ctx, page, oldEtag)
	if err != nil {
		return Page{}, errors.Wrap(err, "failed to update page")
	}

	return newPage, nil
}
