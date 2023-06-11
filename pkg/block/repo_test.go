package block_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/block"
)

func TestFilter_Match(t *testing.T) {
	ctx := context.Background()

	type fields struct {
		Website string
		Name    string
	}
	type args struct {
		block block.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "match website",
			fields: fields{
				Website: "www.example.com",
			},
			args: args{
				block: block.Block{
					Website: "www.example.com",
					Name:    "test",
				},
			},
			want: true,
		},
		{
			name: "match name",
			fields: fields{
				Name: "test",
			},
			args: args{
				block: block.Block{
					Website: "www.example.com",
					Name:    "test",
				},
			},
			want: true,
		},
		{
			name: "match website and name",
			fields: fields{
				Website: "www.example.com",
				Name:    "test",
			},
			args: args{
				block: block.Block{
					Website: "www.example.com",
					Name:    "test",
				},
			},
			want: true,
		},
		{
			name: "no match name",
			fields: fields{
				Website: "www.example.com",
				Name:    "test",
			},
			args: args{
				block: block.Block{
					Website: "www.example.com",
					Name:    "test2",
				},
			},
			want: false,
		},
		{
			name: "no match website",
			fields: fields{
				Website: "www.example.com",
				Name:    "test",
			},
			args: args{
				block: block.Block{
					Website: "www.example.com2",
					Name:    "test",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := block.Filter{
				Website: tt.fields.Website,
				Name:    tt.fields.Name,
			}
			match, err := f.Match(ctx, tt.args.block)
			assert.NoError(t, err)
			assert.Equalf(t, tt.want, match, "Match(%v, %v)", ctx, tt.args.block)
		})
	}
}
