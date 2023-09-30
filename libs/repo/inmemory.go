package repo

import (
	"context"
	"log/slog"
	"sync"

	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/model"
)

type index interface {
	Find(key model.Key) *model.ID
	Add(key model.Key, id model.ID) error
	Delete(key model.Key) error
	Replace(oldKey, newKey model.Key) error
}

type uniqueIndex struct {
	data   map[model.Key]model.ID
	rwLock *sync.RWMutex
}

func newUniqueIndex() index {
	return &uniqueIndex{
		data:   make(map[model.Key]model.ID),
		rwLock: &sync.RWMutex{},
	}
}

func (i *uniqueIndex) Find(key model.Key) *model.ID {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()

	if result, ok := i.data[key]; ok {
		return &result
	}

	return nil
}

func (i *uniqueIndex) Add(key model.Key, id model.ID) error {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	_, ok := i.data[key]
	if ok {
		return arr.New(arr.ApplicationError, "uniq index can not be overwritten")
	}

	i.data[key] = id

	return nil
}

func (i *uniqueIndex) Delete(key model.Key) error {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	_, ok := i.data[key]
	if !ok {
		return arr.New(arr.ApplicationError, "index not found", slog.String("key", key.String()))
	}

	delete(i.data, key)

	return nil
}

func (i *uniqueIndex) Replace(oldKey, newKey model.Key) error {
	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	val, ok := i.data[oldKey]
	if !ok {
		return arr.New(arr.ApplicationError, "target index not found", slog.String("key", string(oldKey)))
	}

	_, ok = i.data[newKey]
	if ok {
		return arr.New(arr.ApplicationError, "replacement key already exists", slog.String("key", string(newKey)))
	}

	delete(i.data, oldKey)
	i.data[newKey] = val

	return nil
}

func NewInMemoryRepo[T model.EntityInterface]() *InMemoryRepo[T] {
	return &InMemoryRepo[T]{
		entities: make(map[model.ID]T),
		indexes:  newUniqueIndex(),
		rwLock:   &sync.RWMutex{},
	}
}

type InMemoryRepo[T model.EntityInterface] struct {
	entities map[model.ID]T
	indexes  index
	rwLock   *sync.RWMutex
}

func (r *InMemoryRepo[T]) Create(ctx context.Context, entity T) (T, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	var t T

	_ = ctx
	id := entity.GetID()

	if !entity.IsComplete() {
		attrs := []slog.Attr{
			slog.String("id", id.String()),
			slog.String("etag", string(entity.GetETag())),
			slog.Time("created_at", entity.GetCreatedAt()),
			slog.Time("updated_at", entity.GetUpdatedAt()),
		}

		return t, arr.New(arr.ApplicationError, "entity not complete", attrs...)
	}

	err := r.indexes.Add(entity.GetUniqueKey(), id)
	if err != nil {
		return t, errors.Wrap(err, "index creation error")
	}

	r.entities[id] = entity

	return entity, nil
}

func (r *InMemoryRepo[T]) GetByID(ctx context.Context, id model.ID) (T, error) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()

	_ = ctx

	t, ok := r.entities[id]
	if !ok {
		return t, arr.New(arr.ResourceNotFound, "entity not found", slog.String("id", id.String()))
	}

	return t, nil
}

func (r *InMemoryRepo[T]) GetByKey(ctx context.Context, key model.Key) (T, error) {
	id := r.indexes.Find(key)
	if id == nil {
		var t T
		return t, arr.New(arr.ResourceNotFound, "index not found", slog.String("key", key.String()))
	}

	return r.GetByID(ctx, *id)
}

func (r *InMemoryRepo[T]) Update(ctx context.Context, entity T, oldETag model.ETag) (T, error) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	_ = ctx
	id := entity.GetID()

	if !entity.IsComplete() {
		var t T
		return t, arr.New(arr.ApplicationError, "entity not complete", slog.String("id", id.String()))
	}

	oldEntity, ok := r.entities[id]
	if !ok {
		return oldEntity, arr.New(arr.ResourceNotFound, "entity not found", slog.String("id", id.String()))
	}

	if oldEntity.GetETag() != oldETag {
		return oldEntity, arr.New(arr.ETagMismatch, "e-tag mismatch", slog.String("id", id.String()), slog.String("stored e-tag", oldEntity.GetETag().String()), slog.String("received e-tag", oldETag.String()))
	}

	oldKey := oldEntity.GetUniqueKey()
	newKey := entity.GetUniqueKey()
	if oldKey != newKey {
		err := r.indexes.Replace(oldKey, newKey)
		if err != nil {
			var t T
			return t, errors.Wrap(err, "index update error")
		}
	}

	r.entities[id] = entity

	return entity, nil
}

func (r *InMemoryRepo[T]) Delete(ctx context.Context, id model.ID, oldETag model.ETag) error {
	_ = ctx

	r.rwLock.Lock()
	defer r.rwLock.Unlock()

	t, ok := r.entities[id]
	if !ok {
		return arr.New(arr.ResourceNotFound, "entity not found", slog.String("id", id.String()))
	}

	if oldETag != t.GetETag() {
		return arr.New(arr.ETagMismatch, "e-tag mismatch", slog.String("id", id.String()), slog.String("stored e-tag", t.GetETag().String()), slog.String("received e-tag", oldETag.String()))
	}

	err := r.indexes.Delete(t.GetUniqueKey())
	if err != nil {
		return errors.Wrap(err, "index deletion error")
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
			return nil, arr.WrapWithType(arr.ApplicationError, err, "failed to match entity", slog.String("id", entity.GetID().String()))
		}

		if match {
			templates = append(templates, entity)
		}
	}

	return templates, nil
}
