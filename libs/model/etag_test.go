package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/model"
)

func TestETagFromString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want model.ETag
	}{
		{
			name: "foo",
			args: args{
				input: "foo",
			},
			want: "0beec",
		},
		{
			name: "bar",
			args: args{
				input: "bar",
			},
			want: "62cdb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, model.ETagFromString(tt.args.input), "ETagFromString(%v)", tt.args.input)
		})
	}
}

func TestETagFromAny(t *testing.T) {
	t.Run("marshaling failure", func(t *testing.T) {
		type foo struct {
			Foo *foo
			Bar string
		}
		input := foo{Bar: "bar"}
		input.Foo = &input

		assert.Panics(t, func() { model.ETagFromAny(input) })
	})

	t.Run("success", func(t *testing.T) {
		type foo struct {
			Foo string `json:"foo,omitempty"`
			Bar string `json:"bar,omitempty"`
		}
		type args struct {
			input any
		}
		tests := []struct {
			name string
			args args
			want model.ETag
		}{
			{
				name: "foo",
				args: args{
					input: foo{
						Foo: "foo",
					},
				},
				want: "07485",
			},
			{
				name: "bar",
				args: args{
					input: foo{
						Bar: "bar",
					},
				},
				want: "99420",
			},
			{
				name: "default",
				args: args{
					input: foo{
						Foo: "foo",
						Bar: "bar",
					},
				},
				want: "d0228",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, model.ETagFromAny(tt.args.input), "ETagFromAny(%v)", tt.args.input)
			})
		}
	})
}
