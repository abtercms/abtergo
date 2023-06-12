package page

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/repo"
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
func NewService(logger *zap.Logger, repo Repo, updater Updater) Service {
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
		return Page{}, arr.Wrap(arr.InvalidUserInput, err, "validation failed")
	}

	return s.repo.Create(ctx, entity)
}

// Get retrieves an existing entity.
func (s *service) Get(ctx context.Context, id string) (Page, error) {
	return s.repo.Retrieve(ctx, id)
}

// List retrieves a collection of existing entities.
func (s *service) List(ctx context.Context, filter Filter) ([]Page, error) {
	return s.repo.List(ctx, filter)
}

// Update updates an existing entity.
func (s *service) Update(ctx context.Context, id string, entity Page, oldETag string) (Page, error) {
	if entity.ID != "" && entity.ID != id {
		return Page{}, arr.New(arr.InvalidUserInput, "path and payload ids do not match", zap.String("id in path", id), zap.String("id in payload", entity.ID))
	}

	if err := entity.Validate(); err != nil {
		return Page{}, arr.Wrap(arr.InvalidUserInput, err, "payload validation failed")
	}

	return s.repo.Update(ctx, entity, oldETag)
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id, oldETag string) error {
	return s.repo.Delete(ctx, id, oldETag)
}

// Transition changes the status of an existing entity.
func (s *service) Transition(ctx context.Context, id string, trigger Trigger, oldEtag string) (Page, error) {
	page, err := s.repo.Retrieve(ctx, id)
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
		return Page{}, errors.Wrap(err, "update failed")
	}

	return newPage, nil
}
