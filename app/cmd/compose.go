package main

import (
	"go.uber.org/zap"

	repo2 "github.com/abtergo/abtergo/libs/repo"
	"github.com/abtergo/abtergo/pkg/block"
	"github.com/abtergo/abtergo/pkg/page"
	"github.com/abtergo/abtergo/pkg/redirect"
	"github.com/abtergo/abtergo/pkg/template"
	"github.com/abtergo/abtergo/pkg/website"
)

func createLogger() *zap.Logger {
	return zap.Must(zap.NewProduction())
}

func createRedirectHandler(logger *zap.Logger) *redirect.Handler {
	repo := repo2.NewInMemoryRepo[redirect.Redirect]()
	service := redirect.NewService(logger, repo)
	handler := redirect.NewHandler(logger, service)

	return handler
}

func createTemplateHandler(logger *zap.Logger) *template.Handler {
	repo := repo2.NewInMemoryRepo[template.Template]()
	service := template.NewService(logger, repo)
	handler := template.NewHandler(logger, service)

	return handler
}

func createPageHandler(logger *zap.Logger) *page.Handler {
	repo := repo2.NewInMemoryRepo[page.Page]()
	updater := page.NewUpdater()
	service := page.NewService(logger, repo, updater)
	handler := page.NewHandler(logger, service)

	return handler
}

func createBlockHandler(logger *zap.Logger) *block.Handler {
	repo := repo2.NewInMemoryRepo[block.Block]()
	service := block.NewService(logger, repo)
	handler := block.NewHandler(logger, service)

	return handler
}

func createRendererHandler(logger *zap.Logger) *website.Handler {
	contentRetriever := website.NewContentRetriever()
	templateRetriever := website.NewTemplateRetriever()
	renderer := website.NewRenderer()
	service := website.NewService(contentRetriever, templateRetriever, renderer)
	handler := website.NewHandler(logger, service)

	return handler
}
