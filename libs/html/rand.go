package html

import "github.com/brianvoe/gofakeit/v6"

// FixAssets ensures that randomly generated assets pass validation.
func FixAssets(a Assets) Assets {
	for i := range a.HeaderMeta {
		if gofakeit.Bool() {
			a.HeaderMeta[i].Name = ""
		} else {
			a.HeaderMeta[i].Property = ""
		}
	}

	return a
}
