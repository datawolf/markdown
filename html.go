//
// html.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

//
//
// HTML rendering backend
//
//

package markdown

import (
	"bytes"
	"fmt"
)

// Html renderer configuration options.
const (
	HTML_SKIP_HTML = 1 << iota
	HTML_SKIP_STYLE
	HTML_SKIP_IMAGES
	HTML_SKIP_LINKS
	HTML_SAFELINK
	HTML_NOFOLLOW_LINKS
	HTML_NOREFERRER_LINKS
	HTML_HREF_TARGET_BLANK
	HTML_TOC
	HTML_OMIT_CONTENTS
	HTML_COMPLETE_PAGE
	HTML_USE_XHTML // generate XHTML output instead of HTML
	HTML_USE_SMARTYPANTS
	HTML_SMARTYPANTS_FRACTIONS
	HTML_SMARTYPANTS_DASHES
	HTML_SMARTYPANTS_LATEX_DASHES
	HTML_SMARTYPANTS_ANGLED_QUOTES
	HTML_FOOTNOTE_RETURN_LINKS
)

// HtmlRendererParameters defines the html renderer parameters
type HtmlRendererParameters struct {
	// prepend this text to each relative URL
	AbsolutePrefix string
	// Add this text to each footnote anchor, to ensure uniqueness
	FootnoteAnchorPrefix string
	// Show this text inside the <a> tag for a footnote return link, if
	// the HTML_FOOTNOTE_RETURN_LINKS flag is enabled. If blank, the string
	// <sup>[return]</sup> is used.
	FootnoteReturnLinkContents string
	// If set, add this text to the front of each Header ID, to ensure uniqueness
	HeaderIDPrefix string
	// If set, add this text to the back of each Header ID, to ensure uniqueness
	HeaderIDSuffix string
}

// Html is a type that implements the Renderer interface for HTML output
//
// Do not create this directly, instead use the HtmlRenderer function
type Html struct {
	flags    int    // HTML_* options
	closeTag string // The close tag: either " />" or ">"
	title    string // The document title
	css      string // Optional css file url

	parameters HtmlRendererParameters

	// table of contents data
	tocMarker    int
	headerCount  int
	currentLevel int
	toc          *bytes.Buffer

	// Track header Ids to prevent ID collision in a single generation
	headerIDs map[string]int
}

const (
	xhtmlClose = " />"
	htmlClose  = ">"
)

// HtmlRenderer creates and configures an Html object, which satisfies the Renderer interface.
//
// flags is a set of HTML_* options ORed together
// title is the title of the document
// css is a URL for the document's stylesheet.
//
// title and css are only used when HTML_COMPLETE_PAGE is selected.
func HtmlRenderer(flags int, title, css string) Renderer {
	return HtmlRendererWithParameters(flags, title, css, HtmlRendererParameters{})
}

func HtmlRendererWithParameters(flags int, title, css string, renderParameters HtmlRendererParameters) Renderer {
	// configure the rendering engine
	closeTag := htmlClose
	if flags&HTML_USE_XHTML != 0 {
		closeTag = xhtmlClose
	}

	if renderParameters.FootnoteReturnLinkContents == "" {
		renderParameters.FootnoteReturnLinkContents = `<sup>[return</sup>]`
	}

	return &Html{
		flags:        flags,
		closeTag:     closeTag,
		title:        title,
		css:          css,
		parameters:   renderParameters,
		headerCount:  0,
		currentLevel: 0,
		toc:          new(bytes.Buffer),
		headerIDs:    make(map[string]int),
	}
}

func (html *Html) GetFlags() int {
	return html.flags
}

func (html *Html) DocumentHeader(out *bytes.Buffer) {
	if html.flags&HTML_COMPLETE_PAGE == 0 {
		return
	}

	ending := ""

	if html.flags&HTML_USE_XHTML != 0 {
		out.WriteString("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" ")
		out.WriteString("\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
		out.WriteString("<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")
		ending = " /"
	} else {
		out.WriteString("<!DOCTYPE html>\n")
		out.WriteString("<html>\n")
	}

	out.WriteString("<head>\n")
	out.WriteString("  <title>")
	html.NormalText(out, []byte(html.title))
	out.WriteString("</title>\n")
	out.WriteString("  <meta name=\"GENERATOR\" content=\"Markdown Processor v")
	out.WriteString(VERSION)
	out.WriteString("\"")
	out.WriteString(">\n")
	out.WriteString("  <meta charset=\"utf-8\"")
	out.WriteString(ending)
	out.WriteString(">\n")

	if html.css != "" {
		out.WriteString("  <link rel=\"stylesheet\" type=\"text/css\" href=\"")
		attrEscape(out, []byte(html.css))
		out.WriteString(ending)
		out.WriteString(">\n")
	}
	out.WriteString("</head>\n")
	out.WriteString("<body>\n")

	html.tocMarker = out.Len()
}

func (html *Html) DocumentFooter(out *bytes.Buffer) {
	// finalize an d insert the table of contents
	if html.flags&HTML_TOC != 0 {
	}

	if html.flags&HTML_COMPLETE_PAGE != 0 {
		out.WriteString("\n</body>\n")
		out.WriteString("</html>\n")
	}
}

func (html *Html) Header(out *bytes.Buffer, header func() bool, level int, id string) {
	marker := out.Len()
	doubleSpace(out)

	if id == "" && html.flags&HTML_TOC != 0 {
		id = fmt.Sprintf("toc_%d", html.headerCount)
	}

	if id != "" {
		out.WriteString(fmt.Sprintf("<h%d id=\"%s\">", level, id))
	} else {
		out.WriteString(fmt.Sprintf("<h%d>", level))
	}

	//tocMarker := out.Len()
	if !header() {
		out.Truncate(marker)
		return
	}

	out.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func (html *Html) NormalText(out *bytes.Buffer, text []byte) {
	if html.flags&HTML_USE_SMARTYPANTS != 0 {
		html.Smartypants(out, text)
	} else {
		attrEscape(out, text)
	}
}

func (html *Html) Smartypants(out *bytes.Buffer, text []byte) {

}

func (html *Html) Emphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
}

func (html *Html) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
}

func (html *Html) TripleEmphasis(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")

}

func (html *Html) StrikeThrough(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}

	out.WriteString("<del>")
	out.Write(text)
	out.WriteString("</del>")
}

func (html *Html) CodeSpan(out *bytes.Buffer, text []byte) {
	if len(text) == 0 {
		return
	}
	out.WriteString("<code>")
	out.Write(text)
	out.WriteString("</code>")
}
func (html *Html) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	doubleSpace(out)
	out.WriteString("<p>")
	if !text() {
		out.Truncate(marker)
		return
	}
	out.WriteString("</p>\n")
}

func (html *Html) LineBreak(out *bytes.Buffer) {
	out.WriteString("<br")
	out.WriteString(html.closeTag)
	out.WriteByte('\n')
}
