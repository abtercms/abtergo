package block

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/model"
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
		panic(fmt.Errorf("failed to generate random redirect. err: %w", err))
	}

	if asNew {
		b.Entity = model.Entity{}
	} else {
		b.Entity = model.FixEntity(b.Entity)
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
