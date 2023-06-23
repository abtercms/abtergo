package website

import (
	"github.com/adelowo/onecache"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/decoder"
	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/pkg/page"
)

func CreateHandler(config Config, logger *zap.Logger, cache onecache.Store, pageRepo page.Repo) *Handler {
	cr := CreateContentRetriever(config.ContentRetriever, logger, pageRepo)
	tr := NewTemplateRetriever()
	p := templ.NewParser("block")
	r := templ.NewRenderer(p, nil, cache)
	srv := NewService(cr, tr, r)

	return NewHandler(srv, logger)
}

func CreateContentRetriever(config ContentRetrieverConfig, logger *zap.Logger, pageRepo page.Repo) ContentRetriever {
	var sources []ContentRetriever
	if config.Monolith.Enabled {
		sources = append(sources, NewMonolithRetriever(pageRepo))
	}
	if config.HTTP.Enabled {
		sources = append(sources, NewHTTPRetriever(fiber.AcquireAgent(), config.HTTP.URL, decoder.NewDecoder()))
	}

	return NewContentRetriever(logger, sources)
}
