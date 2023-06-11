package block

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/fakeit"
)

func init() {
	fakeit.AddWebsiteFaker()
	fakeit.AddCSSURLFaker()
	fakeit.AddJSURLFaker()
	fakeit.AddDateRangeFaker()
	fakeit.AddEtagFaker()
}

// RandomBlock generates a random Block instance.
func RandomBlock(asNew bool) Block {
	b := NewBlock()

	err := gofakeit.Struct(&b)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate random redirect"))
	}

	if asNew {
		return b.AsNew().(Block)
	}

	return b
}

// RandomBlockList generates a collection of random Block instances.
func RandomBlockList(min, max int) []Block {
	blocks := []Block{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		blocks = append(blocks, RandomBlock(false))
	}

	return blocks
}
