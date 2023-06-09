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
func RandomTemplate() Template {
	t := Template{}

	err := gofakeit.Struct(&t)
	if err != nil {
		panic(fmt.Errorf("failed to generate random redirect. err: %w", err))
	}

	FixTemplate(&t)

	return t
}

// RandomTemplateList generates a collection of random Template instances.
func RandomTemplateList(min, max int) []Template {
	templates := []Template{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		templates = append(templates, RandomTemplate())
	}

	return templates
}

// FixTemplate ensures that randomly generated templates pass validation.
func FixTemplate(t *Template) *Template {
	if t == nil {
		return t
	}

	if len(t.HTTPHeader) == 0 {
		t.HTTPHeader = nil
	}

	return t
}
