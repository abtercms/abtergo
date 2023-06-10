package template

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/validation"
)

var templateValidator *validator.Validate

func init() {
	v := validation.NewValidator()

	templateValidator = v
}

type Template struct {
	model.Entity

	Website    string      `json:"website" validate:"required,url" fake:"{website}"`
	Name       string      `json:"name" validate:"required" fake:"{sentence:1}"`
	Body       string      `json:"body" validate:"required" fake:"{paragraph:10}"`
	Assets     html.Assets `json:"assets" validate:"dive"`
	HTTPHeader http.Header `json:"http_header" validate:"dive,required"`
	Version    int64       `json:"version" validate:"required" fake:"{number:1}"`
}

func NewTemplate() Template {
	return Template{
		Entity: model.NewEntity(),
	}
}

func (t Template) Clone() Template {
	c := t.AsNew()
	c.Entity = t.Entity.Clone().(model.Entity)

	return c
}

func (t Template) AsNew() Template {
	return Template{
		Entity:     model.Entity{},
		Website:    t.Website,
		Name:       t.Name,
		Body:       t.Body,
		Assets:     t.Assets.Clone(),
		HTTPHeader: t.HTTPHeader.Clone(),
		Version:    t.Version,
	}
}

func (t Template) Validate() error {
	return templateValidator.Struct(&t)
}
