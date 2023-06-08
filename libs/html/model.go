package html

type Attributes map[string]string

func (a Attributes) Clone() Attributes {
	if a == nil {
		return nil
	}

	clone := make(Attributes, len(a))
	for k, v := range a {
		clone[k] = v
	}
	return clone
}

type Link struct {
	Rel        string
	Href       string
	Attributes Attributes
}

func (l Link) Clone() Link {
	return Link{
		Rel:        l.Rel,
		Href:       l.Href,
		Attributes: l.Attributes.Clone(),
	}
}

type Links []Link

func (l Links) Clone() Links {
	if l == nil {
		return nil
	}

	clone := make(Links, len(l))
	for i, v := range l {
		clone[i] = v.Clone()
	}
	return clone
}

type Meta struct {
	Name       string
	Content    string
	Property   string
	Attributes Attributes
}

func (l Meta) Clone() Meta {
	return Meta{
		Name:       l.Name,
		Content:    l.Content,
		Property:   l.Property,
		Attributes: l.Attributes.Clone(),
	}
}

type MetaList []Meta

func (m MetaList) Clone() MetaList {
	if m == nil {
		return nil
	}

	clone := make(MetaList, len(m))
	for i, v := range m {
		clone[i] = v.Clone()
	}
	return clone
}

type Script struct {
	Src        string
	Attributes Attributes
}

func (s Script) Clone() Script {
	return Script{
		Src:        s.Src,
		Attributes: s.Attributes.Clone(),
	}
}

type Scripts []Script

func (s Scripts) Clone() Scripts {
	if s == nil {
		return nil
	}

	clone := make(Scripts, len(s))
	for i, v := range s {
		clone[i] = v.Clone()
	}
	return clone
}

type Assets struct {
	HeaderCSS  Links
	HeaderJS   Scripts
	HeaderMeta MetaList
	FooterJS   Scripts
}

func (a Assets) Clone() Assets {
	clone := Assets{
		HeaderCSS:  a.HeaderCSS.Clone(),
		HeaderJS:   a.HeaderJS.Clone(),
		HeaderMeta: a.HeaderMeta.Clone(),
		FooterJS:   a.FooterJS.Clone(),
	}
	return clone
}
