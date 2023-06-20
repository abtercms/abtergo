package templ_test

import (
	"fmt"

	"github.com/abtergo/abtergo/libs/templ"
	"github.com/abtergo/abtergo/libs/util"
)

type fakeRetriever struct {
	results map[string]string
}

func newFakeRetriever() *fakeRetriever {
	return &fakeRetriever{
		results: make(map[string]string),
	}
}

func (fr *fakeRetriever) Retrieve(viewTag templ.ViewTag) (string, error) {
	tag := util.ETagAny(viewTag)
	str, ok := fr.results[tag]
	if !ok {
		return "", fmt.Errorf("no result found for %s", tag)
	}

	return str, nil
}

func (fr *fakeRetriever) SetViewTag(viewTag templ.ViewTag, result string) {
	fr.results[util.ETagAny(viewTag)] = result
}
