package html_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	html2 "github.com/abtergo/abtergo/libs/html"
)

func TestFixAssets(t *testing.T) {
	a := html2.Meta{
		Name:     "foo",
		Property: "foo",
	}

	type args struct {
		assets html2.Assets
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "only one of name or property attribute is set for all meta entries",
			args: args{
				assets: html2.Assets{
					HeaderMeta: html2.MetaList{
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
			got := html2.FixAssets(tt.args.assets)

			for _, meta := range got.HeaderMeta {
				assert.NotEqual(t, meta.Name, meta.Property)
				assert.NotEmpty(t, meta.Name+meta.Property)
			}
		})
	}
}
