package website

import (
	"context"
	"fmt"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, website, path string) (string, error)
}

type service struct {
	contentRetriever  ContentRetriever
	templateRetriever TemplateRetriever
	renderer          Renderer
}

// NewService creates a new Service instance.
func NewService(contentRetriever ContentRetriever, templateRetriever TemplateRetriever, renderer Renderer) Service {
	return &service{
		contentRetriever:  contentRetriever,
		templateRetriever: templateRetriever,
		renderer:          renderer,
	}
}

// Get retrieves content for a given website+path combination.
func (s *service) Get(ctx context.Context, website, path string) (string, error) {
	content, err := s.contentRetriever.Retrieve(ctx, website, path)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page, website: %s, path: %s, err: %w", website, path, err)
	}

	template, err := s.templateRetriever.Retrieve(ctx, website, path)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve template, website: %s, path: %s, err: %w", website, path, err)
	}

	return s.renderer.Render(template.Body, content.Render())
}
