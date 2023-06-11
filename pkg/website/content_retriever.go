package website

import (
	"github.com/pkg/errors"
)

type ContentRetriever interface {
	Retrieve(args ...interface{}) (Content, error)
}

type contentRetriever struct{}

func (c *contentRetriever) Retrieve(args ...interface{}) (Content, error) {
	_ = args

	return nil, errors.New("not implemented")
}

func NewContentRetriever() ContentRetriever {
	return &contentRetriever{}
}
