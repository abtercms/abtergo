package redirect

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/abtergo/abtergo/libs/arr"
)

// InMemoryRepo is a repository using a simple map to manage Redirect entities. It primarily serves testing purposes.
type InMemoryRepo struct {
	lock       sync.Mutex
	entityByID map[string]Redirect
}

// NewInMemoryRepo creates a new InMemoryRepo instance.
func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{
		entityByID: make(map[string]Redirect),
	}
}

// Create persists a new Redirect instance.
func (r *InMemoryRepo) Create(ctx context.Context, entity Redirect) (Redirect, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entity = entity.AsNew().WithEtag().WithID().WithTime()

	_, ok := r.entityByID[entity.ID]
	if ok {
		return Redirect{}, fmt.Errorf("generated the same ID twice, uuid: '%s'", entity.ID)
	}

	r.entityByID[entity.ID] = entity

	return entity, nil
}

// Retrieve retrieves an existing Redirect instance.
func (r *InMemoryRepo) Retrieve(ctx context.Context, id string) (Redirect, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	entity, ok := r.entityByID[id]
	if !ok {
		return Redirect{}, arr.New(arr.ResourceNotFound, "entity not found. id: %s", id)
	}

	return entity, nil
}

// List retrieves a list of existing Redirect instances.
func (r *InMemoryRepo) List(ctx context.Context, filter Filter) ([]Redirect, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	list := []Redirect{}
	for _, entity := range r.entityByID {
		if filter.Path != "" && entity.Path != filter.Path {
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

// Update changes an existing Redirect instance.
func (r *InMemoryRepo) Update(ctx context.Context, id string, entity Redirect, etag string) (Redirect, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	old, ok := r.entityByID[id]
	if !ok {
		return Redirect{}, arr.New(arr.ResourceNotFound, "resource not found. id: '%s'", id)
	}

	if old.Etag != etag {
		return Redirect{}, arr.New(arr.InvalidEtag, "invalid etag received. id: '%s', etag expected: '%s', etag got: '%s'", id, old.Etag, etag)
	}

	entity = entity.AsNew().WithEtag().SetID(id).SetCreatedAt(old.CreatedAt).SetUpdatedAt(old.UpdatedAt)

	if etag == entity.Etag {
		return Redirect{}, arr.New(arr.ResourceNotModified, "resource was not modified, received version appears to be the same as stored. id: '%s'", id)
	}

	entity = entity.SetUpdatedAt(time.Now())
	r.entityByID[id] = entity

	return entity, nil
}

// Delete deletes an existing Redirect instance.
func (r *InMemoryRepo) Delete(ctx context.Context, id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.entityByID, id)

	return nil
}
