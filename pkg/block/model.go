package block

import (
	"time"

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

// Clone clones (duplicates) a Block resource.
func (b Block) Clone() Block {
	return Block{
		Entity:  b.Entity.Clone().(model.Entity),
		Website: b.Website,
		Name:    b.Name,
		Body:    b.Body,
		Assets:  b.Assets.Clone(),
		Version: b.Version,
	}
}

func NewBlock() Block {
	return Block{
		Entity: model.NewEntity(),
	}
}

func (b Block) Reset() Block {
	b.Entity = model.NewEntity()

	return b
}

// WithTime returns the entity but also ensures that it has created at and updated at set.
func (b Block) WithTime() Block {
	if b.CreatedAt.Unix() > 0 && b.UpdatedAt.Unix() > 0 {
		return b
	}

	t2 := time.Now()

	b.CreatedAt = t2
	b.UpdatedAt = t2

	return b
}

// Validate validates the entity.
func (b Block) Validate() error {
	return blockValidator.Struct(&b)
}
