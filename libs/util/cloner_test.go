package util_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/util"
)

type foo int

func (f foo) Clone() foo {
	return f + 1
}

func TestClone(t *testing.T) {
	type args[T util.Cloner[T]] struct {
		list []T
	}
	type testCase[T util.Cloner[T]] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[foo]{
		{
			name: "nil",
			args: args[foo]{
				list: nil,
			},
			want: nil,
		},
		{
			name: "empty list",
			args: args[foo]{
				list: []foo{},
			},
			want: []foo{},
		},
		{
			name: "non-empty list",
			args: args[foo]{
				list: []foo{1, 2, 3},
			},
			want: []foo{2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, util.Clone(tt.args.list), "Clone(%v)", tt.args.list)
		})
	}
}

func TestCloneStrings(t *testing.T) {
	type args struct {
		list []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nil",
			args: args{
				list: nil,
			},
			want: nil,
		},
		{
			name: "empty list",
			args: args{
				list: []string{},
			},
			want: []string{},
		},
		{
			name: "non-empty list",
			args: args{
				list: []string{"1", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, util.CloneStrings(tt.args.list), "CloneStrings(%v)", tt.args.list)
		})
	}
}

func TestCloneHttpHeader(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want http.Header
	}{
		{
			name: "nil",
			args: args{
				header: nil,
			},
			want: nil,
		},
		{
			name: "empty list",
			args: args{
				header: make(http.Header),
			},
			want: make(http.Header),
		},
		{
			name: "non-empty list",
			args: args{
				header: http.Header{"foo": []string{"bar"}},
			},
			want: http.Header{"foo": []string{"bar"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, util.CloneHTTPHeader(tt.args.header), "CloneHTTPHeader(%v)", tt.args.header)
		})
	}
}
