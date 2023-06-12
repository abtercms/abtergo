package arr_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/arr"
)

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
			name: "invalid etag",
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
			name: "invalid etag",
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
				t.Errorf("HTTPStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arr_GetSlug(t *testing.T) {
	type fields struct {
		t    arr.ErrorType
		msg  string
		args []zap.Field
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
				args: []zap.Field{
					zap.String("bar", "quix"),
				},
			},
			want: "resource-not-found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := arr.New(tt.fields.t, tt.fields.msg, tt.fields.args...)
			assert.Equalf(t, tt.want, a.GetSlug(), "GetSlug()")
		})
	}
}

func Test_arr_Error(t *testing.T) {
	type fields struct {
		t    arr.ErrorType
		e    error
		msg  string
		args []zap.Field
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple error",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
			},
			want: "foo: " + assert.AnError.Error() + ".",
		},
		{
			name: "complex error",
			fields: fields{
				t:   arr.UnknownError,
				e:   assert.AnError,
				msg: "foo",
				args: []zap.Field{
					zap.Int8("i8", -1),
					zap.Int16("i16", 16),
					zap.Int32("i32", -17),
					zap.Int64("i64", 101),
					zap.Uint("u", 97),
					zap.Uint8("u8", 83),
					zap.Uint16("u16", 32),
					zap.Uint32("u32", 73),
					zap.Uint64("u64", 89),
					zap.Float32("f32", 123.45),
					zap.Float64("f64", 64.91),
					zap.String("greeting", "hello"),
					zap.Bool("is_ok", true),
					zap.Ints("numbers", []int{1, 2, 3}),
					zap.Strings("foobar", []string{"foo", "bar"}),
				},
			},
			want: "foo: " + assert.AnError.Error() + ". i8: -1, i16: 16, i32: -17, i64: 101, u: 97, u8: 83, u16: 32, u32: 73, u64: 89, f32: 123.45, f64: 64.91, greeting: hello, is_ok: true, numbers: [1 2 3], foobar: [foo bar]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := arr.Wrap(tt.fields.t, tt.fields.e, tt.fields.msg, tt.fields.args...)
			assert.Equalf(t, tt.want, a.Error(), "Error()")
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
			want: http.StatusInternalServerError,
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
				e: arr.Wrap(arr.UnknownError, assert.AnError, "foo"),
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "not found error",
			args: args{
				e: arr.Wrap(arr.ResourceNotFound, assert.AnError, "foo"),
			},
			want: http.StatusNotFound,
		},
		{
			name: "not found error wrapped var errors.Wrap",
			args: args{
				e: errors.Wrap(arr.Wrap(arr.ResourceNotFound, assert.AnError, "foo"), "bar"),
			},
			want: http.StatusNotFound,
		},
		{
			name: "not found error wrapped via fmt",
			args: args{
				e: arr.Wrap(arr.ResourceNotFound, assert.AnError, "bar"),
			},
			want: http.StatusNotFound,
		},
		{
			name: "not found error double wrapped",
			args: args{
				e: fmt.Errorf("bar, err: %w", errors.Wrap(arr.Wrap(arr.ResourceNotFound, assert.AnError, "foo"), "bar")),
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
				t.Errorf("HTTPStatusFromError() = %v, want %v", got, tt.want)
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
				e: arr.Wrap(arr.ETagMismatch, assert.AnError, "foo"),
			},
			want: arr.ETagMismatch,
		},
		{
			name: "double wrapped error",
			args: args{
				e: errors.Wrap(arr.Wrap(arr.UpstreamServiceBusy, assert.AnError, "foo"), "bar"),
			},
			want: arr.UpstreamServiceBusy,
		},
		{
			name: "fmt wrapped error",
			args: args{
				e: fmt.Errorf("quix. err: %w", arr.Wrap(arr.ResourceNotModified, assert.AnError, "foo")),
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
