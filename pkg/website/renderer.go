package website

import "github.com/cbroglie/mustache"

// Renderer is an interface enabling rendering of content.
type Renderer interface {
	AddContext(context ...any)
	Render(template string, context ...any) (string, error)
}

// NewRenderer creates a new Renderer instance.
func NewRenderer() Renderer {
	return &renderer{}
}

type renderer struct {
	context []any
}

// AddContext adds context for a renderer.
func (r *renderer) AddContext(context ...any) {
	r.context = append(r.context, context...)
}

// Render renders content using a template library using given template and context.
func (r *renderer) Render(template string, context ...any) (string, error) {
	allContext := append(r.context, context...)

	return mustache.Render(template, allContext...)
}
