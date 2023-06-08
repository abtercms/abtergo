package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// NewValidator creates validator with proper configuration for json tag name extraction.
func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(jsonTagName)

	return v
}

// jsonTagName returns the json name of a property.
func jsonTagName(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 1)[0]
	if name == "-" {
		return ""
	}

	return strings.SplitN(name, ",", 2)[0]
}

func AddNotBeforeValidation(v *validator.Validate) {
	// Register custom validators here
	err := v.RegisterValidation("not_before_date", ValidateNotBeforeDate)
	if err != nil {
		panic(fmt.Errorf("failed to register 'not before date' validator. err: %w", err))
	}
}

func AddEtagValidation(v *validator.Validate) {
	err := v.RegisterValidation("etag", ValidateEtag)
	if err != nil {
		panic(fmt.Errorf("failed to register 'etag' validator. err: %w", err))
	}
}

// ValidateNotBeforeDate validates that a date is not before a reference date.
func ValidateNotBeforeDate(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {
	case reflect.Struct:
		timeType := reflect.TypeOf(time.Time{})

		if field.Type().ConvertibleTo(timeType) {
			p, err := time.Parse(time.DateOnly, param)
			if err != nil {
				panic(fmt.Sprintf("Invalid date: %s", param))
			}

			t := field.Convert(timeType).Interface().(time.Time)

			return t.After(p) || t.Equal(p) || t.Equal(time.Time{})
		}
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// ValidateEtag validates and etag.
func ValidateEtag(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	if val == "" {
		return true
	}

	const pattern = "^[a-z0-9]{5}$"

	match, err := regexp.MatchString(pattern, val)
	if err != nil {
		panic(fmt.Sprintf("Invalid pattern: %s", pattern))
	}

	return match
}
