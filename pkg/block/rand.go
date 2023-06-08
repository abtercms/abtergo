package block

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/html"
)

func init() {
	fakeit.AddWebsiteFaker()
	fakeit.AddCSSURLFaker()
	fakeit.AddJSURLFaker()
	fakeit.AddDateRangeFaker()
	fakeit.AddEtagFaker()
}

// RandomBlock generates a random Block instance.
func RandomBlock() Block {
	t := Block{}

	err := gofakeit.Struct(&t)
	if err != nil {
		panic(fmt.Errorf("failed to generate random redirect. err: %w", err))
	}

	FixBlock(&t)

	return t
}

// RandomBlockList generates a collection of random Block instances.
func RandomBlockList(min, max int) []Block {
	blocks := []Block{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		blocks = append(blocks, RandomBlock())
	}

	return blocks
}

// FixBlock ensures that randomly generated blocks pass validation.
func FixBlock(t *Block) *Block {
	if t == nil {
		return t
	}

	t.Assets = html.FixAssets(t.Assets)

	return t
}
