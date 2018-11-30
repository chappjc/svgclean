// Copyright (c) 2018, The Decred developers
// See LICENSE for details.

package svgclean

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/chappjc/svgclean/whitelist"
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
	// non-nil pointer.
	Elements *elements `xml:",any"`
}

type elementsPromisc []*ElementPromisc

// ElementPromisc is like Element, except that it does not capture anything by
// name; all tags and attributes to to the Elements and Attrs slices. This
// facilitates the whitelist approach.
type ElementPromisc struct {
	XMLName  xml.Name
	Attrs    []*xml.Attr `xml:",any,attr"`
	Comment  string      `xml:",comment"`
	CharData string      `xml:",chardata"`
	//CDATA string `xml:",cdata"`

	// Elements is a pointer to a slice of child elements. It is a pointer to
	// avoid marshaling an empty <Elements></Elements> for a slice with no
	// non-nil pointer.
	Elements *elementsPromisc `xml:",any"`
}

func CleanSVGStringBlack(svg string) string {
	return cleanSVGString(svg, true)
}

func CleanSVGStringWhite(svg string) string {
	return cleanSVGString(svg, false)
}

func cleanSVGString(svg string, blacklist bool) string {
	buf := bytes.NewBuffer([]byte(svg))
	dec := xml.NewDecoder(buf)

	var b []byte

	//e := (*Element)(&x)
	if blacklist {
		var e Element
		err := dec.Decode(&e)
		if err != nil {
			fmt.Printf("failed to unmarshal source svg: %v", err)
			os.Exit(1)
		}
		CleanElementTreeBlack(&e)
		b, err = xml.MarshalIndent(e, "", "  ")
		if err != nil {
			fmt.Printf("failed to marshal cleaned svg: %v", err)
			os.Exit(1)
		}
	} else {
		var e ElementPromisc
		err := dec.Decode(&e)
		if err != nil {
			fmt.Printf("failed to unmarshal source svg: %v", err)
			os.Exit(1)
		}
		CleanElementTreeWhite(&e)
		b, err = xml.MarshalIndent(e, "", "  ")
		if err != nil {
			fmt.Printf("failed to marshal cleaned svg: %v", err)
			os.Exit(1)
		}
	}

	return svgHeader + string(b)
}

// CleanElementTreeBlack recursively processes the given element and all of its
// children via CleanElementsBlackList.
func CleanElementTreeBlack(e *Element) {
	CleanElementsBlackList(elements{e})
}

// CleanElementTreeWhite recursively processes the given element and all of its
// children via CleanElementsWhiteList.
func CleanElementTreeWhite(e *ElementPromisc) {
	CleanElementsWhiteList(elementsPromisc{e})
}

func CleanElementsWhiteList(Elements elementsPromisc) {
	okFun := func(e *ElementPromisc) bool {
		if e == nil {
			return false
		}

		// Check tag
		if !whitelist.IsAllowedTag(e.XMLName.Local) {
			fmt.Printf("Removing disallowed tag \"%s\".\n", e.XMLName.Local)
			return false // FilterElementsWhite will set this element to nil
		}

		// Check all attributes
		var okAttrs []*xml.Attr
		for _, a := range e.Attrs {
			// Allow only whitelisted attributes.
			if !whitelist.IsAllowedAttribute(a.Name.Local) {
				fmt.Printf("Removing disallowed attribute \"%s\".\n",
					a.Name.Local)
				continue
			}

			// Filter values too.
			if strings.Contains(a.Value, "javascript") {
				fmt.Printf("Removing attribute \"%s\" with possible javascript.\n",
					a.Name.Local)
				continue
			}
			okAttrs = append(okAttrs, a)
		}
		e.Attrs = okAttrs

		// Remove newlines and trailing/leading spaces from character data.
		switch e.XMLName.Local {
		case "p", "span", "div":
		default:
			r := strings.NewReplacer("\n", "", "\r", "")
			e.CharData = strings.TrimSpace(r.Replace(e.CharData))
		}

		// Remove redundant xmlns attribute, which will be retained by Attrs.
		e.XMLName.Space = ""

		// signal to keep this cleaned element
		return true
	}

	FilterElementsWhite(Elements, okFun)
}

// CleanElementsBlackListed recursively processes the given elements and all of
// their children. The following are removed: <a> and <script> tags, href
// attributes, and both leading/trailing white space and CR/LF characters from
// character data outside of <p>, <span>, and <div> tags.
func CleanElementsBlackList(Elements elements) {
	okFun := func(e *Element) {
		if e == nil {
			return
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
	}

	FilterElementsBlack(Elements, okFun)
}

// FilterElementsBlack recursively processes the *Element slice with the
// provided function. Processing is recursive in that the child elements are
// also processed with FilterElementsBlack, but only if the function returns true.
func FilterElementsBlack(Elements elements, cleanFun func(*Element)) {
	for _, e := range Elements {
		cleanFun(e)
		if e != nil && e.Elements != nil {
			FilterElementsBlack(*e.Elements, cleanFun)
		}
	}
}

func FilterElementsWhite(Elements elementsPromisc, okFun func(*ElementPromisc) bool) {
	for i, e := range Elements {
		if okFun(e) {
			if e.Elements != nil {
				FilterElementsWhite(*e.Elements, okFun)
			}
		} else {
			// clean up
			Elements[i] = nil
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
