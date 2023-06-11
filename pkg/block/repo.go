package block

import "context"

type Filter struct {
	// TODO: Use nullable types
	Website string `json:"website" validate:"url"`
	Name    string `json:"name" validate:""`
}

func (f Filter) Match(ctx context.Context, block Block) (bool, error) {
	_ = ctx

	if f.Website != "" && f.Website != block.Website {
		return false, nil
	}

	if f.Name != "" && f.Name != block.Name {
		return false, nil
	}

	return true, nil
}
