package arr_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/arr"
)

func TestWrap(t *testing.T) {
	type args struct {
		e     error
		msg   string
		attrs []slog.Attr
	}
	tests := []struct {
		name         string
		args         args
		wantDetailed string
		wantError    string
		wantAttrs    []slog.Attr
	}{
		{
			name: "assert.AnError wrapped",
			args: args{
				e:     assert.AnError,
				msg:   "outdated resource",
				attrs: []slog.Attr{{Key: "id", Value: slog.StringValue("foo")}},
			},
			wantDetailed: "outdated resource: assert.AnError general error for testing. id: foo",
			wantError:    "outdated resource: assert.AnError general error for testing",
			wantAttrs: []slog.Attr{
				{Key: "err", Value: slog.StringValue(errors.Wrap(assert.AnError, "outdated resource").Error())},
				{Key: "type", Value: slog.StringValue(string(arr.UnknownError))},
				{Key: "status", Value: slog.IntValue(http.StatusInternalServerError)},
				{Key: "id", Value: slog.StringValue("foo")},
			},
		},
		{
			name: "arr.Arr wrapped",
			args: args{
				e:     arr.New(arr.ResourceNotFound, "not found"),
				msg:   "outdated resource",
				attrs: []slog.Attr{{Key: "id", Value: slog.StringValue("foo")}},
			},
			wantDetailed: "outdated resource: not found. id: foo",
			wantError:    "outdated resource: not found",
			wantAttrs: []slog.Attr{
				{Key: "err", Value: slog.StringValue(errors.Wrap(arr.New(arr.ResourceNotFound, "not found"), "outdated resource").Error())},
				{Key: "type", Value: slog.StringValue(string(arr.ResourceNotFound))},
				{Key: "status", Value: slog.IntValue(http.StatusNotFound)},
				{Key: "id", Value: slog.StringValue("foo")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := arr.Wrap(tt.args.e, tt.args.msg, tt.args.attrs...)

			assert.Equal(t, tt.wantDetailed, sut.DetailedError())
			assert.Equal(t, tt.wantError, sut.Error())
			assert.Equal(t, tt.wantAttrs, sut.Attrs())
		})
	}
}

func TestWrapWithFallback(t *testing.T) {
	type args struct {
		t     arr.ErrorType
		e     error
		msg   string
		attrs []slog.Attr
	}
	tests := []struct {
		name         string
		args         args
		wantDetailed string
		wantError    string
		wantAttrs    []slog.Attr
	}{
		{
			name: "assert.AnError wrapped",
			args: args{
				t:     arr.ResourceIsOutdated,
				e:     assert.AnError,
				msg:   "outdated resource",
				attrs: []slog.Attr{{Key: "id", Value: slog.StringValue("foo")}},
			},
			wantDetailed: "outdated resource: assert.AnError general error for testing. id: foo",
			wantError:    "outdated resource: assert.AnError general error for testing",
			wantAttrs: []slog.Attr{
				{Key: "err", Value: slog.StringValue(errors.Wrap(assert.AnError, "outdated resource").Error())},
				{Key: "type", Value: slog.StringValue(string(arr.ResourceIsOutdated))},
				{Key: "status", Value: slog.IntValue(http.StatusConflict)},
				{Key: "id", Value: slog.StringValue("foo")},
			},
		},
		{
			name: "arr.Arr wrapped",
			args: args{
				t:     arr.ResourceIsOutdated,
				e:     arr.New(arr.ResourceNotFound, "not found"),
				msg:   "outdated resource",
				attrs: []slog.Attr{{Key: "id", Value: slog.StringValue("foo")}},
			},
			wantDetailed: "outdated resource: not found. id: foo",
			wantError:    "outdated resource: not found",
			wantAttrs: []slog.Attr{
				{Key: "err", Value: slog.StringValue(errors.Wrap(arr.New(arr.ResourceNotFound, "not found"), "outdated resource").Error())},
				{Key: "type", Value: slog.StringValue(string(arr.ResourceNotFound))},
				{Key: "status", Value: slog.IntValue(http.StatusNotFound)},
				{Key: "id", Value: slog.StringValue("foo")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := arr.WrapWithFallback(tt.args.t, tt.args.e, tt.args.msg, tt.args.attrs...)

			assert.Equal(t, tt.wantDetailed, sut.DetailedError())
			assert.Equal(t, tt.wantError, sut.Error())
			gotAttrs := sut.Attrs()
			assert.Equal(t, tt.wantAttrs, gotAttrs)
		})
	}
}

func TestWrapWithType(t *testing.T) {
	type args struct {
		t     arr.ErrorType
		e     error
		msg   string
		attrs []slog.Attr
	}
	tests := []struct {
		name          string
		args          args
		wantDetailed  string
		wantError     string
		wantAttrs     []slog.Attr
		wantErrorType arr.ErrorType
	}{
		{
			name: "error type",
			args: args{
				t:     arr.ResourceIsOutdated,
				e:     assert.AnError,
				msg:   "outdated resource",
				attrs: []slog.Attr{{Key: "id", Value: slog.StringValue("foo")}},
			},
			wantDetailed: "outdated resource: assert.AnError general error for testing. id: foo",
			wantError:    "outdated resource: assert.AnError general error for testing",
			wantAttrs: []slog.Attr{
				{Key: "err", Value: slog.StringValue(errors.Wrap(assert.AnError, "outdated resource").Error())},
				{Key: "type", Value: slog.StringValue(string(arr.ResourceIsOutdated))},
				{Key: "status", Value: slog.IntValue(http.StatusConflict)},
				{Key: "id", Value: slog.StringValue("foo")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := arr.WrapWithType(tt.args.t, tt.args.e, tt.args.msg, tt.args.attrs...)

			assert.Equal(t, tt.wantDetailed, sut.DetailedError())
			assert.Equal(t, tt.wantError, sut.Error())
			assert.Equal(t, tt.wantAttrs, sut.Attrs())
		})
	}
}

func TestErrorType_HTTPStatus(t *testing.T) {
	tests := []struct {
		name string
		et   arr.ErrorType
		want int
	}{
		{
			name: "unknown error",
			et:   arr.UnknownError,
			want: http.StatusInternalServerError,
		},
		{
			name: "resource not found",
			et:   arr.ResourceNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "resource not modified",
			et:   arr.ResourceNotModified,
			want: http.StatusNotModified,
		},
		{
			name: "resource is outdated",
			et:   arr.ResourceIsOutdated,
			want: http.StatusConflict,
		},
		{
			name: "invalid e-tag",
			et:   arr.ETagMismatch,
			want: http.StatusPreconditionFailed,
		},
		{
			name: "invalid user input",
			et:   arr.InvalidUserInput,
			want: http.StatusBadRequest,
		},
		{
			name: "upstream service unavailable",
			et:   arr.UpstreamServiceUnavailable,
			want: http.StatusBadGateway,
		},
		{
			name: "upstream service busy",
			et:   arr.UpstreamServiceBusy,
			want: http.StatusTooManyRequests,
		},
		{
			name: "custom error",
			et:   arr.ErrorType("custom"),
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.et.HTTPStatus(), "HTTPStatus()")
		})
	}
}

func TestErrorType_GetSlug(t *testing.T) {
	tests := []struct {
		name string
		et   arr.ErrorType
		want string
	}{
		{
			name: "unknown error",
			et:   arr.UnknownError,
			want: "unknown-error",
		},
		{
			name: "resource not found",
			et:   arr.ResourceNotFound,
			want: "resource-not-found",
		},
		{
			name: "resource not modified",
			et:   arr.ResourceNotModified,
			want: "resource-not-modified",
		},
		{
			name: "resource is outdated",
			et:   arr.ResourceIsOutdated,
			want: "resource-is-outdated",
		},
		{
			name: "invalid e-tag",
			et:   arr.ETagMismatch,
			want: "e-tag-mismatch",
		},
		{
			name: "invalid user input",
			et:   arr.InvalidUserInput,
			want: "invalid-user-input",
		},
		{
			name: "upstream service unavailable",
			et:   arr.UpstreamServiceUnavailable,
			want: "upstream-service-unavailable",
		},
		{
			name: "upstream service busy",
			et:   arr.UpstreamServiceBusy,
			want: "upstream-service-busy",
		},
		{
			name: "custom error",
			et:   arr.ErrorType("custom"),
			want: "custom",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.et.GetSlug(), "GetSlug()")
		})
	}
}

func Test_arr_HttpStatus(t *testing.T) {
	type fields struct {
		t arr.ErrorType
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "unknown error",
			fields: fields{
				t: arr.UnknownError,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "not found",
			fields: fields{
				t: arr.ResourceNotFound,
			},
			want: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := arr.New(tt.fields.t, tt.name)
			if got := a.HTTPStatus(); got != tt.want {
				t.Errorf("HTTPStatus() = %v, wantDetailed %v", got, tt.want)
			}
		})
	}
}

func Test_arr_GetSlug(t *testing.T) {
	type fields struct {
		t     arr.ErrorType
		msg   string
		attrs []slog.Attr
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "unknown error",
			fields: fields{
				t: arr.UnknownError,
			},
			want: "unknown-error",
		},
		{
			name: "not found",
			fields: fields{
				t:   arr.ResourceNotFound,
				msg: "foo",
				attrs: []slog.Attr{
					{Key: "bar", Value: slog.StringValue("quix")},
				},
			},
			want: "resource-not-found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := arr.New(tt.fields.t, tt.fields.msg, tt.fields.attrs...)
			assert.Equalf(t, tt.want, a.GetSlug(), "GetSlug()")
		})
	}
}

func Test_arr_DetailedError(t *testing.T) {
	type fields struct {
		t     arr.ErrorType
		e     error
		msg   string
		attrs []slog.Attr
	}
	tests := []struct {
		name         string
		fields       fields
		wantDetailed string
		wantError    string
	}{
		{
			name: "simple error",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
			},
			wantDetailed: "foo: " + assert.AnError.Error() + ".",
			wantError:    "foo: " + assert.AnError.Error(),
		},
		{
			name: "numbers",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
				attrs: []slog.Attr{
					{Key: "i", Value: slog.IntValue(-1)},
					{Key: "i64", Value: slog.Int64Value(-101)},
					{Key: "u64", Value: slog.Uint64Value(89)},
					{Key: "f64", Value: slog.Float64Value(64.91)},
				},
			},
			wantDetailed: "foo: " + assert.AnError.Error() + ". i: -1, i64: -101, u64: 89, f64: 64.91",
			wantError:    "foo: " + assert.AnError.Error(),
		},
		{
			name: "non-number scalars",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
				attrs: []slog.Attr{
					{Key: "greeting", Value: slog.StringValue("hello")},
					{Key: "is_ok", Value: slog.BoolValue(true)},
				},
			},
			wantDetailed: "foo: " + assert.AnError.Error() + ". greeting: hello, is_ok: true",
			wantError:    "foo: " + assert.AnError.Error(),
		},
		{
			name: "complex",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
				attrs: []slog.Attr{
					{Key: "numbers", Value: slog.AnyValue([]int{1, 2, 3})},
					{Key: "strings", Value: slog.AnyValue([]string{"foo", "bar"})},
				},
			},
			wantDetailed: "foo: " + assert.AnError.Error() + ". numbers: [1 2 3], strings: [foo bar]",
			wantError:    "foo: " + assert.AnError.Error(),
		},
		{
			name: "date",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
				attrs: []slog.Attr{
					{Key: "date", Value: slog.TimeValue(time.Date(2030, 1, 2, 3, 4, 5, 6, time.UTC))},
				},
			},
			wantDetailed: "foo: " + assert.AnError.Error() + ". date: 2030-01-02 03:04:05.000000006 +0000 UTC",
			wantError:    "foo: " + assert.AnError.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := arr.WrapWithType(tt.fields.t, tt.fields.e, tt.fields.msg, tt.fields.attrs...)
			assert.Equalf(t, tt.wantDetailed, a.DetailedError(), "DetailedError()")
			assert.Equalf(t, tt.wantError, a.Error(), "Error()")
		})
	}
}

func TestHttpStatusFromError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "no error",
			args: args{
				e: nil,
			},
			want: http.StatusOK,
		},
		{
			name: "non-arr error",
			args: args{
				e: assert.AnError,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "unknown error",
			args: args{
				e: arr.WrapWithType(arr.UnknownError, assert.AnError, "foo"),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "not found error",
			args: args{
				e: arr.WrapWithType(arr.ResourceNotFound, assert.AnError, "foo"),
			},
			want: http.StatusNotFound,
		},
		{
			name: "not found error wrapped var errors.WrapWithType",
			args: args{
				e: errors.Wrap(arr.WrapWithType(arr.ResourceNotFound, assert.AnError, "foo"), "bar"),
			},
			want: http.StatusNotFound,
		},
		{
			name: "not found error wrapped via fmt",
			args: args{
				e: arr.WrapWithType(arr.ResourceNotFound, assert.AnError, "bar"),
			},
			want: http.StatusNotFound,
		},
		{
			name: "not found error double wrapped",
			args: args{
				e: fmt.Errorf("bar, err: %w", errors.Wrap(arr.WrapWithType(arr.ResourceNotFound, assert.AnError, "foo"), "bar")),
			},
			want: http.StatusNotFound,
		},
		{
			name: "fiber error",
			args: args{
				e: fiber.ErrConflict,
			},
			want: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arr.HTTPStatusFromError(tt.args.e); got != tt.want {
				t.Errorf("HTTPStatusFromError() = %v, wantDetailed %v", got, tt.want)
			}
		})
	}
}

func TestTypeFromError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name string
		args args
		want arr.ErrorType
	}{
		{
			name: "no error",
			args: args{
				e: nil,
			},
			want: arr.UnknownError,
		},
		{
			name: "an assert error",
			args: args{
				e: assert.AnError,
			},
			want: arr.UnknownError,
		},
		{
			name: "new error",
			args: args{
				e: arr.New(arr.ResourceIsOutdated, "foo"),
			},
			want: arr.ResourceIsOutdated,
		},
		{
			name: "wrapped error",
			args: args{
				e: arr.WrapWithType(arr.ETagMismatch, assert.AnError, "foo"),
			},
			want: arr.ETagMismatch,
		},
		{
			name: "double wrapped error",
			args: args{
				e: errors.Wrap(arr.WrapWithType(arr.UpstreamServiceBusy, assert.AnError, "foo"), "bar"),
			},
			want: arr.UpstreamServiceBusy,
		},
		{
			name: "fmt wrapped error",
			args: args{
				e: fmt.Errorf("quix. err: %w", arr.WrapWithType(arr.ResourceNotModified, assert.AnError, "foo")),
			},
			want: arr.ResourceNotModified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, arr.TypeFromError(tt.args.e), "TypeFromError(%v)", tt.args.e)
		})
	}
}
