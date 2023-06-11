package template

import "context"

type Filter struct {
	Website string `json:"website" form:"website" validate:"required,url" fake:"{website}"`
	Name    string `json:"name" form:"name" validate:"required" fake:"{sentence}"`
}

func (f Filter) Match(ctx context.Context, template Template) (bool, error) {
	_ = ctx

	if f.Website != "" && f.Website != template.Website {
		return false, nil
	}

	if f.Name != "" && f.Name != template.Name {
		return false, nil
	}

	return true, nil
}
