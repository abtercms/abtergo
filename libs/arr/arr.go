package arr

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

type ErrorType string

const (
	UnknownError               ErrorType = "unknown error"
	ResourceNotFound           ErrorType = "resource not found"
	ResourceNotModified        ErrorType = "resource not modified"
	ResourceIsOutdated         ErrorType = "resource is outdated"
	InvalidEtag                ErrorType = "invalid etag"
	InvalidUserInput           ErrorType = "invalid user input"
	UpstreamServiceUnavailable ErrorType = "upstream service unavailable"
	UpstreamServiceBusy        ErrorType = "upstream service busy"
)

func (et ErrorType) HTTPStatus() int {
	switch et {
	case ResourceNotFound:
		// 404
		return http.StatusNotFound
	case ResourceNotModified:
		// 403
		return http.StatusNotModified
	case ResourceIsOutdated:
		// 409
		return http.StatusConflict
	case InvalidEtag:
		// 409
		return http.StatusConflict
	case InvalidUserInput:
		// 400
		return http.StatusBadRequest
	case UpstreamServiceUnavailable:
		// 503
		return http.StatusServiceUnavailable
	case UpstreamServiceBusy:
		// 429
		return http.StatusTooManyRequests
	}

	// 500
	return http.StatusInternalServerError
}

func (et ErrorType) GetTitle() string {
	return string(et)
}

func (et ErrorType) GetSlug() string {
	return slug.Make(et.GetTitle())
}

type Arr interface {
	error

	HTTPStatus() int
}

type arr struct {
	t    ErrorType
	e    error
	args []string
}

func (a arr) HTTPStatus() int {
	return a.t.HTTPStatus()
}

func (a arr) GetSlug() string {
	return "tbd"
}

func (a arr) Error() string {
	res := a.e.Error() + "."
	if len(a.args) > 0 {
		res += " " + strings.Join(a.args, ", ")
	}

	return res
}

func Wrap(t ErrorType, e error, msg string, args ...interface{}) Arr {
	if len(args)%2 != 0 {
		panic("invalid args")
	}

	a := make([]string, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		a = append(a, fmt.Sprintf("%s: %v", args[i], args[i+1]))
	}

	if msg != "" {
		e = errors.Wrap(e, msg)
	}

	return &arr{
		e:    e,
		t:    t,
		args: a,
	}
}

func New(t ErrorType, msg string, args ...interface{}) Arr {
	if len(args)%2 != 0 {
		panic("invalid args")
	}

	a := make([]string, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		a = append(a, fmt.Sprintf("%s: %v", args[i], args[i+1]))
	}

	return &arr{
		e:    errors.New(msg),
		t:    t,
		args: a,
	}
}

func HTTPStatusFromError(e error) int {
	if e == nil {
		return http.StatusInternalServerError
	}

	var a *arr

	if errors.As(e, &a) {
		return a.HTTPStatus()
	}

	return http.StatusInternalServerError
}

func TypeFromError(e error) ErrorType {
	if e == nil {
		return UnknownError
	}

	var a *arr

	if errors.As(e, &a) {
		return a.t
	}

	return UnknownError
}
