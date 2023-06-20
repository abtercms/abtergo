package templ_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/adelowo/onecache/memory"
	"github.com/cbroglie/mustache"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/abtergo/abtergo/libs/templ"
	mocks "github.com/abtergo/abtergo/mocks/libs/templ"
)

func TestRenderer_Render(t *testing.T) {
	type data struct {
		A bool
		B string
	}

	type args struct {
		template string
		context  []interface{}
	}
	type fields struct {
		parsedTemplate string
		viewTags       []templ.ViewTag
		templates      []string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   string
	}{
		{
			name: "simple",
			args: args{
				template: `Hello, {{name}}!
You have just won ${{value}}!`,
				context: []interface{}{
					map[string]string{"name": "World"},
					map[string]int{"value": 12},
				},
			},
			fields: fields{
				parsedTemplate: `Hello, World!`,
			},
			want: `Hello, World!
You have just won $12!`,
		},
		{
			name: "section",
			args: args{
				template: `Hello, {{name}}!
You have just won ${{value}}!

{{#A}}{{B}}{{/A}}`,
				context: []interface{}{
					map[string]string{"name": "World", "B": "Foobar"},
					map[string]int{"value": 12},
					map[string]bool{"A": true},
				},
			},
			fields: fields{
				parsedTemplate: `Hello, World!
You have just won $12!

Foobar`,
			},
			want: `Hello, World!
You have just won $12!

Foobar`,
		},
		{
			name: "section skipped",
			args: args{
				template: `Hello, {{name}}!
You have just won ${{value}}!

{{#A}}{{B}}{{/A}}`,
				context: []interface{}{
					map[string]string{"name": "World", "B": "Foobar"},
					map[string]int{"value": 12},
				},
			},
			fields: fields{
				parsedTemplate: `Hello, World!
You have just won $12!

`,
			},
			want: `Hello, World!
You have just won $12!

`,
		},
		{
			name: "section with data struct",
			args: args{
				template: `Hello, {{name}}!
You have just won ${{value}}!

{{#A}}{{B}}{{/A}}`,
				context: []interface{}{
					map[string]string{"name": "World"},
					map[string]int{"value": 12},
					data{A: true, B: "Foobar"},
				},
			},
			fields: fields{
				parsedTemplate: `Hello, World!
You have just won $12!

Foobar`,
			},
			want: `Hello, World!
You have just won $12!

Foobar`,
		},
		{
			name: "list",
			args: args{
				template: `Your list is:
{{#list}}({{.}}){{/list}}`,
				context: []interface{}{
					map[string]string{"name": "World"},
					map[string]int{"value": 12},
					map[string]interface{}{"list": []string{"a", "b", "c", "d", "e"}},
				},
			},
			fields: fields{
				parsedTemplate: `Your list is:
(a)(b)(c)(d)(e)`,
			},
			want: `Your list is:
(a)(b)(c)(d)(e)`,
		},
		{
			name: "lambda",
			args: args{
				template: `{{#Lambda}}Hello, {{{Name}}}!{{/Lambda}}`,
				context: []interface{}{
					struct {
						Name   string
						Lambda mustache.LambdaFunc
					}{
						Name: "World",
						Lambda: func(text string, render mustache.RenderFunc) (string, error) {
							return render(fmt.Sprintf("%s\n%f", text, math.SqrtPi))
						},
					},
				},
			},
			fields: fields{
				parsedTemplate: fmt.Sprintf("Hello, World!\n%f", math.SqrtPi),
			},
			want: fmt.Sprintf("Hello, World!\n%f", math.SqrtPi),
		},
		{
			name: "view tags",
			args: args{
				template: `{{#Lambda}}Hello, <block module="foo" name="AbterCMS" />!{{/Lambda}}`,
				context: []interface{}{
					struct {
						Name   string
						Lambda mustache.LambdaFunc
					}{
						Name: "World",
						Lambda: func(text string, render mustache.RenderFunc) (string, error) {
							return render(text)
						},
					},
				},
			},
			fields: fields{
				parsedTemplate: `Hello, <block module="foo" name="AbterCMS" />!`,
				viewTags: []templ.ViewTag{
					{
						TagName: "block",
						Needles: []string{`<block module="foo" name="AbterCMS" />`},
						Attributes: []html.Attribute{
							{
								Key: "module",
								Val: "foo",
							},
							{
								Key: "name",
								Val: "AbterCMS",
							},
						},
					},
				},
				templates: []string{"World"},
			},
			want: "Hello, World!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeRetrieverMap := make(map[string]templ.Retriever)
			fr := newRetrieverDouble()
			for i, vt := range tt.fields.viewTags {
				fr.SetViewTag(vt, tt.fields.templates[i])
				fakeRetrieverMap[vt.TagName] = fr
			}

			cacheDouble := memory.New()

			parserMock := new(mocks.Parser)
			parserMock.EXPECT().
				Parse(tt.fields.parsedTemplate).
				Maybe().
				Return(tt.fields.viewTags, nil)
			parserMock.EXPECT().
				Parse(tt.want).
				Maybe().
				Return(nil, nil)
			defer parserMock.AssertExpectations(t)
			r := templ.NewRenderer(parserMock, fakeRetrieverMap, cacheDouble)

			got, err := r.Render(tt.args.template, tt.args.context...)

			assert.NoError(t, err)
			if got != tt.want {
				t.Errorf("Render() got = %v, want %v", got, tt.want)
			}
		})
	}
}
