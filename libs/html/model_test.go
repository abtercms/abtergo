package html_test

import (
	"reflect"
	"testing"

	"github.com/abtergo/abtergo/libs/html"
)

func TestAttributes_Clone(t *testing.T) {
	tests := []struct {
		name string
		a    html.Attributes
		want html.Attributes
	}{
		{
			name: "empty",
			a:    nil,
			want: nil,
		},
		{
			name: "default",
			a: html.Attributes{
				"foo": "Foo",
				"bar": "Bar",
			},
			want: html.Attributes{
				"foo": "Foo",
				"bar": "Bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLink_Clone(t *testing.T) {
	type fields struct {
		Rel        string
		Href       string
		Attributes html.Attributes
	}
	tests := []struct {
		name   string
		fields fields
		want   html.Link
	}{
		{
			name: "default",
			fields: fields{
				Rel:  "foo",
				Href: "foo.css",
				Attributes: html.Attributes{
					"foo": "Foo",
					"bar": "Bar",
				},
			},
			want: html.Link{
				Rel:  "foo",
				Href: "foo.css",
				Attributes: html.Attributes{
					"foo": "Foo",
					"bar": "Bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := html.Link{
				Rel:        tt.fields.Rel,
				Href:       tt.fields.Href,
				Attributes: tt.fields.Attributes,
			}
			if got := l.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinks_Clone(t *testing.T) {
	tests := []struct {
		name string
		l    html.Links
		want html.Links
	}{
		{
			name: "empty",
			l:    nil,
			want: nil,
		},
		{
			name: "default",
			l: html.Links{
				{
					Rel:  "foo",
					Href: "foo.css",
					Attributes: html.Attributes{
						"foo": "Foo",
						"bar": "Bar",
					},
				},
			},
			want: html.Links{
				{
					Rel:  "foo",
					Href: "foo.css",
					Attributes: html.Attributes{
						"foo": "Foo",
						"bar": "Bar",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetaList_Clone(t *testing.T) {
	tests := []struct {
		name string
		ml   html.MetaList
		want html.MetaList
	}{
		{
			name: "empty",
			ml:   nil,
			want: nil,
		},
		{
			name: "default",
			ml: html.MetaList{
				{
					Name:    "foo",
					Content: "Foo",
					Attributes: html.Attributes{
						"bar": "Bar",
						"baz": "Baz",
					},
				},
				{
					Content: "bar",
					Attributes: html.Attributes{
						"baz":  "Baz",
						"quix": "Quix",
					},
				},
			},
			want: html.MetaList{
				{
					Name:    "foo",
					Content: "Foo",
					Attributes: html.Attributes{
						"bar": "Bar",
						"baz": "Baz",
					},
				},
				{
					Content: "bar",
					Attributes: html.Attributes{
						"baz":  "Baz",
						"quix": "Quix",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ml.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeta_Clone(t *testing.T) {
	type fields struct {
		Name       string
		Property   string
		Content    string
		Attributes html.Attributes
	}
	tests := []struct {
		name   string
		fields fields
		want   html.Meta
	}{
		{
			name: "default",
			fields: fields{
				Name:    "foo",
				Content: "Foo",
				Attributes: html.Attributes{
					"foo": "Foo",
					"bar": "Bar",
				},
			},
			want: html.Meta{
				Name:    "foo",
				Content: "Foo",
				Attributes: html.Attributes{
					"foo": "Foo",
					"bar": "Bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := html.Meta{
				Name:       tt.fields.Name,
				Content:    tt.fields.Content,
				Attributes: tt.fields.Attributes,
			}
			if got := m.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScript_Clone(t *testing.T) {
	type fields struct {
		Src        string
		Attributes html.Attributes
	}
	tests := []struct {
		name   string
		fields fields
		want   html.Script
	}{
		{
			name: "default",
			fields: fields{
				Src: "foo.js",
				Attributes: html.Attributes{
					"foo": "Foo",
					"bar": "Bar",
				},
			},
			want: html.Script{
				Src: "foo.js",
				Attributes: html.Attributes{
					"foo": "Foo",
					"bar": "Bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := html.Script{
				Src:        tt.fields.Src,
				Attributes: tt.fields.Attributes,
			}
			if got := s.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScripts_Clone(t *testing.T) {
	tests := []struct {
		name string
		s    html.Scripts
		want html.Scripts
	}{
		{
			name: "empty",
			s:    nil,
			want: nil,
		},
		{
			name: "default",
			s: html.Scripts{
				{
					Src: "foo.js",
					Attributes: html.Attributes{
						"baz":  "Baz",
						"quix": "Quix",
					},
				},
			},
			want: html.Scripts{
				{
					Src: "foo.js",
					Attributes: html.Attributes{
						"baz":  "Baz",
						"quix": "Quix",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssets_Clone(t *testing.T) {
	type fields struct {
		HeaderCSS  html.Links
		HeaderJs   html.Scripts
		HeaderMeta html.MetaList
		FooterJs   html.Scripts
	}
	tests := []struct {
		name   string
		fields fields
		want   html.Assets
	}{
		{
			name: "default",
			fields: fields{
				HeaderCSS:  html.Links{},
				HeaderJs:   html.Scripts{},
				HeaderMeta: html.MetaList{},
				FooterJs:   html.Scripts{},
			},
			want: html.Assets{
				HeaderCSS:  html.Links{},
				HeaderJS:   html.Scripts{},
				HeaderMeta: html.MetaList{},
				FooterJS:   html.Scripts{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := html.Assets{
				HeaderCSS:  tt.fields.HeaderCSS,
				HeaderJS:   tt.fields.HeaderJs,
				HeaderMeta: tt.fields.HeaderMeta,
				FooterJS:   tt.fields.FooterJs,
			}
			if got := a.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}
