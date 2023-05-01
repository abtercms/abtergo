package util

import "time"

// CloneDate clones a time.Time instance.
func CloneDate(d time.Time) time.Time {
	year, month, day := d.Date()

	return time.Date(year, month, day, d.Hour(), d.Minute(), d.Second(), d.Nanosecond(), d.Location())
}
