package redirect

import (
	"fmt"
	"time"

	validator "github.com/go-playground/validator/v10"
	ulid "github.com/oklog/ulid/v2"

	"github.com/abtergo/abtergo/libs/util"
	"github.com/abtergo/abtergo/libs/val"
)

var redirectValidator *validator.Validate

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

	redirectValidator = v
}

// Redirect represents a resource which can be used to redirect traffic from web one resource to another.
type Redirect struct {
	ID        string    `json:"id,omitempty" xml:"id" form:"id" validate:"required_with=Etag CreatedAt UpdatedAt" fake:"{uuid}"`
	Website   string    `json:"website" xml:"website" form:"website" validate:"required,url" fake:"{website}"`
	Path      string    `json:"path" xml:"path" form:"path" validate:"required" fake:"{path}"`
	Target    string    `json:"target,omitempty" xml:"target" form:"target" validate:"required,url" fake:"{url}"`
	Owner     string    `json:"owner" xml:"owner" form:"owner" validate:"required"`
	Etag      string    `json:"etag,omitempty" validate:"required_with=ID CreatedAt UpdatedAt,etag" fake:"{etag}"`
	CreatedAt time.Time `json:"created_at,omitempty" validate:"required_with=ID Etag UpdatedAt,not_before_date=2023-01-01" fake:"{daterange2:[2023-01-01],[2023-12-31]}"`
	UpdatedAt time.Time `json:"updated_at,omitempty" validate:"required_with=ID Etag CreatedAt,gtecsfield=CreatedAt" fake:"{daterange2:[2024-01-01],[2024-12-31]}"`
}

// Clone clones (duplicates) a Redirect resource.
func (r Redirect) Clone() Redirect {
	return Redirect{
		ID:        r.ID,
		Website:   r.Website,
		Path:      r.Path,
		Target:    r.Target,
		Owner:     r.Owner,
		Etag:      r.Etag,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// AsNew returns a clone of the entity but with calculated fields reset to their default.
func (r Redirect) AsNew() Redirect {
	c := r.Clone()

	c.ID = ""
	c.Etag = ""
	c.CreatedAt = time.Time{}
	c.UpdatedAt = time.Time{}

	return c
}

// WithID returns the entity but also ensures that it has an ID.
func (r Redirect) WithID() Redirect {
	if r.ID != "" {
		return r
	}

	r.ID = ulid.Make().String()

	return r
}

// SetID returns the entity but also sets the provided ID.
func (r Redirect) SetID(id string) Redirect {
	r.ID = id

	return r
}

// WithEtag returns the entity but also ensures that it has an etag set.
func (r Redirect) WithEtag() Redirect {
	if r.Etag != "" {
		return r
	}

	r.Etag = util.EtagAny(r.AsNew())

	return r
}

// WithTime returns the entity but also ensures that it has created at and updated at set.
func (r Redirect) WithTime() Redirect {
	if r.CreatedAt.Unix() > 0 && r.UpdatedAt.Unix() > 0 {
		return r
	}

	t := time.Now()

	r.CreatedAt = t
	r.UpdatedAt = t

	return r
}

// SetCreatedAt returns the entity but also sets the provided created at.
func (r Redirect) SetCreatedAt(t time.Time) Redirect {
	r.CreatedAt = t

	return r
}

// SetUpdatedAt returns the entity but also sets the provided updated at.
func (r Redirect) SetUpdatedAt(t time.Time) Redirect {
	r.UpdatedAt = t

	return r
}

// Validate validates the entity.
func (r Redirect) Validate() error {
	return redirectValidator.Struct(&r)
}

// Filter represents a set of filters which can be used to filter the entities returned in a list request.
type Filter struct {
	Website string `json:"website" xml:"website" form:"website" validate:"required,url" fake:"{website}"`
	Path    string `json:"path" xml:"path" form:"path" validate:"required" fake:"{path}"`
}
