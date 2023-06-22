package decoder

import (
	"errors"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/utils"
)

type Decoder struct {
	JSONDecoder utils.JSONUnmarshal `json:"-"`
}

func NewDecoder() *Decoder {
	return &Decoder{
		JSONDecoder: json.Unmarshal,
	}
}

func (d *Decoder) Decode(data []byte, targets ...any) error {
	var errs error
	for _, target := range targets {
		err := d.JSONDecoder(data, target)
		if errs == nil {
			errs = err
		} else {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
