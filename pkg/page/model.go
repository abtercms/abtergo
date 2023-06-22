package page

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/util"
	"github.com/abtergo/abtergo/libs/validation"
)

var pageValidator *validator.Validate

func init() {
	v := validation.NewValidator()

	pageValidator = v
}

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDraft    Status = "draft"
)

type Page struct {
	model.Entity

	Website    string      `json:"website" validate:"required,url" fake:"{website}"`
	Path       string      `json:"path" validate:"required,path" fake:"{path}"`
	Title      string      `json:"title" validate:"required" fake:"{sentence:1}"`
	Lead       string      `json:"lead" validate:"required" fake:"{paragraph:3}"`
	Body       string      `json:"body" validate:"required" fake:"{paragraph:10}"`
	Assets     html.Assets `json:"assets,omitempty" validate:"dive"`
	HTTPHeader http.Header `json:"http_header,omitempty" validate:"dive,required"`
	Status     Status      `json:"status" validate:"required,oneof=active inactive draft" fake:"{randomstring:[active,inactive,draft]}"`
}

func (p Page) Clone() model.EntityInterface {
	return Page{
		Entity:     p.Entity.Clone().(model.Entity),
		Website:    p.Website,
		Path:       p.Path,
		Lead:       p.Lead,
		Title:      p.Title,
		Body:       p.Body,
		Assets:     p.Assets.Clone(),
		HTTPHeader: util.CloneHTTPHeader(p.HTTPHeader),
		Status:     p.Status,
	}
}

func (p Page) Render() string {
	return p.Body
}

func (p Page) GetContext() []any {
	return []any{p}
}

func (p Page) GetTags() []string {
	return []string{"page-" + p.ID}
}

func (p Page) GetUniqueKey() string {
	return util.Key(p.Website, p.Path)
}

func (p Page) Validate() error {
	return pageValidator.Struct(p)
}
