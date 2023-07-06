package arr

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	AsZapFields() []zap.Field
	Unwrap() error
}

type arr struct {
	t    ErrorType
	e    error
	args []zap.Field
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

	if len(a.args) == 0 {
		return res
	}

	b := strings.Builder{}
	for _, arg := range a.args {
		b.WriteString(" ")
		b.WriteString(arg.Key)
		b.WriteString(": ")
		b.WriteString(a.argToString(arg))
		b.WriteString(",")
	}

	str := b.String()

	return res + str[0:len(str)-1]
}

func (a *arr) AsZapFields() []zap.Field {
	result := make([]zap.Field, 0, len(a.args)+5)
	result = append(result, zap.Error(a))
	result = append(result, zap.String("type", string(a.t)))
	result = append(result, zap.Int("status", a.t.HTTPStatus()))
	result = append(result, a.args...)

	return result
}

func (a *arr) boolArgToString(arg zap.Field) string {
	if arg.Integer == 1 {
		return "true"
	}

	return "false"
}

func (a *arr) timeArgToString(arg zap.Field) string {
	loc := arg.Interface.(*time.Location)
	sec := arg.Integer / int64(time.Second)
	nsec := arg.Integer % int64(time.Second)
	t := time.Unix(sec, nsec).In(loc)

	if nsec > 0 {
		return t.Format(time.RFC3339Nano)
	}

	return t.Format(time.RFC3339)
}

func (a *arr) argToString(arg zap.Field) string {
	switch arg.Type {
	case zapcore.StringType:
		return arg.String
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type, zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
		return fmt.Sprintf("%d", arg.Integer)
	case zapcore.Float32Type:
		return fmt.Sprintf("%g", math.Float32frombits(uint32(arg.Integer)))
	case zapcore.Float64Type:
		return fmt.Sprintf("%g", math.Float64frombits(uint64(arg.Integer)))
	case zapcore.BoolType:
		return a.boolArgToString(arg)
	case zapcore.TimeType:
		return a.timeArgToString(arg)
	}

	return fmt.Sprintf("%v", arg.Interface)
}

func (a *arr) Unwrap() error {
	return a.e
}

func Wrap(e error, msg string, args ...zap.Field) Arr {
	t2 := TypeFromError(e)

	return WrapWithType(t2, e, msg, args...)
}

func WrapWithFallback(t ErrorType, e error, msg string, args ...zap.Field) Arr {
	t2 := TypeFromError(e)

	if t2 != UnknownError {
		return WrapWithType(t2, e, msg, args...)
	}

	return WrapWithType(t, e, msg, args...)
}

func WrapWithType(t ErrorType, e error, msg string, args ...zap.Field) Arr {
	if msg != "" {
		e = errors.Wrap(e, msg)
	}

	return newArr(t, e, args...)
}

func New(t ErrorType, msg string, args ...zap.Field) Arr {
	return newArr(t, errors.New(msg), args...)
}

func newArr(t ErrorType, e error, args ...zap.Field) Arr {
	return &arr{
		t:    t,
		e:    e,
		args: args,
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
