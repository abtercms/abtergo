package redirect

import (
	"sync"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pkg/errors"

	"github.com/abtergo/abtergo/libs/fakeit"
	"github.com/abtergo/abtergo/libs/model"
)

func init() {
	fakeit.AddPathFaker()
	fakeit.AddWebsiteFaker()
	fakeit.AddDateRangeFaker()
	fakeit.AddEtagFaker()
}

var lock sync.Mutex

// RandomRedirect generates a random Redirect instance.
func RandomRedirect(asNew bool) Redirect {
	lock.Lock()
	defer lock.Unlock()

	r := Redirect{}

	err := gofakeit.Struct(&r)
	if err != nil {
		panic(errors.Wrap(err, "failed to generate random redirect"))
	}

	if asNew {
		r.Entity = model.Entity{}
	}

	return r
}

// RandomRedirectList generates a collection of random Redirect instances.
func RandomRedirectList(min, max int) []Redirect {
	redirects := []Redirect{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		redirects = append(redirects, RandomRedirect(false))
	}

	return redirects
}
