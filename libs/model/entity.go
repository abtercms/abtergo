package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"

	"github.com/abtergo/abtergo/libs/validation"
)

var entityValidator *validator.Validate

func init() {
	v := validation.NewValidator()

	entityValidator = v
}

type EntityInterface interface {
	SetCreatedAt(t2 time.Time) EntityInterface
	SetUpdatedAt(t2 time.Time) EntityInterface
	SetETag(id string) EntityInterface
	GetETag() string
	SetID(id string) EntityInterface
	GetID() string
	Clone() EntityInterface
	AsNew() EntityInterface
}

type Entity struct {
	ID        string     `json:"id,omitempty" validate:"required_with=ETag CreatedAt UpdatedAt" fake:"{uuid}"`
	ETag      string     `json:"etag,omitempty" validate:"required_with=ID CreatedAt UpdatedAt,etag" fake:"{etag}"`
	CreatedAt time.Time  `json:"created_at,omitempty" validate:"required_with=ID ETag UpdatedAt,not_before_date=2023-01-01" fake:"{daterange2:[2023-01-01],[2023-12-31]}"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" validate:"required_with=ID ETag CreatedAt,gtecsfield=CreatedAt" fake:"{daterange2:[2024-01-01],[2024-12-31]}"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewEntity() Entity {
	// time.Now() returns an extra monotonic clock which we usually don't need.
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())

	return Entity{
		ID:        id(),
		CreatedAt: t,
		UpdatedAt: t,
	}
}

func id() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), nil).String()
}

func (e Entity) AsNew() EntityInterface {
	n := time.Now()
	n2 := time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), n.Location())

	return Entity{
		ID:        id(),
		CreatedAt: n2,
		UpdatedAt: n2,
	}
}

func (e Entity) Validate() error {
	return entityValidator.Struct(&e)
}

func (e Entity) SetCreatedAt(t2 time.Time) EntityInterface {
	e.CreatedAt = t2

	return e
}

func (e Entity) SetUpdatedAt(t2 time.Time) EntityInterface {
	e.UpdatedAt = t2

	return e
}

func (e Entity) SetDeletedAt(t2 time.Time) EntityInterface {
	e.DeletedAt = &t2

	return e
}

func (e Entity) Clone() EntityInterface {
	return Entity{
		ID:        e.ID,
		ETag:      e.ETag,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

func (e Entity) SetID(id string) EntityInterface {
	e.ID = id

	return e
}

func (e Entity) GetID() string {
	return e.ID
}

func (e Entity) SetETag(etag string) EntityInterface {
	e.ETag = etag

	return e
}

func (e Entity) GetETag() string {
	return e.ETag
}
