package block

import "context"

type Filter struct {
	Website string `json:"website" validate:"required,url" fake:"{website}"`
	Name    string `json:"name" validate:"required" fake:"{sentence:1}"`
}

func NewFilter() Filter {
	return Filter{}
}

func (f Filter) Match(ctx context.Context, block Block) bool {
	if f.Website != "" && f.Website != block.Website {
		return false
	}

	if f.Name != "" && f.Name != block.Name {
		return false
	}

	return true
}

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Block) (Block, error)
	Retrieve(ctx context.Context, id string) (Block, error)
	Update(ctx context.Context, id string, entity Block, etag string) (Block, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter Filter) ([]Block, error)
}
