package validation

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// NewValidator creates validator with proper configuration for json tag name extraction.
func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(jsonTagName)

	AddNotBeforeValidation(v)
	AddEtagValidation(v)
	AddPathValidation(v)

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
	mustRegister(err, "not_before_date")
}

func AddEtagValidation(v *validator.Validate) {
	err := v.RegisterValidation("etag", ValidateEtag)
	mustRegister(err, "etag")
}

func AddPathValidation(v *validator.Validate) {
	err := v.RegisterValidation("path", ValidatePath)
	mustRegister(err, "path")
}

func mustRegister(err error, validatorName string) {
	if err == nil {
		return
	}

	panic(errors.Wrapf(err, "failed to register '%s' validator", validatorName))
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
				panic(fmt.Errorf("invalid date: %s", param))
			}

			t := field.Convert(timeType).Interface().(time.Time)

			return t.After(p) || t.Equal(p) || t.Equal(time.Time{})
		}
	}

	panic(fmt.Errorf("bad field type %T", field.Interface()))
}

// ValidateEtag validates an e-tag.
func ValidateEtag(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	if val == "" {
		return true
	}

	const pattern = "^[a-z0-9]{5}$"

	match, err := regexp.MatchString(pattern, val)
	if err != nil {
		panic(errors.Wrapf(err, "invalid pattern: %s", pattern))
	}

	return match
}

// ValidatePath validates a path.
func ValidatePath(fl validator.FieldLevel) bool {
	val := fl.Field().String()

	if val == "" {
		return true
	}

	u, err := url.Parse(val)

	return err == nil && u.Scheme == "" && u.Host == ""
}
