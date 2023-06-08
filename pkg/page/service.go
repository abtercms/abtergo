package page

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, id string) (Page, error)
	List(ctx context.Context, filter Filter) ([]Page, error)
	Create(ctx context.Context, page Page) (Page, error)
	Update(ctx context.Context, id string, page Page, etag string) (Page, error)
	Delete(ctx context.Context, id string) error
	Transition(ctx context.Context, id string, trigger Trigger, oldEtag string) (Page, error)
}

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
		return Page{}, arr.Wrap(arr.InvalidUserInput, errors.New("payload must not include an id"), "id", entity.ID)
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
func (s *service) Update(ctx context.Context, id string, entity Page, etag string) (Page, error) {
	if entity.ID != "" && entity.ID != id {
		return Page{}, arr.Wrap(arr.InvalidUserInput, errors.New("path and payload ids do not match"), "id in path", id, "id in payload", entity.ID)
	}

	if err := entity.Validate(); err != nil {
		return Page{}, arr.Wrap(arr.InvalidUserInput, err, "payload validation failed")
	}

	return s.repo.Update(ctx, id, entity, etag)
}

// Delete deletes an existing entity.
func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Transition changes the status of an existing entity.
func (s *service) Transition(ctx context.Context, id string, trigger Trigger, oldEtag string) (Page, error) {
	page, err := s.repo.Retrieve(ctx, id)
	if err != nil {
		return Page{}, fmt.Errorf("failed to get page. id: %s, err: %w", id, err)
	}

	if page.Etag != oldEtag {
		return Page{}, arr.New(arr.InvalidEtag, "invalid etag found", "id", id, "request etag", oldEtag, "found etag", page.Etag)
	}

	newStatus, err := s.updater.Transition(page.Status, trigger)
	if err != nil {
		return Page{}, arr.Wrap(arr.ResourceIsOutdated, err, "failed to transition status", "id", id, "website", page.Website, "path", page.Path, "trigger", Activate)
	}

	page.Status = newStatus

	newPage, err := s.repo.Update(ctx, id, page, oldEtag)
	if err != nil {
		return Page{}, fmt.Errorf("failed to update page. id: %s, website: %s, path: %s, err: %w", newPage.ID, newPage.Website, newPage.Path, err)
	}

	return newPage, nil
}
