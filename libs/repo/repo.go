package repo

import (
	"context"

	"github.com/abtergo/abtergo/libs/model"
)

type Filter[T model.EntityInterface] interface {
	Match(ctx context.Context, entity T) (bool, error)
}

type Repository[T model.EntityInterface] interface {
	Create(ctx context.Context, entity T) (T, error)
	GetByID(ctx context.Context, id string) (T, error)
	GetByKey(ctx context.Context, key string) (T, error)
	Update(ctx context.Context, entity T, oldETag string) (T, error)
	Delete(ctx context.Context, id string, oldETag string) error
	List(ctx context.Context, filter Filter[T]) ([]T, error)
}
