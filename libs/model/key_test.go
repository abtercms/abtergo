package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/libs/model"
)

func TestKey(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want model.Key
	}{
		{
			name: "foo",
			args: args{
				values: []string{"foo"},
			},
			want: "0beec",
		},
		{
			name: "bar",
			args: args{
				values: []string{"bar"},
			},
			want: "62cdb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, model.KeyFromStrings(tt.args.values...), "Key(%v)", tt.args.values)
		})
	}
}
