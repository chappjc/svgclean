// Copyright (c) 2018, The Decred developers
// See LICENSE for details.

package svgclean

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

const svgHeader = `<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN" "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
`

type elements []*Element

// Element is the representation of an xml element/tag. All attributes, except
// for href, are captured in the Attrs. The AttrHref field captures href
// attributes separately for ease of removal. Comments and chardata are
// retained, but not CDATA. All child elements, except html <a> and <script>
// tags, are captured in Elements. Link and script tags are captured separately
// for ease of removal.
type Element struct {
	XMLName  xml.Name
	Attrs    []*xml.Attr `xml:",any,attr"`
	AttrHref *xml.Attr   `xml:"href,attr,omitempty"`
	Comment  string      `xml:",comment"`
	CharData string      `xml:",chardata"`
	Links    []*Element  `xml:"a"`
	Scripts  []*Element  `xml:"script"`

	// Elements is a pointer to a slice of child elements. It is a pointer to
	// avoid marshaling an empty <Elements></Elements> for a slice with no
	// non-nil pointer..
	Elements *elements `xml:",any"`
}

func CleanSVGString(svg string) string {
	buf := bytes.NewBuffer([]byte(svg))
	dec := xml.NewDecoder(buf)

	var e Element
	err := dec.Decode(&e)
	if err != nil {
		fmt.Printf("failed to unmarshal source svg: %v", err)
		os.Exit(1)
	}

	//e := (*Element)(&x)
	CleanElementTree(&e)

	b, err := xml.MarshalIndent(e, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal cleaned svg: %v", err)
		os.Exit(1)
	}

	return svgHeader + string(b)
}

// CleanElementTree recursively processes the given element and all of its
// children via CleanElements.
func CleanElementTree(e *Element) {
	CleanElements(elements{e})
}

// CleanElements recursively processes the given elements and all of their
// children. The following are removed: <a> and <script> tags, href attributes,
// and both leading/trailing white space and CR/LF characters from character
// data outside of <p>, <span>, and <div> tags.
func CleanElements(Elements elements) {
	okFun := func(e *Element) bool {
		if e == nil {
			return false
		}

		// Eliminate all <a> elements and their children.
		if e.Links != nil {
			fmt.Printf("Removing %d <a> tag(s).\n", len(e.Links))
			e.Links = nil
		}

		// Eliminate all <script> elements and their children.
		if e.Scripts != nil {
			fmt.Printf("Removing %d <script> tag(s).\n", len(e.Scripts))
			e.Scripts = nil
		}

		// Eliminate href attributes, in case they are in elements other than
		// <a>, which is not correct, but who knows.
		if e.AttrHref != nil {
			fmt.Printf(`Removing stray href attribute "%s".\n`, e.AttrHref.Value)
			e.AttrHref = nil
		}

		for i, a := range e.Attrs {
			if strings.HasPrefix(a.Name.Local, "on") {
				fmt.Printf("Removing on* attribute \"%s\".\n", a.Name.Local)
				e.Attrs[i] = nil
			}
			if strings.Contains(a.Value, "javascript") {
				fmt.Printf("Removing attribute \"%s\" with possible javascript.\n",
					a.Name.Local)
				e.Attrs[i] = nil
			}
		}

		// Remove newlines and trailing/leading spaces from character data.
		switch e.XMLName.Local {
		case "p", "span", "div":
		default:
			r := strings.NewReplacer("\n", "", "\r", "")
			e.CharData = strings.TrimSpace(r.Replace(e.CharData))
		}

		// Remove redundant xmlns attribute, which will be retained by Attrs.
		e.XMLName.Space = ""

		// Only OK to process children if slice has elements.
		return e.Elements != nil && len(*e.Elements) > 0
	}

	WalkElements(Elements, okFun)
}

// WalkElements recursively processes the *Element slice with the provided
// function. Processing is recursive in that the child elements are also
// processed with WalkElements, but only if the function returns true.
func WalkElements(Elements elements, okFun func(*Element) bool) {
	for _, e := range Elements {
		if okFun(e) {
			WalkElements(*e.Elements, okFun)
		}
	}
}

// type XEl Element

// UnmarshalXML implements xml.Unmarshaler for XEl.
// func (n *XEl) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	//n.Attrs = start.Attr
// 	err := d.DecodeElement((*Element)(n), &start)
// 	if err != nil {
// 		return err
// 	}
// 	if n.Scripts != nil {
// 		spew.Dump(n.Scripts)
// 	}
// 	return nil
// }
