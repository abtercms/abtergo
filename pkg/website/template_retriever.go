package website

import (
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/pkg/template"
)

type TemplateRetriever interface {
	Retrieve(args ...interface{}) (template.Template, error)
}

type templateRetriever struct{}

func (t *templateRetriever) Retrieve(args ...interface{}) (template.Template, error) {
	_ = args
	return template.Template{}, errors.New("not implemented")
}

func NewTemplateRetriever() TemplateRetriever {
	return &templateRetriever{}
}
