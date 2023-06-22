package model

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

func init() {
	AddETagFaker()
}

func AddETagFaker() {
	gofakeit.AddFuncLookup("etag", gofakeit.Info{
		Category:    "abtergo",
		Description: "E-tag",
		Example:     "aiso2",
		Output:      "string",
		Params:      []gofakeit.Param{},
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			result := ETagFromString(gofakeit.Word())

			return result, nil
		},
	})
}
