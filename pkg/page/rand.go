package page

import (
	"sync"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"

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

var lock sync.Mutex

// RandomPage generates a random Page instance.
func RandomPage() Page {
	p := mustCreatePage()

	return p
}

func mustCreatePage() Page {
	lock.Lock()
	defer lock.Unlock()

	p := Page{}

	err := gofakeit.Struct(&p)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate random page"))
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
