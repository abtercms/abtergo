package template

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/abtergo/abtergo/libs/arr"
)

// InMemoryRepo is a repository using a simple map to manage Template entities. It primarily serves testing purposes.
type InMemoryRepo struct {
	lock       sync.Mutex
	entityByID map[string]Template
}

// NewInMemoryRepo creates a new InMemoryRepo instance.
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		entityByID: make(map[string]Template),
	}
}

// Create persists a new Template instance.
func (r *InMemoryRepo) Create(ctx context.Context, entity Template) (Template, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entity = entity.AsNew().WithEtag().WithID().WithTime()

	_, ok := r.entityByID[entity.ID]
	if ok {
		return Template{}, fmt.Errorf("generated the same ID twice, uuid: '%s'", entity.ID)
	}

	r.entityByID[entity.ID] = entity

	return entity, nil
}

// Retrieve retrieves an existing Template instance.
func (r *InMemoryRepo) Retrieve(ctx context.Context, id string) (Template, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entity, ok := r.entityByID[id]
	if !ok {
		return Template{}, arr.New(arr.ResourceNotFound, "entity not found. id: %s", id)
	}

	return entity, nil
}

// List retrieves a list of existing Template instances.
func (r *InMemoryRepo) List(ctx context.Context, filter Filter) ([]Template, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	list := []Template{}
	for _, entity := range r.entityByID {
		if filter.Name != "" && entity.Name != filter.Name {
			continue
		}

		if filter.Website != "" && entity.Website != filter.Website {
			continue
		}

		list = append(list, entity)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].CreatedAt.Before(list[j].CreatedAt)
	})

	return list, nil
}

// Update changes an existing Template instance.
func (r *InMemoryRepo) Update(ctx context.Context, id string, entity Template, etag string) (Template, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	old, ok := r.entityByID[id]
	if !ok {
		return Template{}, arr.New(arr.ResourceNotFound, "resource not found. id: '%s'", id)
	}

	if old.Etag != etag {
		return Template{}, arr.New(arr.InvalidEtag, "invalid etag received. id: '%s', etag expected: '%s', etag got: '%s'", id, old.Etag, etag)
	}

	entity = entity.AsNew().WithEtag().SetID(id).SetCreatedAt(old.CreatedAt).SetUpdatedAt(old.UpdatedAt)

	if etag == entity.Etag {
		return Template{}, arr.New(arr.ResourceNotModified, "resource was not modified, received version appears to be the same as stored. id: '%s'", id)
	}

	entity = entity.SetUpdatedAt(time.Now())
	r.entityByID[id] = entity

	return entity, nil
}

// Delete deletes an existing Template instance.
func (r *InMemoryRepo) Delete(ctx context.Context, id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.entityByID, id)

	return nil
}
