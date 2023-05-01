package block

import (
	"fmt"
	"time"

	validator "github.com/go-playground/validator/v10"
	ulid "github.com/oklog/ulid/v2"

	"github.com/abtergo/abtergo/libs/util"
	"github.com/abtergo/abtergo/libs/val"
	"github.com/abtergo/abtergo/pkg/html"
)

var blockValidator *validator.Validate

func init() {
	v := val.NewValidator()

	// Register custom validators here
	err := v.RegisterValidation("not_before_date", val.ValidateNotBeforeDate)
	if err != nil {
		panic(fmt.Errorf("failed to register 'not before date' validator. err: %w", err))
	}
	err = v.RegisterValidation("etag", val.ValidateEtag)
	if err != nil {
		panic(fmt.Errorf("failed to register 'etag' validator. err: %w", err))
	}

	blockValidator = v
}

// Block represents a resource ready which can be used as wrapper for page content.
type Block struct {
	ID        string      `json:"id,omitempty" validate:"required_with=Etag CreatedAt UpdatedAt" fake:"{uuid}"`
	Website   string      `json:"website" validate:"required,url" fake:"{website}"`
	Name      string      `json:"name" validate:"required" fake:"{sentence:1}"`
	Body      string      `json:"body" validate:"required" fake:"{paragraph:5}"`
	Assets    html.Assets `json:"assets" validate:"dive"`
	Owner     string      `json:"owner" validate:"required"`
	Etag      string      `json:"etag,omitempty" validate:"required_with=ID CreatedAt UpdatedAt,etag" fake:"{etag}"`
	CreatedAt time.Time   `json:"created_at,omitempty" validate:"required_with=ID Etag UpdatedAt,not_before_date=2023-01-01" fake:"{daterange2:[2023-01-01],[2023-12-31]}"`
	UpdatedAt time.Time   `json:"updated_at,omitempty" validate:"required_with=ID Etag CreatedAt,gtecsfield=CreatedAt" fake:"{daterange2:[2024-01-01],[2024-12-31]}"`
}

// Clone clones (duplicates) a Block resource.
func (t Block) Clone() Block {
	return Block{
		ID:        t.ID,
		Website:   t.Website,
		Name:      t.Name,
		Body:      t.Body,
		Assets:    t.Assets.Clone(),
		Owner:     t.Owner,
		Etag:      t.Etag,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// AsNew returns a clone of the entity but with calculated fields reset to their default.
func (t Block) AsNew() Block {
	c := t.Clone()

	c.ID = ""
	c.Etag = ""
	c.CreatedAt = time.Time{}
	c.UpdatedAt = time.Time{}

	return c
}

// WithID returns the entity but also ensures that it has an ID.
func (t Block) WithID() Block {
	if t.ID != "" {
		return t
	}

	t.ID = ulid.Make().String()

	return t
}

// SetID returns the entity but also sets the provided ID.
func (t Block) SetID(id string) Block {
	t.ID = id

	return t
}

// WithEtag returns the entity but also ensures that it has an etag set.
func (t Block) WithEtag() Block {
	if t.Etag != "" {
		return t
	}

	t.Etag = util.EtagAny(t.Clone().AsNew())

	return t
}

// WithTime returns the entity but also ensures that it has created at and updated at set.
func (t Block) WithTime() Block {
	if t.CreatedAt.Unix() > 0 && t.UpdatedAt.Unix() > 0 {
		return t
	}

	t2 := time.Now()

	t.CreatedAt = t2
	t.UpdatedAt = t2

	return t
}

// SetCreatedAt returns the entity but also sets the provided created at.
func (t Block) SetCreatedAt(t2 time.Time) Block {
	t.CreatedAt = t2

	return t
}

// SetUpdatedAt returns the entity but also sets the provided updated at.
func (t Block) SetUpdatedAt(t2 time.Time) Block {
	t.UpdatedAt = t2

	return t
}

// Validate validates the entity.
func (t Block) Validate() error {
	return blockValidator.Struct(&t)
}

// Filter represents a set of filters which can be used to filter the entities returned in a list request.
type Filter struct {
	Website string `json:"website" validate:"required,url" fake:"{website}"`
	Name    string `json:"name" validate:"required" fake:"{sentence:1}"`
}
