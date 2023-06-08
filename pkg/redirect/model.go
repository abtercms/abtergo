package redirect

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oklog/ulid/v2"

	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/util"
	"github.com/abtergo/abtergo/libs/validation"
)

var redirectValidator *validator.Validate

func init() {
	v := validation.NewValidator()
	redirectValidator = v
}

// Redirect represents a resource which can be used to redirect traffic from web one resource to another.
type Redirect struct {
	model.Entity

	ID      string `json:"id,omitempty" xml:"id" form:"id" validate:"required_with=Etag CreatedAt UpdatedAt" fake:"{uuid}"`
	Website string `json:"website" xml:"website" form:"website" validate:"required,url" fake:"{website}"`
	Path    string `json:"path" xml:"path" form:"path" validate:"required" fake:"{path}"`
	Target  string `json:"target,omitempty" xml:"target" form:"target" validate:"required,url" fake:"{url}"`
	Owner   string `json:"owner" xml:"owner" form:"owner" validate:"required"`
	Etag    string `json:"etag,omitempty" validate:"required_with=ID CreatedAt UpdatedAt,etag" fake:"{etag}"`
}

// Clone clones (duplicates) a Redirect resource.
func (r Redirect) Clone() Redirect {
	return Redirect{
		ID:      r.ID,
		Website: r.Website,
		Path:    r.Path,
		Target:  r.Target,
		Owner:   r.Owner,
		Etag:    r.Etag,
		Entity: model.Entity{
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
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
