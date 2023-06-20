package http

import (
	"github.com/adelowo/onecache"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/repo"
	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/pkg/block"
	"github.com/abtergo/abtergo/pkg/page"
	"github.com/abtergo/abtergo/pkg/redirect"
	"github.com/abtergo/abtergo/pkg/template"
	"github.com/abtergo/abtergo/pkg/website"
)

func createRedirectHandler(logger *zap.Logger) *redirect.Handler {
	repository := repo.NewInMemoryRepo[redirect.Redirect]()
	service := redirect.NewService(repository, logger)
	handler := redirect.NewHandler(service, logger)

	return handler
}

func createTemplateHandler(logger *zap.Logger) *template.Handler {
	repository := repo.NewInMemoryRepo[template.Template]()
	service := template.NewService(repository, logger)
	handler := template.NewHandler(service, logger)

	return handler
}

func createPageHandler(logger *zap.Logger) *page.Handler {
	repository := repo.NewInMemoryRepo[page.Page]()
	updater := page.NewUpdater()
	service := page.NewService(repository, updater, logger)
	handler := page.NewHandler(service, logger)

	return handler
}

func createBlockHandler(logger *zap.Logger) *block.Handler {
	repository := repo.NewInMemoryRepo[block.Block]()
	service := block.NewService(repository, logger)
	handler := block.NewHandler(service, logger)

	return handler
}

func createRendererHandler(logger *zap.Logger, cache onecache.Store) *website.Handler {
	contentRetriever := website.NewContentRetriever()
	templateRetriever := website.NewTemplateRetriever()
	parser := templ.NewParser("block")
	renderer := templ.NewRenderer(parser, nil, cache)
	service := website.NewService(contentRetriever, templateRetriever, renderer)
	handler := website.NewHandler(service, logger)

	return handler
}
