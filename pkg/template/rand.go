package template

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
)

func init() {
	fakeit.AddWebsiteFaker()
	fakeit.AddCSSURLFaker()
	fakeit.AddJSURLFaker()
	fakeit.AddDateRangeFaker()
	fakeit.AddEtagFaker()
}

// RandomTemplate generates a random Template instance.
func RandomTemplate(asNew bool) Template {
	t := Template{}

	err := gofakeit.Struct(&t)
	if err != nil {
		panic(fmt.Errorf("failed to generate random redirect. err: %w", err))
	}

	if len(t.HTTPHeader) == 0 {
		t.HTTPHeader = nil
	}

	if asNew {
		return t.AsNew()
	}

	return t
}

// RandomTemplateList generates a collection of random Template instances.
func RandomTemplateList(min, max int) []Template {
	templates := []Template{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		templates = append(templates, RandomTemplate(false))
	}

	return templates
}
