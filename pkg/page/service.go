package page

import (
	"context"
	"log/slog"
	"time"

	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/repo"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id model.ID) (Page, error)
	List(ctx context.Context, filter Filter) ([]Page, error)
	Create(ctx context.Context, page Page) (Page, error)
	Update(ctx context.Context, page Page, oldETag model.ETag) (Page, error)
	Delete(ctx context.Context, id model.ID, oldETag model.ETag) error
	Transition(ctx context.Context, id model.ID, trigger Trigger, oldETag model.ETag) (Page, error)
}

type Repo = repo.Repository[Page]

type service struct {
	logger  *slog.Logger
	repo    Repo
	updater Updater
}

// NewService creates a new Service instance.
func NewService(repo Repo, updater Updater, logger *slog.Logger) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		updater: updater,
	}
}

// Create persists a new entity.
func (s *service) Create(ctx context.Context, entity Page) (Page, error) {
	if entity.ID != "" {
		return Page{}, arr.New(arr.InvalidUserInput, "payload must not include an id", slog.String("id", entity.ID.String()))
	}

	if err := entity.Validate(); err != nil {
		return Page{}, arr.WrapWithType(arr.InvalidUserInput, err, "validation failed")
	}

	entity.Entity = model.NewEntity()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Create(ctx, entity)
	if err != nil {
		return Page{}, errors.Wrap(err, "failed to create page")
	}

	return entity, nil
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id model.ID) (Page, error) {
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
func (s *service) Update(ctx context.Context, entity Page, oldETag model.ETag) (Page, error) {
	if err := entity.Validate(); err != nil {
		return Page{}, arr.WrapWithType(arr.InvalidUserInput, err, "payload validation failed")
	}

	entity.UpdatedAt = time.Now()
	entity.ETag = model.ETagFromAny(entity)

	entity, err := s.repo.Update(ctx, entity, oldETag)
	if err != nil {
		return Page{}, errors.Wrap(err, "failed to update page")
	}

	return entity, nil
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id model.ID, oldETag model.ETag) error {
	err := s.repo.Delete(ctx, id, oldETag)
	if err != nil {
		return errors.Wrap(err, "failed to delete page")
	}

	return nil
}

// Transition changes the status of an existing entity.
func (s *service) Transition(ctx context.Context, id model.ID, trigger Trigger, oldETag model.ETag) (Page, error) {
	page, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Page{}, errors.Wrap(err, "page not found")
	}

	if page.ETag != oldETag {
		return Page{}, arr.New(arr.ETagMismatch, "invalid e-tag found", slog.String("id", id.String()), slog.String("request e-tag", oldETag.String()), slog.String("found e-tag", page.ETag.String()))
	}

	newStatus, err := s.updater.Transition(page.Status, trigger)
	if err != nil {
		return Page{}, errors.Wrap(err, "transition failed")
	}

	page.Status = newStatus
	// TODO: update the e-tag

	newPage, err := s.repo.Update(ctx, page, oldETag)
	if err != nil {
		return Page{}, errors.Wrap(err, "failed to update page")
	}

	return newPage, nil
}
