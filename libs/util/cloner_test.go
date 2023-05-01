package util_test

import (
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
