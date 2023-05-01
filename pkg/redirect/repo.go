package redirect

import "context"

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Redirect) (Redirect, error)
	Retrieve(ctx context.Context, id string) (Redirect, error)
	Update(ctx context.Context, id string, entity Redirect, etag string) (Redirect, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter Filter) ([]Redirect, error)
}
