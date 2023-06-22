package model

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/validation"
)

var entityValidator *validator.Validate

func init() {
	fakeit.AddDateRangeFaker()
	AddETagFaker()

	v := validation.NewValidator()

	entityValidator = v
}

type EntityInterface interface {
	GetCreatedAt() time.Time
	SetCreatedAt(t2 time.Time) EntityInterface
	GetUpdatedAt() time.Time
	SetUpdatedAt(t2 time.Time) EntityInterface
	GetDeletedAt() *time.Time
	SetDeletedAt(t2 *time.Time) EntityInterface
	GetETag() ETag
	SetETag(etag ETag) EntityInterface
	ResetETag(etag ETag) EntityInterface
	GetID() ID
	SetID(id ID) EntityInterface
	Clone() EntityInterface
	Validate() error
	IsComplete() bool
	GetUniqueKey() Key
}

type ID string

func (id ID) String() string {
	return string(id)
}

type ETag string

func (etag ETag) String() string {
	return string(etag)
}

type Key string

func (key Key) String() string {
	return string(key)
}

func id() ID {
	return ID(ulid.MustNew(ulid.Timestamp(time.Now()), nil).String())
}

type Entity struct {
	ID        ID         `json:"id,omitempty" validate:"required_with=ETag CreatedAt UpdatedAt" fake:"{uuid}"`
	ETag      ETag       `json:"etag,omitempty" validate:"required_with=ID CreatedAt UpdatedAt,etag" fake:"{etag}"`
	CreatedAt time.Time  `json:"created_at,omitempty" validate:"required_with=ID ETag UpdatedAt,not_before_date=2023-01-01" fake:"{daterange2:[2023-01-01],[2023-12-31]}"`
	UpdatedAt time.Time  `json:"updated_at,omitempty" validate:"required_with=ID ETag CreatedAt,gtecsfield=CreatedAt" fake:"{daterange2:[2024-01-01],[2024-12-31]}"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func NewEntity() Entity {
	n := time.Now()
	n2 := time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), n.Location())

	return Entity{
		ID:        id(),
		CreatedAt: n2,
		UpdatedAt: n2,
	}
}

func (e Entity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

func (e Entity) SetCreatedAt(t2 time.Time) EntityInterface {
	if e.ETag != "" {
		panic("entity is sealed.")
	}

	e.CreatedAt = t2

	return e
}

func (e Entity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

func (e Entity) SetUpdatedAt(t2 time.Time) EntityInterface {
	if e.ETag != "" {
		panic("entity is sealed.")
	}

	e.UpdatedAt = t2

	return e
}

func (e Entity) GetDeletedAt() *time.Time {
	return e.DeletedAt
}

func (e Entity) SetDeletedAt(t2 *time.Time) EntityInterface {
	if e.ETag != "" {
		panic("entity is sealed.")
	}

	e.DeletedAt = t2

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

func (e Entity) IsComplete() bool {
	return e.ETag != "" && e.ID != ""
}

func (e Entity) SetID(id ID) EntityInterface {
	if e.ETag != "" {
		panic("entity is sealed.")
	}

	e.ID = id

	return e
}

func (e Entity) GetID() ID {
	return e.ID
}

func (e Entity) ResetETag(eTag ETag) EntityInterface {
	e.ETag = eTag

	return e
}

func (e Entity) SetETag(eTag ETag) EntityInterface {
	if e.ETag != "" {
		panic("entity is sealed.")
	}

	e.ETag = eTag

	return e
}

func (e Entity) GetETag() ETag {
	return e.ETag
}

func (e Entity) GetUniqueKey() Key {
	panic(arr.New(arr.ApplicationError, "invalid unique key request", zap.String("entity", string(e.ID))))
}

func (e Entity) Validate() error {
	return entityValidator.Struct(&e)
}

func RandomEntity() Entity {
	b := Entity{}
	err := gofakeit.Struct(&b)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate random entity"))
	}

	return b
}
