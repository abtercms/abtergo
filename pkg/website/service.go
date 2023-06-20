package website

import (
	"context"

	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/templ"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, website, path string) (string, error)
}

type service struct {
	contentRetriever  ContentRetriever
	templateRetriever TemplateRetriever
	renderer          templ.Renderer
}

// NewService creates a new Service instance.
func NewService(contentRetriever ContentRetriever, templateRetriever TemplateRetriever, renderer templ.Renderer) Service {
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
		return "", arr.WrapWithFallback(arr.ResourceNotFound, err, "failed to retrieve content", zap.String("website", website), zap.String("path", path))
	}

	template, err := s.templateRetriever.Retrieve(ctx, website, path)
	if err != nil {
		return "", arr.WrapWithFallback(arr.ResourceNotFound, err, "failed to retrieve template", zap.String("website", website), zap.String("path", path))
	}

	raw, err := s.renderer.Render(template.Body, content.Render())
	if err != nil {
		return "", arr.WrapWithFallback(arr.TemplateError, err, "failed to render content", zap.String("website", website), zap.String("path", path))
	}

	return raw, nil
}
