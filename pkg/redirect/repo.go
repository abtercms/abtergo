package redirect

import "context"

type Filter struct {
	Website string `json:"website" form:"website" validate:"required,url" fake:"{website}"`
	Path    string `json:"path" form:"path" validate:"required" fake:"{path}"`
}

func (f Filter) Match(ctx context.Context, redirect Redirect) bool {
	if f.Website != "" && redirect.Website != f.Website {
		return false
	}

	if f.Path != "" && redirect.Path != f.Path {
		return false
	}

	return true
}

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Redirect) (Redirect, error)
	Retrieve(ctx context.Context, id string) (Redirect, error)
	Update(ctx context.Context, id string, entity Redirect, etag string) (Redirect, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter Filter) ([]Redirect, error)
}
