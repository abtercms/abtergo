package repo

import (
	"context"
	"sync"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
)

func NewInMemoryRepo[T model.EntityInterface]() *InMemoryRepo[T] {
	return &InMemoryRepo[T]{
		entities: make(map[string]T),
		rwLock:   &sync.RWMutex{},
	}
}

type InMemoryRepo[T model.EntityInterface] struct {
	entities map[string]T
	rwLock   *sync.RWMutex
}

func (r *InMemoryRepo[T]) Create(ctx context.Context, entity T) (T, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	_ = ctx

	e2 := entity.AsNew().(T)

	r.entities[e2.GetID()] = e2

	return e2, nil
}

func (r *InMemoryRepo[T]) Retrieve(ctx context.Context, id string) (T, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	_ = ctx

	t, ok := r.entities[id]
	if !ok {
		return t, arr.New(arr.ResourceNotFound, "entity not found", "id", id)
	}

	return t, nil
}

func (r *InMemoryRepo[T]) Update(ctx context.Context, entity T, oldETag string) (T, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	_ = ctx

	t, ok := r.entities[entity.GetID()]
	if !ok {
		return t, arr.New(arr.ResourceNotFound, "entity not found", "id", entity.GetID())
	}

	if t.GetETag() != oldETag {
		return t, arr.New(arr.ETagMismatch, "e-tag mismatch", "id", entity.GetID(), "stored e-tag", t.GetETag(), "received e-tag", oldETag)
	}

	r.entities[entity.GetID()] = entity

	return entity, nil
}

func (r *InMemoryRepo[T]) Delete(ctx context.Context, id string, oldETag string) error {
	_ = ctx

	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	t, ok := r.entities[id]
	if !ok {
		return arr.New(arr.ResourceNotFound, "entity not found", "id", id)
	}

	if oldETag != t.GetETag() {
		return arr.New(arr.ETagMismatch, "e-tag mismatch", "id", id, "stored e-tag", t.GetETag(), "received e-tag", oldETag)
	}

	delete(r.entities, id)

	return nil
}

func (r *InMemoryRepo[T]) List(ctx context.Context, filter Filter[T]) ([]T, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	templates := make([]T, 0)
	for _, entity := range r.entities {
		match, err := filter.Match(ctx, entity)
		if err != nil {
			return nil, arr.Wrap(arr.ApplicationError, err, "failed to match entity", "id", entity.GetID())
		}

		if match {
			templates = append(templates, entity)
		}
	}

	return templates, nil
}
