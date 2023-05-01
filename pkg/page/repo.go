package page

import "context"

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Page) (Page, error)
	Retrieve(ctx context.Context, id string) (Page, error)
	RetrieveByWebsiteAndPage(ctx context.Context, website, path string) (Page, error)
	List(ctx context.Context, filter Filter) ([]Page, error)
	Update(ctx context.Context, id string, entity Page, etag string) (Page, error)
	Transition(ctx context.Context, id string, oldStatus, newStatus Status, etag string) (Page, error)
	Delete(ctx context.Context, id string) error
}
