package website

import (
	"context"
	"log/slog"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
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
	page, err := s.contentRetriever.Retrieve(ctx, model.KeyFromStrings(website, path))
	if err != nil {
		return "", arr.WrapWithFallback(arr.ResourceNotFound, err, "failed to retrieve page", slog.String("website", website), slog.String("path", path))
	}

	template, err := s.templateRetriever.Retrieve(ctx, website, path)
	if err != nil {
		return "", arr.WrapWithFallback(arr.ResourceNotFound, err, "failed to retrieve template", slog.String("website", website), slog.String("path", path))
	}

	allContext := s.mergeContexts(page.GetContext(), template.GetContext())
	raw, err := s.renderer.RenderInLayout(page.Render(), template.Render(), allContext...)
	if err != nil {
		return "", arr.WrapWithFallback(arr.TemplateError, err, "failed to render content", slog.String("website", website), slog.String("path", path))
	}

	return raw, nil
}

func (s *service) mergeContexts(c1, c2 []any) []any {
	if len(c1) == 0 {
		return c2
	}

	if len(c2) == 0 {
		return c1
	}

	merged := make([]any, 0, len(c1)+len(c2))
	merged = append(merged, c1...)
	merged = append(merged, c2...)

	return merged
}
