package template

import (
	"context"
	"sync"
)

func NewInMemoryRepo() Repo {
	return &inMemoryRepo{
		templates: make(map[string]Template),
		rwLock:    &sync.RWMutex{},
	}
}

type inMemoryRepo struct {
	templates map[string]Template
	rwLock    *sync.RWMutex
}

func (r *inMemoryRepo) Create(ctx context.Context, entity Template) (Template, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.templates[entity.ID] = entity

	return entity, nil
}

func (r *inMemoryRepo) Retrieve(ctx context.Context, id string) (Template, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	return r.templates[id], nil
}

func (r *inMemoryRepo) Update(ctx context.Context, id string, entity Template, etag string) (Template, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	r.templates[id] = entity

	return entity, nil
}

func (r *inMemoryRepo) Delete(ctx context.Context, id string) error {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	delete(r.templates, id)

	return nil
}

func (r *inMemoryRepo) List(ctx context.Context, filter Filter) ([]Template, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	templates := make([]Template, 0)
	for _, template := range r.templates {
		if filter.Match(ctx, template) {
			templates = append(templates, template)
		}
	}

	return templates, nil
}
