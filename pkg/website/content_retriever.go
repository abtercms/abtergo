package website

import (
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/res"
)

type ContentRetriever interface {
	Retrieve(args ...interface{}) (templ.CacheableContent, error)
}

type contentRetriever struct{}

func (c *contentRetriever) Retrieve(args ...interface{}) (templ.CacheableContent, error) {
	_ = args

	fs, err := res.Read("content/404.html")
	if err != nil {
		return nil, errors.Wrap(err, "failed to read content/404.html")
	}

	cc := templ.NewCacheableContent(string(fs), nil)

	return cc, nil
}

func NewContentRetriever() ContentRetriever {
	return &contentRetriever{}
}
