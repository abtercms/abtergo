package arr

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

// TODO: check multi-error support

type ErrorType string

const (
	ApplicationError           ErrorType = "application error"
	UnknownError               ErrorType = "unknown error"
	ResourceNotFound           ErrorType = "resource not found"
	ResourceNotModified        ErrorType = "resource not modified"
	ResourceIsOutdated         ErrorType = "resource is outdated"
	ETagMismatch               ErrorType = "e-tag mismatch"
	InvalidUserInput           ErrorType = "invalid user input"
	UpstreamServiceUnavailable ErrorType = "upstream service unavailable"
	UpstreamServiceBusy        ErrorType = "upstream service busy"
	TemplateError              ErrorType = "template error"
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
	case ETagMismatch:
		// 409
		return http.StatusPreconditionFailed
	case InvalidUserInput:
		// 400
		return http.StatusBadRequest
	case UpstreamServiceUnavailable:
		// 503
		return http.StatusBadGateway
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
	GetSlug() string
	DetailedError() string
	Attrs() []slog.Attr
	Unwrap() error
}

type arr struct {
	t          ErrorType
	e          error
	attributes []slog.Attr
}

func (a *arr) HTTPStatus() int {
	return a.t.HTTPStatus()
}

func (a *arr) GetSlug() string {
	return a.t.GetSlug()
}

func (a *arr) Error() string {
	return a.e.Error()
}

func (a *arr) DetailedError() string {
	res := a.e.Error() + "."

	if len(a.attributes) == 0 {
		return res
	}

	b := strings.Builder{}
	for _, attr := range a.attributes {
		b.WriteString(" ")
		b.WriteString(attr.Key)
		b.WriteString(": ")
		b.WriteString(attr.Value.String())
		b.WriteString(",")
	}

	str := b.String()

	return res + str[0:len(str)-1]
}

func (a *arr) Attrs() []slog.Attr {
	result := make([]slog.Attr, 0, len(a.attributes)+5)
	result = append(result, slog.Attr{Key: "err", Value: slog.StringValue(a.Error())})
	result = append(result, slog.Attr{Key: "type", Value: slog.StringValue(string(a.t))})
	result = append(result, slog.Attr{Key: "status", Value: slog.IntValue(a.HTTPStatus())})
	result = append(result, a.attributes...)

	return result
}

func (a *arr) Unwrap() error {
	return a.e
}

func Wrap(e error, msg string, attrs ...slog.Attr) Arr {
	t2 := TypeFromError(e)

	return WrapWithType(t2, e, msg, attrs...)
}

func WrapWithFallback(t ErrorType, e error, msg string, attrs ...slog.Attr) Arr {
	t2 := TypeFromError(e)

	if t2 != UnknownError {
		return WrapWithType(t2, e, msg, attrs...)
	}

	return WrapWithType(t, e, msg, attrs...)
}

func WrapWithType(t ErrorType, e error, msg string, attrs ...slog.Attr) Arr {
	if msg != "" {
		e = errors.Wrap(e, msg)
	}

	return newArr(t, e, attrs...)
}

func New(t ErrorType, msg string, args ...slog.Attr) Arr {
	return newArr(t, errors.New(msg), args...)
}

func newArr(t ErrorType, e error, args ...slog.Attr) Arr {
	return &arr{
		t:          t,
		e:          e,
		attributes: args,
	}
}

func HTTPStatusFromError(e error) int {
	if e == nil {
		return http.StatusOK
	}

	var a *arr
	if errors.As(e, &a) {
		return a.HTTPStatus()
	}

	var fe *fiber.Error
	if errors.As(e, &fe) {
		return fe.Code
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
