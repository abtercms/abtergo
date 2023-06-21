package website

import (
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/res"
)

type TemplateRetriever interface {
	Retrieve(args ...interface{}) (templ.CacheableContent, error)
}

type templateRetriever struct{}

func (t *templateRetriever) Retrieve(args ...interface{}) (templ.CacheableContent, error) {
	_ = args

	fs, err := res.Read("templates/404.html.mustache")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read templates/404.html.mustache")
	}

	cc := templ.NewCacheableContent(string(fs), nil)

	return cc, nil
}

func NewTemplateRetriever() TemplateRetriever {
	return &templateRetriever{}
}
