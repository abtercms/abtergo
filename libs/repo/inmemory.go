package repo

import (
	"context"
	"sync"

	"go.uber.org/zap"

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

	if !entity.IsComplete() {
		var (
			t    T
			args = []zap.Field{
				zap.String("id", entity.GetID()),
				zap.String("etag", entity.GetETag()),
				zap.Time("created_at", entity.GetCreatedAt()),
				zap.Time("updated_at", entity.GetUpdatedAt()),
			}
		)

		return t, arr.New(arr.ApplicationError, "entity not complete", args...)
	}

	r.entities[entity.GetID()] = entity

	return entity, nil
}

func (r *InMemoryRepo[T]) Retrieve(ctx context.Context, id string) (T, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	_ = ctx

	t, ok := r.entities[id]
	if !ok {
		return t, arr.New(arr.ResourceNotFound, "entity not found", zap.String("id", id))
	}

	return t, nil
}

func (r *InMemoryRepo[T]) Update(ctx context.Context, entity T, oldETag string) (T, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	_ = ctx

	if !entity.IsComplete() {
		var t T
		return t, arr.New(arr.ApplicationError, "entity not complete", zap.String("id", entity.GetID()))
	}

	t, ok := r.entities[entity.GetID()]
	if !ok {
		return t, arr.New(arr.ResourceNotFound, "entity not found", zap.String("id", entity.GetID()))
	}

	if t.GetETag() != oldETag {
		return t, arr.New(arr.ETagMismatch, "e-tag mismatch", zap.String("id", entity.GetID()), zap.String("stored e-tag", t.GetETag()), zap.String("received e-tag", oldETag))
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
		return arr.New(arr.ResourceNotFound, "entity not found", zap.String("id", id))
	}

	if oldETag != t.GetETag() {
		return arr.New(arr.ETagMismatch, "e-tag mismatch", zap.String("id", id), zap.String("stored e-tag", t.GetETag()), zap.String("received e-tag", oldETag))
	}

	delete(r.entities, id)

	return nil
}

func (r *InMemoryRepo[T]) List(ctx context.Context, filter Filter[T]) ([]T, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	templates := make([]T, 0)
	for _, entity := range r.entities {
		if entity.GetDeletedAt() != nil {
			continue
		}

		match, err := filter.Match(ctx, entity)
		if err != nil {
			return nil, arr.WrapWithType(arr.ApplicationError, err, "failed to match entity", zap.String("id", entity.GetID()))
		}

		if match {
			templates = append(templates, entity)
		}
	}

	return templates, nil
}
