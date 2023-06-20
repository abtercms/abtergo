package templ

import (
	"strings"

	"github.com/adelowo/onecache"
	"github.com/cbroglie/mustache"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
)

// Renderer is an interface enabling rendering of content.
type Renderer interface {
	AddContext(context ...any)
	SetContext(context ...any)
	Render(template string, context ...any) (string, error)
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
	context    []any
	retrievers map[string]Retriever
	parser     Parser
	cache      onecache.Store
}

func (r *renderer) SetContext(context ...any) {
	r.context = context
}

func (r *renderer) AddContext(context ...any) {
	r.context = append(r.context, context...)
}

// Render renders content using a template library using given template and context.
func (r *renderer) Render(template string, context ...any) (string, error) {
	allContext := append(r.context, context...)

	parsedTemplate, err := mustache.Render(template, allContext...)
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
			return "", nil, arr.Wrap(err, "failed to retrieve template", zap.String("tag", viewTag.TagName), zap.String("needle example", viewTag.Needles[0]))
		}

		ccList = append(ccList, cc)

		for _, needle := range viewTag.Needles {
			template = strings.Replace(template, needle, cc.Content, -1)
		}
	}

	return template, ccList, nil
}
