package templ_test

import (
	"fmt"

	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/libs/util"
)

type retrieverDouble struct {
	results map[string]string
}

func newRetrieverDouble() *retrieverDouble {
	return &retrieverDouble{
		results: make(map[string]string),
	}
}

func (fr *retrieverDouble) Retrieve(viewTag templ.ViewTag) (templ.CacheableContent, error) {
	tag := util.ETagAny(viewTag)
	str, ok := fr.results[tag]
	if !ok {
		return nil, fmt.Errorf("no result found for %s", tag)
	}

	return templ.NewCacheableContent(str, nil), nil
}

func (fr *retrieverDouble) SetViewTag(viewTag templ.ViewTag, result string) {
	fr.results[util.ETagAny(viewTag)] = result
}
