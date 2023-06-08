package page

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
)

var (
	NotFound = errors.New("page not found")
)

func NewInMemoryRepo() Repo {
	return &inMemoryRepo{
		pages: make(map[string]Page),
	}
}

type inMemoryRepo struct {
	pages  map[string]Page
	rwLock *sync.RWMutex
}

func (r *inMemoryRepo) Create(ctx context.Context, entity Page) (Page, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.pages[entity.ID] = entity

	return entity, nil
}

func (r *inMemoryRepo) Retrieve(ctx context.Context, id string) (Page, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	return r.pages[id], nil
}

func (r *inMemoryRepo) RetrieveByWebsiteAndPage(ctx context.Context, website, path string) (Page, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	for _, page := range r.pages {
		if page.Website == website && page.Path == path {
			return page, nil
		}
	}

	return Page{}, arr.Wrap(arr.ResourceNotFound, NotFound, "website", website, "path", path)
}

func (r *inMemoryRepo) List(ctx context.Context, filter Filter) ([]Page, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	matches := make([]Page, 0)
	for _, page := range r.pages {
		if filter.Match(ctx, page) {
			matches = append(matches, page)
		}
	}

	return matches, nil
}

func (r *inMemoryRepo) Update(ctx context.Context, id string, entity Page, etag string) (Page, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.pages[id] = entity

	return entity, nil
}

func (r *inMemoryRepo) Transition(ctx context.Context, id string, oldStatus, newStatus Status, etag string) (Page, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	entity, err := r.Retrieve(ctx, id)
	if err != nil {
		return Page{}, err
	}

	r.pages[id] = entity

	return entity, nil
}

func (r *inMemoryRepo) Delete(ctx context.Context, id string) error {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	delete(r.pages, id)

	return nil
}
