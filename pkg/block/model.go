package block

import (
	"github.com/go-playground/validator/v10"

	"github.com/abtergo/abtergo/libs/html"
	"github.com/abtergo/abtergo/libs/model"
	"github.com/abtergo/abtergo/libs/validation"
)

var blockValidator *validator.Validate

func init() {
	v := validation.NewValidator()

	blockValidator = v
}

// Block represents a resource ready which can be used as wrapper for page content.
type Block struct {
	model.Entity `validate:"dive"`

	Website string      `json:"website" validate:"required,url" fake:"{website}"`
	Name    string      `json:"name" validate:"required" fake:"{sentence:1}"`
	Body    string      `json:"body" validate:"required" fake:"{paragraph:5}"`
	Assets  html.Assets `json:"assets" validate:"dive"`
	Version int64       `json:"version" validate:"required" fake:"{number:1}"`
}

func NewBlock(entity model.Entity) Block {
	return Block{
		Entity: entity,
	}
}

func (b Block) Clone() model.EntityInterface {
	c := b.c()
	c.Entity = b.Entity.Clone().(model.Entity)

	return c
}

func (b Block) c() Block {
	return Block{
		Website: b.Website,
		Name:    b.Name,
		Body:    b.Body,
		Assets:  b.Assets.Clone(),
		Version: b.Version,
	}
}

// Validate validates the entity.
func (b Block) Validate() error {
	return blockValidator.Struct(&b)
}
