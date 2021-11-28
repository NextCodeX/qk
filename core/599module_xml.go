package core

import (
	"fmt"
	"github.com/beevik/etree"
)

func (fns *InternalFunctionSet) Xml() Value {
	doc := etree.NewDocument()
	obj := &XmlDocument{}
	obj.raw = doc
	res := newClass("XmlDocument", &obj)
	obj.instance = res
	return res
}

type XmlDocument struct {
	instance Value
	raw      *etree.Document
}

func (xd *XmlDocument) Load(content interface{}) Value {
	switch v := content.(type) {
	case []byte:
		err := xd.raw.ReadFromBytes(v)
		assert(err != nil, "XmlDocument.load(content) failed: ", err)
	case string:
		err := xd.raw.ReadFromString(v)
		assert(err != nil, "XmlDocument.load(content) failed: ", err)
	default:
		panic("XmlDocument.load(content): content must be String/ByteArray")
	}
	return xd.instance
}

func (xd *XmlDocument) LoadFile(path string) {
	err := xd.raw.ReadFromFile(path)
	assert(err != nil, "XmlDocument.load(content) failed: ", err)
}
func (xd *XmlDocument) Root() Value {
	return newXmlElement(xd.raw.Root())
}

func (xd *XmlDocument) DefaultHead() {
	xd.raw.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
}
func (xd *XmlDocument) DefHead() {
	xd.DefaultHead()
}
func (xd *XmlDocument) SetHead(target, inst string) {
	xd.raw.CreateProcInst(target, inst)
}
func (xd *XmlDocument) NewElem(tag string) Value {
	e := xd.raw.CreateElement(tag)
	return newXmlElement(e)
}
func (xd *XmlDocument) NewAttr(key, value string) {
	xd.raw.CreateAttr(key, value)
}
func (xd *XmlDocument) NewComment(comment string) {
	xd.raw.CreateComment(comment)
}
func (xd *XmlDocument) NewCData(data string) {
	xd.raw.CreateCData(data)
}
func (xd *XmlDocument) SetCData(data string) {
	xd.raw.SetCData(data)
}
func (xd *XmlDocument) NewText(text string) {
	xd.raw.CreateText(text)
}
func (xd *XmlDocument) SetText(text string) {
	xd.raw.SetText(text)
}

func (xd *XmlDocument) Indent(space int) {
	xd.raw.Indent(space)
}

func (xd *XmlDocument) NoIndent() {
	xd.raw.Indent(etree.NoIndent)
}

func (xd *XmlDocument) Name() string {
	return xd.raw.Tag
}
func (xd *XmlDocument) Finds(xpath string) Value {
	elements := xd.raw.FindElements(xpath)
	var res []Value
	for _, e := range elements {
		res = append(res, newXmlElement(e))
	}
	return array(res)
}
func (xd *XmlDocument) Find(xpath string) Value {
	element := xd.raw.FindElement(xpath)
	if element == nil {
		return nil
	}
	return newXmlElement(element)
}

func (xd *XmlDocument) Str() string {
	res, err := xd.raw.WriteToString()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return res
}

func (xd *XmlDocument) Bytes() []byte {
	res, err := xd.raw.WriteToBytes()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (xd *XmlDocument) Save(path string) {
	err := xd.raw.WriteToFile(path)
	if err != nil {
		fmt.Println("XmlDocument.save() failed: ", err)
	}
}

func (xd *XmlDocument) Elems(tag string) Value {
	elements := xd.raw.SelectElements(tag)
	var res []Value
	for _, e := range elements {
		res = append(res, newXmlElement(e))
	}
	return array(res)
}

func (xd *XmlDocument) Elem(tag string) Value {
	element := xd.raw.SelectElement(tag)
	return newXmlElement(element)
}

func (xd *XmlDocument) Attr(key string) string {
	return xd.raw.SelectAttrValue(key, "")
}

func (xd *XmlDocument) Text() string {
	return xd.raw.Text()
}

type XmlElement struct {
	raw *etree.Element
}

func newXmlElement(raw *etree.Element) Value {
	obj := &XmlElement{raw}
	return newClass("XmlElement", &obj)
}

func (xe *XmlElement) Name() string {
	return xe.raw.Tag
}
func (xe *XmlElement) Finds(xpath string) Value {
	elements := xe.raw.FindElements(xpath)
	var res []Value
	for _, e := range elements {
		res = append(res, newXmlElement(e))
	}
	return array(res)
}
func (xe *XmlElement) Find(xpath string) Value {
	element := xe.raw.FindElement(xpath)
	if element == nil {
		return nil
	}
	return newXmlElement(element)
}

func (xe *XmlElement) Elems(tag string) Value {
	elements := xe.raw.SelectElements(tag)
	var res []Value
	for _, e := range elements {
		res = append(res, newXmlElement(e))
	}
	return array(res)
}
func (xe *XmlElement) Elem(tag string) Value {
	element := xe.raw.SelectElement(tag)
	return newXmlElement(element)
}
func (xe *XmlElement) Attr(key string) string {
	return xe.raw.SelectAttrValue(key, "")
}
func (xe *XmlElement) Text() string {
	return xe.raw.Text()
}

func (xe *XmlElement) NewElem(tag string) Value {
	e := xe.raw.CreateElement(tag)
	return newXmlElement(e)
}
func (xe *XmlElement) NewAttr(key, value string) {
	xe.raw.CreateAttr(key, value)
}
func (xe *XmlElement) NewComment(comment string) {
	xe.raw.CreateComment(comment)
}
func (xe *XmlElement) NewCData(data string) {
	xe.raw.CreateCData(data)
}
func (xe *XmlElement) SetCData(data string) {
	xe.raw.SetCData(data)
}
func (xe *XmlElement) NewText(text string) {
	xe.raw.CreateText(text)
}
func (xe *XmlElement) SetText(text string) {
	xe.raw.SetText(text)
}
