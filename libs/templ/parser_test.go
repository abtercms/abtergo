package templ_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/abtergo/abtergo/libs/templ"
)

func TestRenderer_ParseBlocks(t *testing.T) {
	type args struct {
		template string
	}
	tests := []struct {
		name string
		args args
		want []templ.ViewTag
	}{
		{
			name: "no attributes, no content",
			args: args{
				template: `Hello, <block/>!`,
			},
			want: []templ.ViewTag{
				{
					TagName: "block",
					Needles: []string{"<block/>"},
				},
			},
		},
		{
			name: "no attributes, but content",
			args: args{
				template: `Hello, <block   >foo</block>!`,
			},
			want: []templ.ViewTag{
				{
					TagName: "block",
					Needles: []string{`<block   >foo</block>`},
					Content: "foo",
				},
			},
		},
		{
			name: "simple",
			args: args{
				template: `Hello, <block module="moduleName" id="entityName" arg1="this" arg2="that"/>!`,
			},
			want: []templ.ViewTag{
				{
					TagName: "block",
					Needles: []string{`<block module="moduleName" id="entityName" arg1="this" arg2="that"/>`},
					Attributes: []html.Attribute{
						{Key: "module", Val: "moduleName"},
						{Key: "id", Val: "entityName"},
						{Key: "arg1", Val: "this"},
						{Key: "arg2", Val: "that"},
					},
				},
			},
		},
		{
			name: "retrieverDouble",
			args: args{
				template: `Hello, <block module="moduleName" id="entityName" arg1="this" arg2="that"/>! <block module="moduleName2" id="entityName1" arg1="that" arg2="this">foo</block>!`,
			},
			want: []templ.ViewTag{
				{
					TagName: "block",
					Needles: []string{`<block module="moduleName" id="entityName" arg1="this" arg2="that"/>`},
					Attributes: []html.Attribute{
						{Key: "module", Val: "moduleName"},
						{Key: "id", Val: "entityName"},
						{Key: "arg1", Val: "this"},
						{Key: "arg2", Val: "that"},
					},
				},
				{
					TagName: "block",
					Needles: []string{`<block module="moduleName2" id="entityName1" arg1="that" arg2="this">foo</block>`},
					Content: "foo",
					Attributes: []html.Attribute{
						{Key: "module", Val: "moduleName2"},
						{Key: "id", Val: "entityName1"},
						{Key: "arg1", Val: "that"},
						{Key: "arg2", Val: "this"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := templ.NewParser("block")
			got, err := sut.Parse(tt.args.template)

			assert.NoError(t, err)
			assert.Len(t, got, len(tt.want))
			assert.Equal(t, tt.want, got)
		})
	}
}
