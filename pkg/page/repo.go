package page

import "context"

type Filter struct {
	Website string `json:"website" validate:"url"`
	Path    string `json:"path" validate:"required"`
}

func (f Filter) Match(ctx context.Context, page Page) bool {
	if f.Website != "" && f.Website != page.Website {
		return false
	}

	if f.Path != "" && f.Path != page.Path {
		return false
	}

	return true
}

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
