package website

import (
	"context"
	stdErrors "errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/decoder"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/pkg/page"
	"github.com/abtergo/abtergo/res"
)

type ContentRetriever interface {
	Retrieve(ctx context.Context, key model.Key) (templ.CacheableContent, error)
}

type contentRetriever struct {
	logger  *slog.Logger
	sources []ContentRetriever
}

func NewContentRetriever(logger *slog.Logger, sources []ContentRetriever) ContentRetriever {
	return &contentRetriever{
		logger:  logger,
		sources: sources,
	}
}

func (c *contentRetriever) Retrieve(ctx context.Context, key model.Key) (templ.CacheableContent, error) {
	for _, s := range c.sources {
		cc, err := s.Retrieve(ctx, key)
		if err == nil {
			return cc, nil
		}
	}

	return c.notFound()
}

func (c *contentRetriever) notFound() (templ.CacheableContent, error) {
	fs, err := res.Read("content/404.html")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read content/404.html")
	}

	cc := templ.NewCacheableContent(string(fs), nil)

	return cc, nil
}

type HTTPRetriever struct {
	agent       *fiber.Agent
	urlTemplate string
	decoder     *decoder.Decoder
}

func NewHTTPRetriever(agent *fiber.Agent, url string, decoder *decoder.Decoder) *HTTPRetriever {
	return &HTTPRetriever{
		agent:       agent,
		urlTemplate: url,
		decoder:     decoder,
	}
}

func (hr *HTTPRetriever) Retrieve(ctx context.Context, key model.Key) (templ.CacheableContent, error) {
	_ = ctx

	req := hr.agent.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(fmt.Sprintf(hr.urlTemplate, key))

	if err := hr.agent.Parse(); err != nil {
		panic(err)
	}

	code, body, errs := hr.agent.Bytes()
	if len(errs) > 0 {
		err := stdErrors.Join(errs...)

		return nil, arr.Wrap(err, "failed to get content (http)", slog.Int("code", code))
	}

	p := &page.Page{}
	err := hr.decoder.Decode(body, p)
	if err != nil {
		return nil, arr.Wrap(err, "failed to parse content (http)", slog.String("body", string(body)))
	}

	return p, nil
}

type MonolithRetriever struct {
	pageRepo page.Repo
}

func NewMonolithRetriever(pageRepo page.Repo) ContentRetriever {
	return &MonolithRetriever{
		pageRepo: pageRepo,
	}
}

func (mr *MonolithRetriever) Retrieve(ctx context.Context, key model.Key) (templ.CacheableContent, error) {
	result, err := mr.pageRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, arr.Wrap(err, "failed to get content (monolith)")
	}

	return result, nil
}
