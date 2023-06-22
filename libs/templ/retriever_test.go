package templ_test

import (
	"fmt"

	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/templ"
)

type retrieverDouble struct {
	results map[model.ETag]string
}

func newRetrieverDouble() *retrieverDouble {
	return &retrieverDouble{
		results: make(map[model.ETag]string),
	}
}

func (fr *retrieverDouble) Retrieve(viewTag templ.ViewTag) (templ.CacheableContent, error) {
	tag := model.ETagFromAny(viewTag)
	str, ok := fr.results[tag]
	if !ok {
		return nil, fmt.Errorf("no result found for %s", tag)
	}

	return templ.NewCacheableContent(str, nil), nil
}

func (fr *retrieverDouble) SetViewTag(viewTag templ.ViewTag, result string) {
	fr.results[model.ETagFromAny(viewTag)] = result
}
