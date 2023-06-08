package page

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/html"
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
	switch gofakeit.Number(0, 2) {
	case 1:
		return RandomPageWithoutTemplate()
	case 2:
		return RandomPageWithTempTemplate()
	}

	return RandomPageWithTemplate()
}

// RandomPageWithoutTemplate generates a random Page instance without a template.
func RandomPageWithoutTemplate() Page {
	p := mustCreatePage()

	p.Assets = html.FixAssets(p.Assets)

	return p
}

// RandomPageWithTemplate generates a random Page instance with a template.
func RandomPageWithTemplate() Page {
	p := mustCreatePage()

	p.Assets = html.FixAssets(p.Assets)

	return p
}

// RandomPageWithTempTemplate generates a random Page instance with temporary template.
func RandomPageWithTempTemplate() Page {
	p := mustCreatePage()

	p.Assets = html.FixAssets(p.Assets)

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
