package page

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/pkg/html"
	"github.com/abtergo/abtergo/pkg/template"
)

func init() {
	fakeit.AddPathFaker()
	fakeit.AddWebsiteFaker()
	fakeit.AddCSSURLFaker()
	fakeit.AddJsURLFaker()
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
	p.Template = nil
	p.TemplateName = ""
	p.TemplateTempName = ""
	p.TemplateTempFrom = nil
	p.TemplateTempUntil = nil

	return p
}

// RandomPageWithTemplate generates a random Page instance with a template.
func RandomPageWithTemplate() Page {
	p := mustCreatePage()

	p.Assets = html.FixAssets(p.Assets)
	p.Template = template.FixTemplate(p.Template)

	p.TemplateName = p.Template.Name
	p.TemplateTempName = ""
	p.TemplateTempFrom = nil
	p.TemplateTempUntil = nil

	return p
}

// RandomPageWithTempTemplate generates a random Page instance with temporary template.
func RandomPageWithTempTemplate() Page {
	p := mustCreatePage()

	p.Assets = html.FixAssets(p.Assets)
	p.Template = template.FixTemplate(p.Template)

	p.TemplateTempName = p.Template.Name

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
