package redirect

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/abtergo/abtergo/libs/fakeit"
)

func init() {
	fakeit.AddPathFaker()
	fakeit.AddWebsiteFaker()
	fakeit.AddDateRangeFaker()
	fakeit.AddEtagFaker()
}

// RandomRedirect generates a random Redirect instance.
func RandomRedirect() Redirect {
	r := Redirect{}

	err := gofakeit.Struct(&r)
	if err != nil {
		panic(fmt.Errorf("failed to generate random redirect. err: %w", err))
	}

	return r
}

// RandomRedirectList generates a collection of random Redirect instances.
func RandomRedirectList(min, max int) []Redirect {
	redirects := []Redirect{}

	for i := 0; i < gofakeit.Number(min, max); i++ {
		redirects = append(redirects, RandomRedirect())
	}

	return redirects
}
