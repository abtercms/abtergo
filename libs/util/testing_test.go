package util

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataToReaderHelper(t *testing.T) {
	type foo struct {
		A int `json:"bar"`
	}

	type args struct {
		t *testing.T
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want io.Reader
	}{
		{
			name: "default",
			args: args{
				t: &testing.T{},
				v: foo{A: 123},
			},
			want: strings.NewReader(`{"bar":123}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, DataToReaderHelper(tt.args.t, tt.args.v), "DataToReaderHelper(%v, %v)", tt.args.t, tt.args.v)
		})
	}
}

func TestParseResponseHelper(t *testing.T) {
	type foo struct {
		A int `json:"bar"`
	}

	type args[V any] struct {
		t  *testing.T
		rc io.Reader
		v  V
	}
	type testCase[V any] struct {
		name string
		args args[V]
		want V
	}
	tests := []testCase[foo]{
		{
			name: "default",
			args: args[foo]{
				t:  &testing.T{},
				rc: strings.NewReader(`{"bar":123}`),
				v:  foo{},
			},
			want: foo{
				A: 123,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{Body: io.NopCloser(tt.args.rc)}

			ParseResponseHelper(tt.args.t, resp, &tt.args.v)

			assert.Equal(t, tt.want, tt.args.v)
		})
	}
}
