package redirect

import "context"

type Filter struct {
	Website string `json:"website" form:"website" validate:"required,url" fake:"{website}"`
	Path    string `json:"path" form:"path" validate:"required" fake:"{path}"`
}

func (f Filter) Match(ctx context.Context, redirect Redirect) (bool, error) {
	_ = ctx

	if f.Website != "" && redirect.Website != f.Website {
		return false, nil
	}

	if f.Path != "" && redirect.Path != f.Path {
		return false, nil
	}

	return true, nil
}
