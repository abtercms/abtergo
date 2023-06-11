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
	Lead       string      `json:"lead" validate:"required" fake:"{paragraph:3}"`
	Title      string      `json:"title" validate:"required" fake:"{sentence:1}"`
	Body       string      `json:"body" validate:"required" fake:"{paragraph:10}"`
	Assets     html.Assets `json:"assets" validate:"dive"`
	HTTPHeader http.Header `json:"http_header" validate:"dive,required"`
	Status     Status      `json:"status" validate:"required,oneof=active inactive draft" fake:"{randomstring:[active,inactive,draft]}"`
	Version    int64       `json:"version" validate:"required" fake:"{number:1}"`
}

func NewPage() Page {
	return Page{
		Entity: model.NewEntity(),
	}
}

func (p Page) Clone() model.EntityInterface {
	c := p.c()
	c.Entity = p.Entity.Clone().(model.Entity)

	return c
}

func (p Page) c() Page {
	return Page{
		Website:    p.Website,
		Path:       p.Path,
		Lead:       p.Lead,
		Title:      p.Title,
		Body:       p.Body,
		Assets:     p.Assets.Clone(),
		HTTPHeader: util.CloneHTTPHeader(p.HTTPHeader),
		Status:     p.Status,
		Version:    p.Version,
	}
}

func (p Page) Validate() error {
	return pageValidator.Struct(p)
}
