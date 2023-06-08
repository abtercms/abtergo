package redirect

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

var (
	NotFound = errors.New("redirect not found")
)

func NewInMemoryRepo() Repo {
	return &inMemoryRepo{
		redirects: make(map[string]Redirect),
	}
}

type inMemoryRepo struct {
	redirects map[string]Redirect
	rwLock    *sync.RWMutex
}

func (r *inMemoryRepo) Create(ctx context.Context, entity Redirect) (Redirect, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.redirects[entity.ID] = entity

	return entity, nil
}

func (r *inMemoryRepo) Retrieve(ctx context.Context, id string) (Redirect, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	return r.redirects[id], nil
}

func (r *inMemoryRepo) Update(ctx context.Context, id string, entity Redirect, etag string) (Redirect, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.redirects[id] = entity

	return entity, nil
}

func (r *inMemoryRepo) Delete(ctx context.Context, id string) error {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	delete(r.redirects, id)

	return nil
}

func (r *inMemoryRepo) List(ctx context.Context, filter Filter) ([]Redirect, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	matches := make([]Redirect, 0)
	for _, redirect := range r.redirects {
		if filter.Match(ctx, redirect) {
			matches = append(matches, redirect)
		}
	}

	return matches, nil
}
