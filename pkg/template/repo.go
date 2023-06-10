package template

import "context"

type Filter struct {
	Website string `json:"website" form:"website" validate:"required,url" fake:"{website}"`
	Name    string `json:"name" form:"name" validate:"required" fake:"{sentence}"`
}

func (f Filter) Match(ctx context.Context, template Template) bool {
	if f.Website != "" && f.Website != template.Website {
		return false
	}

	if f.Name != "" && f.Name != template.Name {
		return false
	}

	return true
}

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Template) (Template, error)
	Retrieve(ctx context.Context, id string) (Template, error)
	Update(ctx context.Context, id string, entity Template, etag string) (Template, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter Filter) ([]Template, error)
}
