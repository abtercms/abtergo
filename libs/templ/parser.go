package templ

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"

	"github.com/abtergo/abtergo/libs/model"
)

type ViewTag struct {
	TagName    string
	Needles    []string
	Content    string
	Attributes []html.Attribute
}

type Parser interface {
	Parse(template string) ([]ViewTag, error)
}

func NewParser(tagNames ...string) Parser {
	return &parser{
		tagNames: tagNames,
	}
}

type parser struct {
	tagNames []string
}

func (p *parser) Parse(template string) ([]ViewTag, error) {
	var result []ViewTag
	for _, tagName := range p.tagNames {
		matches := p.findBlockTags(template, tagName)
		viewTags, err := p.collectViewTags(tagName, matches)
		if err != nil {
			return nil, errors.Wrap(err, "failed to collect view tags")
		}

		result = append(result, viewTags...)
	}

	return result, nil
}

func (p *parser) findBlockTags(template, tagName string) [][]string {
	regex := regexp.MustCompile(fmt.Sprintf(`<%s\s*([^/>]*?)>(.*?)</%s\s*>`, tagName, tagName))
	tagsWithContent := regex.FindAllStringSubmatch(template, -1)

	for _, match := range tagsWithContent {
		template = strings.Replace(template, match[0], "", -1)
	}

	regex = regexp.MustCompile(fmt.Sprintf(`<%s\s*([^/>]*?)/>()`, tagName))
	selfClosedTags := regex.FindAllStringSubmatch(template, -1)

	return append(selfClosedTags, tagsWithContent...)
}

func (p *parser) collectViewTags(tagName string, matches [][]string) ([]ViewTag, error) {
	needleMap := make(map[model.ETag][]string)
	viewTagMap := make(map[model.ETag]ViewTag)
	for _, match := range matches {
		attributes, err := p.parseAttributes(match[1], tagName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse attributes")
		}

		vt := ViewTag{
			TagName:    tagName,
			Content:    match[2],
			Attributes: attributes,
		}
		eTag := model.ETagFromAny(vt)

		if _, ok := needleMap[eTag]; !ok {
			needleMap[eTag] = make([]string, 0)
			viewTagMap[eTag] = vt
		}

		needleMap[eTag] = append(needleMap[eTag], match[0])
	}

	result := make([]ViewTag, 0, len(viewTagMap))
	for eTag, needles := range needleMap {
		vt := viewTagMap[eTag]
		vt.Needles = needles
		result = append(result, vt)
	}

	return result, nil
}

func (p *parser) parseAttributes(rawAttributes, tagName string) ([]html.Attribute, error) {
	div := fmt.Sprintf(`<%s %s />`, tagName, rawAttributes)
	nodes, err := html.ParseFragment(strings.NewReader(div), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse attributes")
	}
	if len(nodes) == 0 {
		return nil, nil
	}

	var f func(*html.Node) []html.Attribute
	f = func(n *html.Node) []html.Attribute {
		if n.Type == html.ElementNode && n.Data == tagName {
			return n.Attr
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			attr := f(c)
			if len(attr) > 0 {
				return attr
			}
		}

		return nil
	}

	return f(nodes[0]), nil
}
