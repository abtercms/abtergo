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

func (r Redirect) Clone() model.EntityInterface {
	c := r.c()
	c.Entity = r.Entity.Clone().(model.Entity)

	return c
}

func (r Redirect) c() Redirect {
	return Redirect{
		Website: r.Website,
		Path:    r.Path,
		Target:  r.Target,
	}
}

func (r Redirect) Validate() error {
	return redirectValidator.Struct(&r)
}
