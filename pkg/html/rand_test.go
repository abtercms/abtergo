package html_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/abtergo/abtergo/pkg/html"
)

func TestFixAssets(t *testing.T) {
	a := html.Meta{
		Name:     "foo",
		Property: "foo",
	}

	type args struct {
		assets html.Assets
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "only one of name or property attribute is set for all meta entries",
			args: args{
				assets: html.Assets{
					HeaderMeta: html.MetaList{
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
						a, a, a, a, a,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := html.FixAssets(tt.args.assets)

			for _, meta := range got.HeaderMeta {
				assert.NotEqual(t, meta.Name, meta.Property)
				assert.NotEmpty(t, meta.Name+meta.Property)
			}
		})
	}
}
