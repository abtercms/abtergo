package block

import "context"

// Repo is an interface for repositories.
type Repo interface {
	Create(ctx context.Context, entity Block) (Block, error)
	Retrieve(ctx context.Context, id string) (Block, error)
	Update(ctx context.Context, id string, entity Block, etag string) (Block, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter Filter) ([]Block, error)
}
