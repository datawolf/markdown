//
// markdown.go
// Copyright (C) 2016 wanglong <wanglong@laoqinren.net>
//
// Distributed under terms of the MIT license.
//

//
//
// Markdown parsing and processing
//
//

package markdown

import (
	"bytes"
)

const VERSION = "0.1"

// These are the suupported markdown parsing extensions.
// OR these values together to select multiple extensions
const (
	EXTENSION_NO_INTRA_EMPHASIS = 1 << iota // ignore emphasis marker inside words
	EXTENSION_TABLES
	EXTENSION_FENCED_CODE
	EXTENSION_AUTOLINK
	EXTENSION_STRIKETHROUGH
	EXTENSION_LAX_HTML_BLOCKS
	EXTENSION_SPACE_HEADERS // be strict about prefix header rules
	EXTENSION_HARD_LINE_BREAK
	EXTENSION_TAB_SIZE_EIGHT
	EXTENSION_FOOTNOTES
	EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	EXTENSION_HEADER_IDS
	EXTENSION_TITLEBLOCK
	EXTENSION_AUTO_HEADER_IDS
	EXTENSION_BACKSLASH_LINE_BREAK
	EXTENSION_DEFINITION_LISTS

	commonHtmlFlags = 0 |
		HTML_USE_XHTML |
		HTML_USE_SMARTYPANTS |
		HTML_SMARTYPANTS_FRACTIONS |
		HTML_SMARTYPANTS_DASHES |
		HTML_SMARTYPANTS_LATEX_DASHES

	commonExtension = 0 |
		EXTENSION_NO_INTRA_EMPHASIS |
		EXTENSION_TABLES |
		EXTENSION_FENCED_CODE |
		EXTENSION_AUTOLINK |
		EXTENSION_STRIKETHROUGH |
		EXTENSION_SPACE_HEADERS |
		EXTENSION_HEADER_IDS |
		EXTENSION_BACKSLASH_LINE_BREAK |
		EXTENSION_DEFINITION_LISTS
)

// The size of a tab stop.
const (
	TAB_SIZE_DEFAULT = 4
	TAB_SIZE_EIGHT   = 8
)

// Renderer is the rendering interface
// This is the mostly interest if you are implementing a new renderering format.
//
// When a byte slice is provided, it contains the contents fo the element.
//
// Currently Html implementation is provided
type Renderer interface {
	// block-level callbacks
	//	BlockCode(out *bytes.Buffer, text []byte, lang string)
	//	BlockQuote(out *bytes.Buffer, text []byte)
	Header(out *bytes.Buffer, text func() bool, level int, id string)

	// span-level callbacks
	//	CodeSpan(out *bytes.Buffer, text []byte)

	// Low-level callbacks
	NormalText(out *bytes.Buffer, entity []byte)

	// Header and footer
	DocumentHeader(out *bytes.Buffer)
	DocumentFooter(out *bytes.Buffer)

	GetFlags() int
}

// Callback functions for inline parsing. One such function is defined
// for each character that triggers a response when parsing inline data
type inlineParser func(p *parser, out *bytes.Buffer, data []byte, offset int) int

// Parser holds runtime state used by the parser
type parser struct {
	r              Renderer
	refOverride    ReferenceOverrideFunc
	refs           map[string]*reference
	inlineCallback [256]inlineParser
	flags          int
	nesting        int
	maxNesting     int
	insideLink     bool
	notes          []*reference
}

// Reference represents the details of a link
type Reference struct {
	// Link is usually the URL the reference points to.
	Link string
	// Title is the alternate text describing the link in more details
	Title string
	// Text is the optional text to override the ref with if the syntax used was
	// [refid][]
	Tetxt string
}

// References are parsed and stored in this struct
type reference struct {
	link     []byte
	title    []byte
	noteId   int // o if no ta footnote ref
	hasBlock bool
	text     []byte
}

// ReferenceOverrideFunc is expected to be called with a reference string and
// return either a valid Reference type that the reference string maps to or
// nil. If overridden is false, the default reference logic will be executed.
// see the documentation in Options for more details on use-case
type ReferenceOverrideFunc func(reference string) (ref *Reference, overridden bool)

// Options represents configurable overrides and callbacks for configuring a Markdown parse
type Options struct {
	// Extensions is a flags set of bit-wise ORed extension bits. See the
	// EXTENSIONS_* flags defined in this package
	Extensions int

	// ReferenceOveerride is an optional function callback that is called every time
	// a reference is resolved.
	ReferenceOverride ReferenceOverrideFunc
}

// MarkdownBasic is a convenience function for simple renderring
// It processes markdowwn input with no extensions enabled
func MarkdownBasic(input []byte) []byte {
	// set up the HTML renderer
	htmlFlags := HTML_USE_XHTML
	renderer := HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	return MarkdownOptions(input, renderer, Options{Extensions: 0})
}

// Markdown is the main rendering function.
func Markdown(input []byte, renderer Renderer, extensions int) []byte {
	return MarkdownOptions(input, renderer, Options{Extensions: extensions})
}

// MarkdownOptions is just like Markdown but takes additional options through the Options struct
func MarkdownOptions(input []byte, renderer Renderer, opts Options) []byte {
	// If renderer is nil, we can not render
	if renderer == nil {
		return nil
	}

	extensions := opts.Extensions

	// fill in the render structure
	p := new(parser)
	p.r = renderer
	p.flags = extensions
	p.refOverride = opts.ReferenceOverride
	p.refs = make(map[string]*reference)
	p.maxNesting = 16
	p.insideLink = false

	// register inline parsers
	//	p.inlineCallback['*'] = emphasis
	//	p.inlineCallback['_'] = emphasis
	//	if extensions&EXTENSION_STRIKETHROUGH != 0 {
	//p.inlineCallback['~'] = emphasis
	//	}
	//	p.inlineCallback['`'] = codeSpan
	//	p.inlineCallback['\n'] = linkBreak
	//	p.inlineCallback['['] = link
	//	p.inlineCallback['<'] = leftAngle
	//	p.inlineCallback['\\'] = escape
	//	p.inlineCallback['&'] = entity

	//	if extensions&EXTENSION_AUTOLINK != 0 {
	//		p.inlineCallback[':'] = autoLink
	//	}

	//	if extensions&EXTENSION_FOOTNOTES != 0 {
	p.notes = make([]*reference, 0)
	//	}

	first := firstRender(p, input)
	second := secondRender(p, first)

	return second
}

// firstRender only does the following:
// - extrace references
// - expand tabs
// - normalize newlines
// - copy everything else
func firstRender(p *parser, input []byte) []byte {
	var out bytes.Buffer
	tabSize := TAB_SIZE_DEFAULT
	if p.flags&EXTENSION_TAB_SIZE_EIGHT != 0 {
		tabSize = TAB_SIZE_EIGHT
	}

	begin, end := 0, 0

	lastFencedCodeBlockEnd := 0

	for begin < len(input) { // iterate over lines
		//		if end = isReference(p, input[begin:], tabSize); end > 0 {
		//			begin += end
		//		} else { // skip to the next line
		//			end = begin
		for end < len(input) && input[end] != '\n' && input[end] != '\r' {
			end++
		}

		//		if p.flags&EXTENSION_FENCED_CODE != 0 {
		//			// track fenced code block boundaries to suppress tab expansion inside them
		//			if begin >= lastFencedCodeBlockEnd {
		//				if i := p.fencedCode(&out, input[begin:], false); i > 0 {
		//					lastFencedCodeBlockEnd = begin + i
		//				}
		//			}
		//		}

		// add the line body if present
		if end > begin {
			if end < lastFencedCodeBlockEnd { // do not expand tabs while inside fenced code blocks.
				out.Write(input[begin:end])
			} else {
				expandTabs(&out, input[begin:end], tabSize)
			}
		}
		out.WriteByte('\n')

		if end < len(input) && input[end] == '\r' {
			end++
		}
		if end < len(input) && input[end] == '\n' {
			end++
		}

		begin = end
		//}
	}

	// empty input>
	if out.Len() == 0 {
		out.WriteByte('\n')
	}

	return out.Bytes()
}

// secondRender: actual renderring
func secondRender(p *parser, input []byte) []byte {
	var out bytes.Buffer

	p.r.DocumentHeader(&out)
	p.block(&out, input)
	p.r.DocumentFooter(&out)

	if p.nesting != 0 {
		panic("Nesting level did not end at zero")
	}

	return out.Bytes()
}
