package renderer_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/cbroglie/mustache"

	"github.com/abtergo/abtergo/pkg/renderer"
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
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
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
			want: `Hello, World!
You have just won $12!`,
			wantErr: false,
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
			want: `Hello, World!
You have just won $12!

Foobar`,
			wantErr: false,
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
			want: `Hello, World!
You have just won $12!

`,
			wantErr: false,
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
			want: `Hello, World!
You have just won $12!

Foobar`,
			wantErr: false,
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
			want: `Your list is:
(a)(b)(c)(d)(e)`,
			wantErr: false,
		},
		{
			name: "lamda",
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
			want:    fmt.Sprintf("Hello, World!\n%f", math.SqrtPi),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := renderer.NewRenderer()
			got, err := r.Render(tt.args.template, tt.args.context...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Render() got = %v, want %v", got, tt.want)
			}
		})

		t.Run(tt.name+" - with AddContext", func(t *testing.T) {
			r := renderer.NewRenderer()
			r.AddContext(tt.args.context...)
			got, err := r.Render(tt.args.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Render() got = %v, want %v", got, tt.want)
			}
		})
	}
}
