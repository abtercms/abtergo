package model

import (
	"strings"
)

func KeyFromStrings(values ...string) Key {
	return Key(ETagFromString(strings.Trim(strings.Join(values, "-"), "-")))
}
