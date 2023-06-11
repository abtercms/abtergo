package page

import "context"

type Filter struct {
	Website string `json:"website" validate:"url"`
	Path    string `json:"path" validate:"required"`
}

func (f Filter) Match(ctx context.Context, page Page) (bool, error) {
	_ = ctx

	if f.Website != "" && f.Website != page.Website {
		return false, nil
	}

	if f.Path != "" && f.Path != page.Path {
		return false, nil
	}

	return true, nil
}
