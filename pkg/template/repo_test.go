package template_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/template"
)

func TestFilter_Match(t *testing.T) {
	type fields struct {
		Website string
		Name    string
	}
	type args struct {
		template template.Template
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
				template: template.Template{Website: "https://example.com", Name: "Foo Template"},
			},
			want: true,
		},
		{
			name: "match with website only",
			fields: fields{
				Website: "https://example.com",
			},
			args: args{
				template: template.Template{Website: "https://example.com", Name: "Foo Template"},
			},
			want: true,
		},
		{
			name: "match with path only",
			fields: fields{
				Website: "https://example.com",
				Name:    "Foo Template",
			},
			args: args{
				template: template.Template{Website: "https://example.com", Name: "Foo Template"},
			},
			want: true,
		},
		{
			name: "match with website and path",
			fields: fields{
				Website: "https://example.com",
				Name:    "Foo Template",
			},
			args: args{
				template: template.Template{Website: "https://example.com", Name: "Foo Template"},
			},
			want: true,
		},
		{
			name: "no match - name over-defined",
			fields: fields{
				Website: "https://example.com",
				Name:    "Foo Template 2",
			},
			args: args{
				template: template.Template{Website: "https://example.com", Name: "Foo Template"},
			},
			want: false,
		},
		{
			name: "no match - subdomain mismatch",
			fields: fields{
				Website: "https://https://example.com",
				Name:    "Foo Template",
			},
			args: args{
				template: template.Template{Website: "https://example.com", Name: "Foo Template"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := template.Filter{
				Website: tt.fields.Website,
				Name:    tt.fields.Name,
			}

			result := f.Match(context.Background(), tt.args.template)

			assert.Equalf(t, tt.want, result, "Match(%v, %v)", context.Background(), tt.args.template)
		})
	}
}
