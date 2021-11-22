package core

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io"
	"strings"
)

func (fns *InternalFunctionSet) HtmlParser(raw interface{}) Value {
	var docReader io.Reader
	if bs, ok := raw.([]byte); ok {
		docReader = bytes.NewReader(bs)
	} else if str, ok := raw.(string); ok {
		docReader = strings.NewReader(str)
	} else {
		panic("invalid input")
	}

	reader, err := goquery.NewDocumentFromReader(docReader)
	if err != nil {
		return NULL
	}
	obj := &HtmlDocument{reader}
	return newClass("HtmlDocument", &obj)
}

type HtmlDocument struct {
	doc *goquery.Document
}

func (d *HtmlDocument) Find(seletor string) Value {
	return newHtmlSelection(d.doc.Find(seletor))
}

type HtmlSelection struct {
	raw *goquery.Selection
}

func newHtmlSelection(raw *goquery.Selection) Value {
	obj := &HtmlSelection{raw}
	return newClass("HtmlSelection", &obj)
}

func (s *HtmlSelection) Find(seletor string) Value {
	return newHtmlSelection(s.raw.Find(seletor))
}
func (s *HtmlSelection) Not(seletor string) Value {
	return newHtmlSelection(s.raw.Not(seletor))
}

func (s *HtmlSelection) Each(fn Function) {
	s.raw.Each(func(i int, s *goquery.Selection) {
		var args = []Value{
			newQKValue(i),
			newHtmlSelection(s),
		}
		fn.setArgs(args)

		fn.execute()
	})
}

func (s *HtmlSelection) Html() string {
	html, err := s.raw.Html()
	if err != nil {
		return ""
	}
	return html
}
func (s *HtmlSelection) Text() string {
	return s.raw.Text()
}

func (s *HtmlSelection) Attr(name string) string {
	attr, exists := s.raw.Attr(name)
	if !exists {
		return ""
	}
	return attr
}
