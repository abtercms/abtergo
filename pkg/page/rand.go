package page

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
)

func init() {
	fakeit.AddPathFaker()
	fakeit.AddWebsiteFaker()
	fakeit.AddCSSURLFaker()
	fakeit.AddJSURLFaker()
	fakeit.AddDateRangeFaker()
	fakeit.AddEtagFaker()
}

// RandomPage generates a random Page instance.
func RandomPage() Page {
	p := mustCreatePage()

	return p
}

func mustCreatePage() Page {
	p := Page{}

	err := gofakeit.Struct(&p)
	if err != nil {
		panic(fmt.Errorf("failed to generate random page. err: %w", err))
	}

	if len(p.HTTPHeader) == 0 {
		p.HTTPHeader = nil
	}

	return p
}

// RandomPageList generates a collection of random Page instances.
func RandomPageList(min, max int) []Page {
	pages := []Page{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		pages = append(pages, RandomPage())
	}

	return pages
}
