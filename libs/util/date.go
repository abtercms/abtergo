package util

import "time"

// CloneDate clones a time.Time instance.
func CloneDate(d time.Time) time.Time {
	year, month, day := d.Date()

	return time.Date(year, month, day, d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
}

func MustParseDate(value, layout string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return t
}
