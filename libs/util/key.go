package util

import (
	"strings"
)

func Key(values ...string) string {
	return ETag(strings.Trim(strings.Join(values, "-"), "-"))
}
