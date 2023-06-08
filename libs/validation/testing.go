package validation

import (
	"fmt"
	"sort"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertFieldErrorsOn asserts that the ValidationErrors given are in the list of expected fields given
// and that all the fields expected to have an error do.
func AssertFieldErrorsOn(t *testing.T, err error, fields []string) {
	require.Error(t, err)
	fieldErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		assert.Fail(t, fmt.Sprintf("invalid error: %s", err))
	}

	foundFields := make([]string, 0, len(fields))
	for _, fErr := range fieldErrors {
		assert.Contains(t, fields, fErr.Field())

		for _, field := range fields {
			if field == fErr.Field() {
				foundFields = append(foundFields, field)
			}
		}
	}

	sort.Strings(fields)
	sort.Strings(foundFields)
	assert.Equal(t, fields, foundFields)
}
