package util

// Cloner is an interface for types which can clone themselves.
type Cloner[T any] interface {
	Clone() T
}

// Clone is an implementation of cloning function for slices (collections).
func Clone[T Cloner[T]](list []T) []T {
	if list == nil {
		return nil
	}

	c := make([]T, 0, len(list))

	for key := range list {
		c = append(c, list[key].Clone())
	}

	return c
}
