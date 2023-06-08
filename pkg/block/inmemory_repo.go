package block

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	NotFound = errors.New("block not found")
)

func NewInMemoryRepo() Repo {
	return &inMemoryRepo{
		blocks: make(map[string]Block),
	}
}

type inMemoryRepo struct {
	blocks map[string]Block
	rwlock *sync.RWMutex
}

func (r *inMemoryRepo) Create(ctx context.Context, entity Block) (Block, error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	r.blocks[entity.ID] = entity

	return entity, nil
}

func (r *inMemoryRepo) Retrieve(ctx context.Context, id string) (Block, error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	return r.blocks[id], nil
}

func (r *inMemoryRepo) Update(ctx context.Context, id string, entity Block, etag string) (Block, error) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	entity.UpdatedAt = time.Now()

	r.blocks[id] = entity

	return entity, nil
}

func (r *inMemoryRepo) Delete(ctx context.Context, id string) error {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	delete(r.blocks, id)

	return nil
}

func (r *inMemoryRepo) List(ctx context.Context, filter Filter) ([]Block, error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	matches := make([]Block, 0)
	for _, block := range r.blocks {
		if filter.Match(ctx, block) {
			matches = append(matches, block)
		}
	}

	return matches, nil
}
