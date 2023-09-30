package templ

import (
	"log/slog"
	"strings"

	"github.com/adelowo/onecache"
	"github.com/cbroglie/mustache"
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
)

// Renderer is an interface enabling rendering of content.
type Renderer interface {
	Render(template string, context ...any) (string, error)
	RenderInLayout(content, template string, context ...any) (string, error)
}

// NewRenderer creates a new Renderer instance.
func NewRenderer(parser Parser, retrievers map[string]Retriever, cache onecache.Store) Renderer {
	return &renderer{
		retrievers: retrievers,
		parser:     parser,
		cache:      cache,
	}
}

type renderer struct {
	retrievers map[string]Retriever
	parser     Parser
	cache      onecache.Store
}

func (r *renderer) RenderInLayout(content, template string, context ...any) (string, error) {
	parsedTemplate, err := mustache.RenderInLayout(content, template, context...)
	if err != nil {
		return "", errors.Wrap(err, "mustache failed to render template")
	}

	return r.Render(parsedTemplate, context...)
}

// Render renders content using a template library using given template and context.
func (r *renderer) Render(template string, context ...any) (string, error) {
	parsedTemplate, err := mustache.Render(template, context...)
	if err != nil {
		return "", errors.Wrap(err, "mustache failed to render template")
	}

	if len(r.retrievers) == 0 {
		return parsedTemplate, nil
	}

	viewTags, err := r.parser.Parse(parsedTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}

	if len(viewTags) == 0 {
		return parsedTemplate, nil
	}

	parsedTemplate, _, err = r.resolveViewTags(parsedTemplate, viewTags)
	if err != nil {
		return "", errors.Wrap(err, "failed to process template")
	}

	parsedTemplate, err = r.Render(parsedTemplate, context...)
	if err != nil {
		return "", errors.Wrap(err, "failed to render re-rendered template")
	}

	return parsedTemplate, nil
}

func (r *renderer) resolveViewTags(parsedTemplate string, viewTags []ViewTag) (string, []CacheableContent, error) {
	template := parsedTemplate
	ccList := []CacheableContent{}
	for _, viewTag := range viewTags {
		cc, err := r.retrievers[viewTag.TagName].Retrieve(viewTag)
		if err != nil {
			return "", nil, arr.Wrap(err, "failed to retrieve template", slog.String("tag", viewTag.TagName), slog.String("needle example", viewTag.Needles[0]))
		}

		ccList = append(ccList, cc)

		for _, needle := range viewTag.Needles {
			template = strings.Replace(template, needle, cc.Render(), -1)
		}
	}

	return template, ccList, nil
}
