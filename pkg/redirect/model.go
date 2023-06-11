package redirect

import (
	"github.com/go-playground/validator/v10"

	"github.com/abtergo/abtergo/libs/model"
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

	Website string `json:"website" xml:"website" form:"website" validate:"required,url" fake:"{website}"`
	Path    string `json:"path" xml:"path" form:"path" validate:"required" fake:"{path}"`
	Target  string `json:"target,omitempty" xml:"target" form:"target" validate:"required,url" fake:"{url}"`
}

func NewRedirect() Redirect {
	return Redirect{
		Entity: model.NewEntity(),
	}
}

// Clone clones (duplicates) a Redirect resource.
func (r Redirect) Clone() model.EntityInterface {
	c := r.AsNew().(Redirect)
	c.Entity = r.Entity.Clone().(model.Entity)

	return c
}

// AsNew returns a clone of the entity but with calculated fields reset to their default.
func (r Redirect) AsNew() model.EntityInterface {
	return Redirect{
		Entity:  model.Entity{},
		Website: r.Website,
		Path:    r.Path,
		Target:  r.Target,
	}
}

// Validate validates the entity.
func (r Redirect) Validate() error {
	return redirectValidator.Struct(&r)
}
