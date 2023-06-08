package renderer

import (
	"context"
	"fmt"

	"github.com/abtergo/abtergo/pkg/block"
	"github.com/abtergo/abtergo/pkg/page"
	"github.com/abtergo/abtergo/pkg/redirect"
	"github.com/abtergo/abtergo/pkg/template"
)

// Service provides basic service functionality for Handler.
type Service interface {
	Get(ctx context.Context, website, path string) (string, error)
}

type service struct {
	renderer     Renderer
	pageRepo     page.Repo
	templateRepo template.Repo
	blockRepo    block.Repo
	redirectRepo redirect.Repo
}

// NewService creates a new Service instance.
func NewService(renderer Renderer, pageRepo page.Repo, templateRepo template.Repo, blockRepo block.Repo, redirectRepo redirect.Repo) Service {
	return &service{
		renderer:     renderer,
		pageRepo:     pageRepo,
		templateRepo: templateRepo,
		blockRepo:    blockRepo,
		redirectRepo: redirectRepo,
	}
}

// Get retrieves content for a given website+path combination.
func (s *service) Get(ctx context.Context, website, path string) (string, error) {
	p, err := s.pageRepo.RetrieveByWebsiteAndPage(ctx, website, path)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page, website: %s, path: %s, err: %w", website, path, err)
	}

	t, err := s.getTemplate(ctx, p)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve template, website: %s, path: %s, err: %w", website, path, err)
	}

	if t != nil {
		s.renderer.AddContext("content", p.Body)
		// s.renderer.Render()
	}

	return s.renderer.Render(p.Body)
}

func (s *service) getTemplate(ctx context.Context, p page.Page) (*template.Template, error) {
	// TODO: create logic to retrieve template according to custom logic

	return nil, nil
}
