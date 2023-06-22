package template

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/util"
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
}

func (t Template) Clone() model.EntityInterface {
	return Template{
		Entity:     t.Entity.Clone().(model.Entity),
		Website:    t.Website,
		Name:       t.Name,
		Body:       t.Body,
		Assets:     t.Assets.Clone(),
		HTTPHeader: t.HTTPHeader.Clone(),
	}
}

func (t Template) Validate() error {
	return templateValidator.Struct(&t)
}

func (t Template) Render() string {
	return t.Body
}

func (t Template) GetContext() []any {
	return []any{t}
}

func (t Template) GetTags() []string {
	return []string{"template-" + t.ID}
}

func (t Template) GetUniqueKey() string {
	return util.Key(t.Website, t.Name)
}
