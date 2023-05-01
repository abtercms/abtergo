package template

import "context"

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Template) (Template, error)
	Retrieve(ctx context.Context, id string) (Template, error)
	Update(ctx context.Context, id string, entity Template, etag string) (Template, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter Filter) ([]Template, error)
}
