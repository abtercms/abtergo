package redirect_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/redirect"
)

func TestFilter_Match(t *testing.T) {
	type fields struct {
		Website string
		Path    string
	}
	type args struct {
		redirect redirect.Redirect
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "match empty filter",
			fields: fields{},
			args: args{
				redirect: redirect.Redirect{Website: "https://example.com", Path: "/path"},
			},
			want: true,
		},
		{
			name: "match with website only",
			fields: fields{
				Website: "https://example.com",
			},
			args: args{
				redirect: redirect.Redirect{Website: "https://example.com", Path: "/path"},
			},
			want: true,
		},
		{
			name: "match with path only",
			fields: fields{
				Website: "https://example.com",
				Path:    "/path",
			},
			args: args{
				redirect: redirect.Redirect{Website: "https://example.com", Path: "/path"},
			},
			want: true,
		},
		{
			name: "match with website and path",
			fields: fields{
				Website: "https://example.com",
				Path:    "/path",
			},
			args: args{
				redirect: redirect.Redirect{Website: "https://example.com", Path: "/path"},
			},
			want: true,
		},
		{
			name: "no match - path over-defined",
			fields: fields{
				Website: "https://example.com",
				Path:    "/path-over-defined",
			},
			args: args{
				redirect: redirect.Redirect{Website: "https://example.com", Path: "/path"},
			},
			want: false,
		},
		{
			name: "no match - subdomain mismatch",
			fields: fields{
				Website: "https://https://example.com",
				Path:    "/path",
			},
			args: args{
				redirect: redirect.Redirect{Website: "https://example.com", Path: "/path"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := redirect.Filter{
				Website: tt.fields.Website,
				Path:    tt.fields.Path,
			}

			result := f.Match(context.Background(), tt.args.redirect)

			assert.Equalf(t, tt.want, result, "Match(%v, %v)", context.Background(), tt.args.redirect)
		})
	}
}
