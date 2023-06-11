package redirect

import (
	"fmt"
	"sync"

	"github.com/brianvoe/gofakeit/v6"

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
		panic(fmt.Errorf("failed to generate random redirect. err: %w", err))
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
